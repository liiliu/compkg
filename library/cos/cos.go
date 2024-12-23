package cos

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
	"weihu_server/library/common"
	"weihu_server/library/config"
	"weihu_server/library/logger"
)

var cosClient *cos.Client

func Initial() {
	u, _ := url.Parse(config.GetString("cos.url"))
	b := &cos.BaseURL{BucketURL: u}
	// 1.永久密钥
	cosClient = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.GetString("cos.secretId"),
			SecretKey: config.GetString("cos.secretKey"),
		},
	})
}

// PutFile 上传文件
func PutFile(localFilePath, targetFile string) (string, error) {
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			XCosStorageClass: "STANDARD_IA",
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			XCosACL: "public-read",
		},
	}
	_, err := cosClient.Object.PutFromFile(context.Background(), targetFile, localFilePath, opt)
	if err != nil {
		logger.Error(common.LogTagCosError, fmt.Sprintf("PutFromFile Error: %s", err.Error()))
		return "", err
	}
	return fmt.Sprintf("%s/%s", config.GetString("cos.url"), targetFile), nil
}

// Upload 上传文件
func Upload(localFilePath, targetFile string) (string, error) {
	_, _, err := cosClient.Object.Upload(context.Background(), targetFile, localFilePath, nil)
	if err != nil {
		logger.Error(common.LogTagCosError, fmt.Sprintf("Upload Error: %s", err.Error()))
		return "", err
	}
	return fmt.Sprintf("%s/%s", config.GetString("cos.url"), targetFile), nil
}

// PutContent 直接上传内容
func PutContent(targetFile string, content io.Reader) error {
	_, err := cosClient.Object.Put(context.Background(), targetFile, content, nil)
	if err != nil {
		logger.Error(common.LogTagCosError, fmt.Sprintf("Put Error: %s", err.Error()))
		return err
	}
	return nil
}

// DownloadFile 下载文件
func DownloadFile(localFilePath, targetFile string) error {
	_, err := cosClient.Object.Download(context.Background(), targetFile, localFilePath, nil)
	if err != nil {
		logger.Error(common.LogTagCosError, fmt.Sprintf("Download Error: %s", err.Error()))
		return err
	}
	return nil
}

// GetPreSignedURL 获取预授权url
func GetPreSignedURL(targetFile string, expire time.Duration) (signedURL string, err error) {
	preSignedURL, err := cosClient.Object.GetPresignedURL(context.Background(), http.MethodPut, targetFile, config.GetString("cos.secretId"), config.GetString("cos.secretKey"), expire, nil)
	if err != nil {
		logger.Error(common.LogTagCosError, fmt.Sprintf("GetPresignedURL Error: %s", err.Error()))
		return
	}
	signedURL = preSignedURL.String()
	return
}

// DeleteFile 删除文件
func DeleteFile(targetFile string) error {
	_, err := cosClient.Object.Delete(context.Background(), targetFile)
	if err != nil {
		logger.Error(common.LogTagCosError, fmt.Sprintf("Delete Error: %s", err.Error()))
		return err
	}
	return nil
}

// GetFileUrl 获取文件临时地址
func GetFileUrl(targetFile string, expire time.Duration) (fileUrl string, err error) {
	preSignedURL, err := cosClient.Object.GetPresignedURL(context.Background(), http.MethodGet, targetFile, config.GetString("cos.secretId"), config.GetString("cos.secretKey"), expire, nil)
	if err != nil {
		return
	}
	fileUrl = preSignedURL.String()
	return
}

// GetCredential 获取临时签名信息
func GetCredential() (cs *sts.CredentialResult, err error) {
	c := sts.NewClient(
		// 通过环境变量获取密钥, os.Getenv 方法表示获取环境变量
		config.GetString("cos.secretId"),
		config.GetString("cos.secretKey"),
		nil,
	)
	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          config.GetString("cos.region"),
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						"*",
					},
					Effect: "allow",
					Resource: []string{
						"*",
					},
				},
			},
		},
	}
	res, err := c.GetCredential(opt)
	if err != nil {
		log.Printf("GetCredential err = %s\n", err.Error())
		return nil, err
	}
	return res, nil
}
