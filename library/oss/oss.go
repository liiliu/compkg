package oss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"hash"
	"io"
	"strings"
	"time"
	"weihu_server/api/view/desk"
	"weihu_server/library/common"
	"weihu_server/library/config"
	"weihu_server/library/logger"
)

type Config struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
	HostUrl         string
	ExpireTime      int64
}

var ossConfig *Config

func Initial() {
	ossConfig = &Config{
		Endpoint:        config.GetString("oss.endPoint"),
		AccessKeyId:     config.GetString("oss.accessKeyId"),
		AccessKeySecret: config.GetString("oss.accessKeySecret"),
		Bucket:          config.GetString("oss.bucket"),
		HostUrl:         config.GetString("oss.hostUrl"),
		ExpireTime:      config.GetInt64("oss.expireTime"),
	}
}

func getGmtIso8601(expireEnd int64) string {
	var tokenExpire = time.Unix(expireEnd, 0).UTC().Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

// GetPolicyToken 获取临时Token
func GetPolicyToken(prefixPath string) *desk.PolicyToken {
	now := time.Now().Unix()
	expireEnd := now + ossConfig.ExpireTime
	tokenExpire := getGmtIso8601(expireEnd)

	//create post policy json
	configStruct := new(ConfigStruct)
	configStruct.Expiration = tokenExpire
	condition := make([]string, 0)
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, prefixPath)
	configStruct.Conditions = append(configStruct.Conditions, condition)

	//计算签名
	result, _ := json.Marshal(configStruct)
	deByte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(ossConfig.AccessKeySecret))
	io.WriteString(h, deByte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	policyToken := new(desk.PolicyToken)
	policyToken.AccessKeyId = ossConfig.AccessKeyId
	policyToken.Host = ossConfig.HostUrl
	policyToken.Expire = expireEnd
	policyToken.Signature = signedStr
	policyToken.Directory = prefixPath
	policyToken.Policy = deByte

	return policyToken
}

// PutFile 上传文件
func PutFile(localFilePath, targetFile string) (fullPath string, err error) {
	client, err := oss.New(ossConfig.Endpoint, ossConfig.AccessKeyId, ossConfig.AccessKeySecret)
	if err != nil {
		logger.Error(common.LogTagOssError, fmt.Sprintf("oss.New error: %s", err.Error()))
		return
	}

	bucket, err := client.Bucket(ossConfig.Bucket)
	if err != nil {
		logger.Error(common.LogTagOssError, fmt.Sprintf("Bucket Error: %s", err.Error()))
		return
	}

	err = bucket.PutObjectFromFile(targetFile, localFilePath)
	if err != nil {
		logger.Error(common.LogTagOssError, fmt.Sprintf("PutObjectFromFile Error: %s", err.Error()))
		return
	}

	fullPath = fmt.Sprintf("%s/%s", ossConfig.HostUrl, targetFile)

	return fullPath, nil
}

// CheckpointPutFile 断点续传 50M
func CheckpointPutFile(localFilePath, targetFile string) (fullPath string, err error) {
	client, err := oss.New(ossConfig.Endpoint, ossConfig.AccessKeyId, ossConfig.AccessKeySecret)
	if err != nil {
		logger.Error(common.LogTagOssError, fmt.Sprintf("oss.New error: %s", err.Error()))
		return
	}

	bucket, err := client.Bucket(ossConfig.Bucket)
	if err != nil {
		logger.Error(common.LogTagOssError, fmt.Sprintf("Bucket Error: %s", err.Error()))
		return
	}

	err = bucket.UploadFile(targetFile, localFilePath, 50*1024*1024, oss.Routines(3), oss.Checkpoint(true, ""))
	if err != nil {
		logger.Error(common.LogTagOssError, fmt.Sprintf("CheckpointPutFile error: %s", err.Error()))
		return
	}

	fullPath = fmt.Sprintf("%s/%s", ossConfig.HostUrl, targetFile)
	logger.Info(common.LogTagOssError, fmt.Sprintf("CheckpointPutFile success: %s", fullPath))
	return fullPath, nil
}

// GetTempUrl 获取临时访问地址
func GetTempUrl(targetFile string, expiredInSec int64) (signedURL string, err error) {
	client, err := oss.New(ossConfig.Endpoint, ossConfig.AccessKeyId, ossConfig.AccessKeySecret)
	if err != nil {
		logger.Error(common.LogTagOssError, fmt.Sprintf("oss.New error: %s", err.Error()))
		return
	}

	bucket, err := client.Bucket(ossConfig.Bucket)
	if err != nil {
		logger.Error(common.LogTagOssError, fmt.Sprintf("Bucket Error: %s", err.Error()))
		return
	}

	// 生成用于下载的签名URL，并指定签名URL的有效时间为60秒。
	signedURL, err = bucket.SignURL(targetFile, oss.HTTPGet, expiredInSec)
	if err != nil {
		logger.Error(common.LogTagOssError, fmt.Sprintf("SignURL Error: %s", err.Error()))
		return
	}
	return
}

// PutFileForTempUrl 上传文件
func PutFileForTempUrl(localFilePath, targetFile string) (signedURL string, err error) {
	client, err := oss.New(ossConfig.Endpoint, ossConfig.AccessKeyId, ossConfig.AccessKeySecret)
	if err != nil {
		logger.Error(common.LogTagOssError, fmt.Sprintf("oss.New error: %s", err.Error()))
		return
	}

	bucket, err := client.Bucket(config.GetString("ossPrivate.bucket"))
	if err != nil {
		logger.Error(common.LogTagOssError, fmt.Sprintf("Bucket Error: %s", err.Error()))
		return
	}

	options := oss.ContentDisposition("attachment; filename=" + targetFile)
	err = bucket.PutObjectFromFile(targetFile, localFilePath, options)
	if err != nil {
		logger.Error(common.LogTagOssError, fmt.Sprintf("PutObjectFromFile Error: %s", err.Error()))
		return
	}

	// 生成用于下载的签名URL，并指定签名URL的有效时间为60秒。
	var expiredInSec int64
	expiredInSec = 60

	signedURL, err = bucket.SignURL(targetFile, oss.HTTPGet, expiredInSec)
	if err != nil {
		logger.Error(common.LogTagOssError, fmt.Sprintf("SignURL Error: %s", err.Error()))
		return "", err
	}
	if config.GetString("server.env") == "PROD" {
		signedURL = strings.Replace(signedURL, "http://", "https://", 1)
	}
	return
}
