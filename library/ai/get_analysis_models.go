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

type GetAnalysisModelsResult struct {
	Code int `json:"code"`
	Data struct {
		Data []string `json:"data"`
	} `json:"data"`
	Msg    string `json:"msg"`
	Result int    `json:"result"`
}

func GetAnalysisModels(token string) (models []string, err error) {
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token

	u := config.GetString("ai.getAnalysisModels")
	respBody, err := util.HttpRequest(u, "GET", nil, headers)
	logMap := make(map[string]interface{})
	logMap["url"] = u
	logMap["req"] = map[string]interface{}{
		"RequestHeader": headers,
	}
	if gjson.ValidBytes(respBody) {
		bodyMap := util.ByteToMap(respBody)
		logMap["respBody"] = bodyMap
	} else {
		logMap["respBody"] = string(respBody)
	}
	logger.Info(common.LogGetAnalysisModels, util.JsonToString(logMap))

	if err != nil {
		return
	}

	var resp GetAnalysisModelsResult
	if err = json.Unmarshal(respBody, &resp); err != nil {
		logger.Error(common.LogTagAiAnalysis, err.Error())
		return
	}

	//回调记录
	record := new(model.CallbackRecord)
	record.CallbackType = model.CallbackGetAnalysisModels
	record.CallbackUrl = u
	record.CallbackBody = ""
	record.CallbackResult = string(respBody)
	record.CallbackStatus = model.CallbackStatusSuccess
	if resp.Code != 1000 {
		err = errors.New(fmt.Sprintf("请求失败，错误码：%d，错误信息：%s", resp.Code, resp.Msg))
		record.CallbackStatus = model.CallbackStatusFail
	} else {
		models = resp.Data.Data
		record.CallbackStatus = model.CallbackStatusSuccess
	}

	if Err := database.MasterDb.Create(record).Error; Err != nil {
		logger.Error(common.LogTagAiAnalysis, Err.Error())
		err = Err
		return
	}

	return
}
