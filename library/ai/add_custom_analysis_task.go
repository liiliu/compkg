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

type CustomAnalysisParams struct {
	//任务名称
	TaskName string `json:"taskName"`
	//回调地址
	CallbackUrl string `json:"callbackUrl"`
	//回调参数 {"request_id": "1043"}
	CallbackParams string `json:"callbackParams"`
	//配置 {"input":{"llm":"Doubao-pro-32k","prompts":["以下是一段对话记录:\n{content}","将上述对话记录提炼出摘要信息，字数限制在200字以内"],"input_params":'{"content":content}',"out_formatter":'{"summary":"摘要内容"}'}}
	Config string `json:"config"`
}

type CustomAnalysisCallbackParams struct {
	ConversationId int64  `json:"conversationId"`
	AiAnalysisId   int64  `json:"aiAnalysisId"`
	Identify       string `json:"identify"`
}

type CustomAnalysisConfig struct {
	Input struct {
		Llm          string   `json:"llm"`
		Prompts      []string `json:"prompts"`
		InputParams  string   `json:"input_params"`
		OutFormatter string   `json:"out_formatter"`
	} `json:"input"`
}

type CustomAnalysisResult struct {
	Code int `json:"code"`
	Data struct {
		TaskID string `json:"taskId"`
	} `json:"data"`
	Msg    string `json:"msg"`
	Result int    `json:"result"`
}

func AddCustomAnalysisTask(params CustomAnalysisParams, token string) (taskId string, err error) {
	jsonStr := util.JsonToString(params)
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token

	u := config.GetString("ai.addAnalysisTaskUrl")
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
	logger.Info(common.LogTagCustomAiAnalysis, util.JsonToString(logMap))

	if err != nil {
		return
	}

	if respCode != 200 {
		err = errors.New("网络错误")
		return
	}

	var resp CustomAnalysisResult
	if err = json.Unmarshal(respBody, &resp); err != nil {
		logger.Error(common.LogTagCustomAiAnalysis, err.Error())
		return
	}

	//回调记录
	record := new(model.CallbackRecord)
	record.CallbackType = model.CallbackTypeCustomAiAnalysis
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
		logger.Error(common.LogTagCustomAiAnalysis, Err.Error())
		err = Err
		return
	}

	return
}
