package ai

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"weihu_server/library/common"
	"weihu_server/library/config"
	"weihu_server/library/database"
	"weihu_server/library/logger"
	"weihu_server/library/util"
	"weihu_server/model"
)

type AsrParams struct {
	//任务名称
	TaskName string `json:"taskName"`
	//语音文件url
	AudioFileUrl string `json:"audioFileUrl"`
	//回调地址
	CallbackUrl string `json:"callbackUrl"`
	//回调参数 {"request_id": "3045"}
	CallbackParams string `json:"callbackParams"`
	//配置 火山 {"input": {"need_normalize": false, "need_denoise": false, "is_stereo": true, "asr_component": "volc"},"out": {"audio_type": "mp3", "asr_type": "file"}}
	//配置 funasr {"input": {"need_normalize": false, "need_denoise": false, "is_stereo": true, "asr_component": "funasr","interlocutor_url": interlocutor_url},"out": {"audio_type": "mp3", "asr_type": "file"}}
	//配置 tencent {"input": {"need_normalize": false, "need_denoise": false, "is_stereo": true, "asr_component": "tencent"},"out": {"audio_type": "mp3", "asr_type": "file"}}
	Config string `json:"config"`
}

type AsrCallbackParams struct {
	ConversationId int64 `json:"conversationId"`
	AiAnalysisId   int64 `json:"aiAnalysisId"`
}

type AsrCfg struct {
	Input struct {
		NeedNormalize bool   `json:"need_normalize"`
		NeedDenoise   bool   `json:"need_denoise"`
		IsStereo      bool   `json:"is_stereo"`
		AsrComponent  string `json:"asr_component"`
	} `json:"input"`
	Out struct {
		AudioType string `json:"audio_type"`
		AsrType   string `json:"asr_type"`
	} `json:"out"`
}

type AsrResult struct {
	Code int `json:"code"`
	Data struct {
		TaskID string `json:"taskId"`
	} `json:"data"`
	Msg    string `json:"msg"`
	Result int    `json:"result"`
}

func AddAsrTask(params AsrParams, token string) (taskId string, err error) {
	jsonStr := util.JsonToString(params)
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token

	u := config.GetString("ai.addAsrTaskUrl")
	respCode, respBody, err := util.HttpPost(u, jsonStr, headers)
	logMap := make(map[string]interface{})
	logMap["url"] = u
	logMap["req"] = map[string]interface{}{
		"RequestHeader": headers,
		"RequestBody":   params,
	}
	logMap["respCode"] = respCode
	if gjson.ValidBytes(respBody) {
		bodyMap := util.ByteToMap(respBody)
		logMap["respBody"] = bodyMap
	} else {
		logMap["respBody"] = string(respBody)
	}
	logger.Info(common.LogTagAiASR, util.JsonToString(logMap))

	if err != nil {
		return
	}
	if respCode != 200 {
		err = errors.New("网络错误")
		return
	}

	var resp AsrResult
	if err = json.Unmarshal(respBody, &resp); err != nil {
		logger.Error(common.LogTagAiASR, err.Error())
		return
	}

	//回调记录
	record := new(model.CallbackRecord)
	record.CallbackType = model.CallbackTypeAiASR
	record.CallbackUrl = u
	record.CallbackBody = jsonStr
	record.CallbackResult = string(respBody)
	record.CallbackStatus = model.CallbackStatusSuccess
	if resp.Code != 1000 {
		err = errors.New(fmt.Sprintf("请求失败，错误码：%d，错误信息：%s", resp.Code, resp.Msg))
		record.CallbackStatus = model.CallbackStatusFail
	} else {
		taskId = resp.Data.TaskID
		record.CallbackStatus = model.CallbackStatusSuccess
	}

	if Err := database.MasterDb.Create(record).Error; Err != nil {
		logger.Error(common.LogTagAiAnalysis, Err.Error())
		err = Err
		return
	}

	return
}
