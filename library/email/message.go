package email

import (
	"bytes"
	"fmt"
	"github.com/emersion/go-imap"
	id "github.com/emersion/go-imap-id"
	"github.com/emersion/go-imap/client"
	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"io"
	"io/ioutil"
	"net"
	"strings"
	"time"
	"weihu_server/library/common"
	"weihu_server/library/cos"
	"weihu_server/library/util"
)

type Email struct {
	From               []*mail.Address // 发件人
	To                 []*mail.Address // 收件人
	Cc                 []*mail.Address // 抄送
	Bcc                []*mail.Address // 密送
	Subject            string          // 主题
	Date               time.Time       // 发送日期
	MessageID          string          // 邮件唯一标识
	ReplyTo            string          // 回复默认收件人地址
	InReplyTo          string          // 回复的邮箱
	References         []string        // 引用的邮件ID列表
	ReturnPath         string          // 反向路径或envelop sender
	Received           []string        // 传输路径记录
	MIMEVersion        string          // MIME版本
	ContentType        string          // 内容类型
	ContentDisposition string          // 内容处置方式，如附件处理
	Body               string          // 邮件正文内容，可能是纯文本或需要进一步解析的多部分
	Attachments        []Attachment    // 附件列表
}

type Attachment struct {
	Filename    string // 附件文件名
	Data        []byte // 附件数据
	ContentType string // 附件的MIME类型
	FileUrl     string // 附件的下载地址
}

func ListByUid(server, UserName, Password string) (err error, result []*Email) {
	c, err := loginEmail(server, UserName, Password)
	if err != nil {
		fmt.Println(err)
		return
	}
	idClient := id.NewClient(c)
	idClient.ID(
		id.ID{
			id.FieldName:    "IMAPClient",
			id.FieldVersion: "2.1.0",
		},
	)

	defer c.Close()

	mailBoxes := make(chan *imap.MailboxInfo, 10)
	mailBoxDone := make(chan error, 1)
	go func() {
		mailBoxDone <- c.List("", "*", mailBoxes)
	}()
	for box := range mailBoxes {
		if box.Name != "INBOX" {
			continue
		}
		fmt.Println("切换目录:", box.Name)
		mbox, err := c.Select(box.Name, false)
		// 选择收件箱
		if err != nil {
			fmt.Println("select inbox err: ", err)
			continue
		}
		if mbox.Messages == 0 {
			continue
		}

		// 选择收取邮件的时间段
		criteria := imap.NewSearchCriteria()
		// 收取7天之内的邮件
		criteria.Since = time.Now().Add(-3 * time.Hour * 24)
		//criteria.Since = time.Now().Add(-4 * time.Hour)
		// 按条件查询邮件
		ids, err := c.UidSearch(criteria)
		fmt.Println(len(ids))
		if err != nil {
			continue
		}
		if len(ids) == 0 {
			continue
		}
		seqSet := new(imap.SeqSet)
		seqSet.AddNum(ids...)
		sect := &imap.BodySectionName{Peek: true}

		messages := make(chan *imap.Message, 100)
		messageDone := make(chan error, 1)

		go func() {
			messageDone <- c.UidFetch(seqSet, []imap.FetchItem{sect.FetchItem()}, messages)
		}()
		for msg := range messages {
			fmt.Printf("邮件UID: %d\n", msg.Uid)
			//fmt.Printf("msg: %v\n", msg)
			r := msg.GetBody(sect)
			mr, err := mail.CreateReader(r)
			if err != nil {
				fmt.Println(err)
				continue
			}
			header := mr.Header

			emailData := new(Email)
			if from, err := header.AddressList("From"); err == nil {
				emailData.From = from
			}
			if to, err := header.AddressList("To"); err == nil {
				emailData.To = to
			}
			if cc, err := header.AddressList("Cc"); err == nil {
				emailData.Cc = cc
			}
			subject, _ := header.Subject()
			emailData.Subject = subject

			// 提取Date字段
			date, err := header.Date()
			if err != nil {
				fmt.Println("无法解析日期:", err)
				continue
			}
			emailData.Date = date

			//提取messageId
			messageID, err := header.MessageID()
			if err != nil {
				fmt.Println("无法解析messageId:", err)
				//continue
			}
			emailData.MessageID = messageID

			//回复默认收件人地址
			emailData.ReplyTo = header.Get("Reply-To")

			//回复的邮箱
			emailData.InReplyTo = header.Get("In-Reply-To")

			//引用的邮件ID列表
			references, err := header.Text("References")
			if err != nil {
				fmt.Println("无法解析References:", err)
				continue
			}
			references = util.RemoveMultiChar(references, []string{"<", ">"})
			emailData.References = strings.Split(references, " ")

			// 获取并打印Reply-To
			if from, err := header.AddressList("Reply-To"); err == nil {
				for _, addr := range from {
					fmt.Println("Reply-To:", addr.Address)
				}
			}

			_, fileMap, _ := parseEmail(mr)

			//result = append(result, results...)
			for _, attachment := range fileMap {
				fmt.Println("收取到附件:", attachment.Filename)
				// 保存附件
				filePath, err := saveAttachment(attachment.Filename, attachment.Data)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("附件保存路径:", filePath)

				emailData.Attachments = append(emailData.Attachments, Attachment{
					Filename: attachment.Filename,
					//Data:        attachment.Data,
					ContentType: attachment.ContentType,
					FileUrl:     filePath,
				})
			}
			result = append(result, emailData)
			fmt.Printf("emailData = %+v\n", emailData)
		}
	}
	return
}

func parseEmail(mr *mail.Reader) (body []byte, fileMap []Attachment, results []string) {
	for {
		p, err := mr.NextPart()

		if err == io.EOF {
			break
		} else if err != nil {
			break
		}
		if p != nil {
			switch h := p.Header.(type) {
			case *mail.InlineHeader:
				body, err = ioutil.ReadAll(p.Body)
				if err != nil {
					fmt.Println("read body err:", err.Error())
				}
				//fmt.Println("<<<---------------------------------------------------------")
				//fmt.Println(string(body))
				//fmt.Println("--------------------------------------------------------->>>")
				results = append(results, string(body))

			case *mail.AttachmentHeader:
				contentType := strings.Split(p.Header.Get("Content-Type"), ";")[0]
				fileName, _ := h.Filename()
				fileContent, _ := ioutil.ReadAll(p.Body)
				fileMap = append(fileMap, Attachment{
					Filename:    fileName,
					Data:        fileContent,
					ContentType: contentType,
				})
			}
		}
	}
	return
}

func loginEmail(server, UserName, Password string) (*client.Client, error) {
	dial := new(net.Dialer)
	dial.Timeout = time.Duration(3) * time.Second
	c, err := client.DialWithDialerTLS(dial, server, nil)
	if err != nil {
		c, err = client.DialWithDialer(dial, server) // 非加密登录
	}
	if err != nil {
		return nil, err
	}
	// 登陆
	if err = c.Login(UserName, Password); err != nil {
		return nil, err
	}
	return c, nil
}

func saveAttachment(fileName string, fileContent []byte) (path string, err error) {
	// 文件上传到云端
	reader := bytes.NewBuffer(fileContent)
	path = fmt.Sprintf("%s/%s", common.EmailAttachment, fileName)
	err = cos.PutContent(path, reader)
	return
}

// parseEmailDate 尝试解析邮件的Date字段
func parseEmailDate(dateStr string) (time.Time, error) {
	// 尝试多种解析格式，因为邮件Date的格式可能有所不同
	for _, layout := range []string{
		time.RFC822,           // 常见的邮件日期格式
		time.RFC822Z,          // 带时区的RFC822格式
		"02 Jan 06 15:04 MST", // ANSI C 格式
	} {
		date, err := time.Parse(layout, dateStr)
		if err == nil {
			return date, nil
		}
	}

	// 如果所有尝试都失败，则返回错误
	return time.Time{}, fmt.Errorf("date format not recognized: %s", dateStr)
}
