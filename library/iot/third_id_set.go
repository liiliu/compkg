package iot

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

type ThirdIdSetParams struct {
	//设备号
	DeviceNo string `json:"deviceNo"`
	//id1
	Id1 string `json:"id1"`
	//id2
	Id2 string `json:"id2"`
	//id3
	Id3 string `json:"id3"`
}

type ThirdIdSetResult struct {
	Result int64    `json:"result"`
	Msg    string   `json:"msg"`
	Code   int64    `json:"code"`
	Data   struct{} `json:"data"`
}

func ThirdIdSet(params ThirdIdSetParams, token string) error {
	jsonStr := util.JsonToString(params)
	headers := make(map[string]string)
	headers["Authorization"] = token

	u := config.GetString("iot.thirdIdSet")
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
	logger.Info(common.LogTagSetThirdPartyID, util.JsonToString(logMap))

	if err != nil {
		return err
	}
	if respCode != 200 {
		return errors.New("网络错误")
	}

	var resp ThirdIdSetResult
	if err = json.Unmarshal(respBody, &resp); err != nil {
		logger.Error(common.LogTagSetThirdPartyID, err.Error())
		return err
	}

	//回调记录
	record := new(model.CallbackRecord)
	record.CallbackType = model.CallbackTypeSetThirdId
	record.CallbackUrl = u
	record.CallbackBody = jsonStr
	record.CallbackResult = string(respBody)
	record.CallbackStatus = model.CallbackStatusSuccess
	if resp.Code != 1000 {
		err = errors.New(fmt.Sprintf("请求失败，错误码：%d，错误信息：%s", resp.Code, resp.Msg))
		record.CallbackStatus = model.CallbackStatusFail
	} else {
		record.CallbackStatus = model.CallbackStatusSuccess
	}

	if Err := database.MasterDb.Create(record).Error; Err != nil {
		logger.Error(common.LogTagSetThirdPartyID, Err.Error())
		return Err
	}

	return err
}
