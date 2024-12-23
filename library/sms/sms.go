package sms

import (
	"encoding/json"
	aliyunsmsclient "github.com/KenmyZhang/aliyun-communicate"
	"weihu_server/library/common"
	"weihu_server/library/config"
	"weihu_server/library/logger"
)

// SendSms 发送短信
func SendSms(mobile string, param string, template string) (resp string, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(common.LogTagSendSms, "SendSms panic: %v", err)
		}
	}()
	debug := config.GetBool("sms.debug")
	if debug {
		logger.Debug(common.LogTagSendSms, "模拟发送验证码短信")
		return
	}
	smsClient := aliyunsmsclient.New(config.GetString("sms.method"))
	result, err := smsClient.Execute(config.GetString("sms.appkey"), config.GetString("sms.appSecret"), mobile, config.GetString("sms.v"), template, param)
	if err != nil {
		logger.Error(common.LogTagSendSms, "SendSms error: %v", err)
		return
	}
	if result == nil {
		logger.Error(common.LogTagSendSms, "SendSms error: result is nil")
		return
	}

	logger.Debug("Got raw response from server:", string(result.RawResponse))
	if err != nil {
		logger.Error(common.LogTagSendSms, "SendSms error: %v", err)
		return
	}
	resultJson, err := json.Marshal(result)
	if err != nil {
		logger.Error(common.LogTagSendSms, "SendSms error: %v", err)
		return
	}
	if result.IsSuccessful() {
		logger.Debug(common.LogTagSendSms, "A SMS is sent successfully:", resultJson)
	} else {
		logger.Debug(common.LogTagSendSms, "Failed to send a SMS:", resultJson)
	}
	resp = string(resultJson)

	return
}
