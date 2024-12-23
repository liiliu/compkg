package aws

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	"weihu_server/library/common"
	"weihu_server/library/config"
	"weihu_server/library/logger"
)

type Config struct {
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
	DefaultRegion   string
	ExpireTime      time.Duration
}

var awsConfig *Config
var Sess *session.Session

func Initial() {
	awsConfig = &Config{
		AccessKeyId:     config.GetString("aws.accessKeyId"),
		AccessKeySecret: config.GetString("aws.accessKeySecret"),
		Bucket:          config.GetString("aws.bucket"),
		DefaultRegion:   config.GetString("aws.defaultRegion"),
		ExpireTime:      time.Duration(config.GetInt64("aws.expireTime")) * time.Second,
	}

	// 初始化AWS Session，使用默认区域和配置
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsConfig.DefaultRegion), // 替换为你的AWS区域
		Credentials: credentials.NewStaticCredentials(
			awsConfig.AccessKeyId,     // 替换为你的 AWS 访问密钥 ID
			awsConfig.AccessKeySecret, // 替换为你的 AWS 密钥
			"",                        // 如果使用的是 IAM 角色，此处可以留空；否则填入会话令牌（Session Token）
		),
	})
	if err != nil {
		logger.Error(common.LogTagAwsError, fmt.Sprintf("Error creating AWS session: %s", err.Error()))
		return
	}
	Sess = sess
}

func Upload(localFilePath, targetFile string) (fileUrl string, err error) {
	uploader := s3manager.NewUploader(Sess)

	file, err := os.Open(localFilePath)
	if err != nil {
		logger.Error(common.LogTagAwsError, fmt.Sprintf("Error opening file: %s", err.Error()))
		return
	}
	defer file.Close()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(awsConfig.Bucket),
		Key:    aws.String(targetFile),
		Body:   file,
	})
	if err != nil {
		logger.Error(common.LogTagAwsError, fmt.Sprintf("Error uploading file: %s", err.Error()))
		return
	}

	// 默认返回公有权限的url
	//https://<bucket-name>.s3.<region>.amazonaws.com/<object-key>
	fileUrl = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", awsConfig.Bucket, config.GetString("aws.defaultRegion"), targetFile)

	return
}

// GetFileUrl 获取文件临时地址
func GetFileUrl(targetFile string, expire time.Duration) (fileUrl string, err error) {
	// 创建S3服务客户端
	svc := s3.New(Sess)

	// 设置URL过期时间
	if expire == 0 {
		expire = awsConfig.ExpireTime
	}

	// Generate a pre-signed URL for the uploaded object that is publicly accessible.
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(awsConfig.Bucket),
		Key:    aws.String(targetFile),
	})
	fileUrl, err = req.Presign(expire) // Set the expiration time to 15 minutes
	if err != nil {
		logger.Error(common.LogTagAwsError, fmt.Sprintf("Error generating presigned URL: %s", err.Error()))
		return
	}
	return
}

// DownloadFile 下载文件
func DownloadFile(targetFile string, localFilePath string) (err error) {
	//文件路径不存在则创建
	err = os.MkdirAll(filepath.Dir(localFilePath), 0755)
	if err != nil {
		logger.Error(common.LogTagAwsError, fmt.Sprintf("Error creating directory for downloaded file: %s", err.Error()))
		return
	}

	// 创建 S3 客户端
	svc := s3.New(Sess)

	// 下载文件
	input := &s3.GetObjectInput{
		Bucket: aws.String(awsConfig.Bucket),
		Key:    aws.String(targetFile),
	}

	// 执行 GetObject 请求
	result, err := svc.GetObject(input)
	if err != nil {
		logger.Error(common.LogTagAwsError, fmt.Sprintf("Error downloading file: %s", err.Error()))
		return
	}
	defer result.Body.Close()

	// 将文件内容读取到内存中
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, result.Body)
	if err != nil {
		logger.Error(common.LogTagAwsError, fmt.Sprintf("Error reading downloaded file content: %s", err.Error()))
		return
	}

	// 将内容保存到本地文件或者进一步处理
	err = ioutil.WriteFile(localFilePath, buf.Bytes(), 0644)
	if err != nil {
		logger.Error(common.LogTagAwsError, fmt.Sprintf("Error writing downloaded file content to local file: %s", err.Error()))
	}
	return
}

// DeleteFile 删除文件
func DeleteFile(targetFile string) (err error) {
	// 创建 S3 客户端
	svc := s3.New(Sess)

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(awsConfig.Bucket),
		Key:    aws.String(targetFile),
	})
	if err != nil {
		logger.Error(common.LogTagAwsError, fmt.Sprintf("Error waiting for file to be deleted: %s", err.Error()))
		return
	}
	return
}

// GetPreSignedURL 获取预签名URL
func GetPreSignedURL(targetFile string, expire time.Duration) (fileUrl string, err error) {
	// 创建S3服务客户端
	svc := s3.New(Sess)

	// 设置URL过期时间
	if expire == 0 {
		expire = awsConfig.ExpireTime
	}

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(awsConfig.Bucket),
		Key:    aws.String(targetFile),
	})
	fileUrl, err = req.Presign(expire)
	if err != nil {
		logger.Error(common.LogTagAwsError, fmt.Sprintf("Error generating presigned URL: %s", err.Error()))
		return
	}
	return
}
