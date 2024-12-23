package ai

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"time"
	"weihu_server/library/common"
	"weihu_server/library/config"
	"weihu_server/library/logger"
	"weihu_server/library/util"
)

type GetTokenParams struct {
	AppID     string `json:"appId"`
	Sign      string `json:"sign"`
	Timestamp int64  `json:"timestamp"`
}

type GetTokenResult struct {
	Code int `json:"code"`
	Data struct {
		AccessToken string `json:"accessToken"`
		ExpiredTime int64  `json:"expiredTime"`
	} `json:"data"`
	Msg    string `json:"msg"`
	Result int    `json:"result"`
}

func GetToken() (accessToken string, expiredTime int64, err error) {
	timestamp := time.Now().Unix()
	sign := util.Md5(config.GetString("ai.appId") + util.Int64ToString(timestamp) + config.GetString("ai.appSecret"))
	params := GetTokenParams{
		AppID:     config.GetString("ai.appId"),
		Sign:      sign,
		Timestamp: timestamp,
	}
	jsonStr := util.JsonToString(params)
	u := config.GetString("ai.getToken")
	respCode, respBody, err := util.HttpPost(u, jsonStr, nil)
	logMap := make(map[string]interface{})
	logMap["url"] = u
	logMap["req"] = map[string]interface{}{
		"RequestBody": params,
	}
	logMap["respCode"] = respCode
	if gjson.ValidBytes(respBody) {
		bodyMap := util.ByteToMap(respBody)
		logMap["respBody"] = bodyMap
	} else {
		logMap["respBody"] = string(respBody)
	}
	logger.Info(common.LogTagGetToken, util.JsonToString(logMap))

	if err != nil {
		return
	}
	if respCode != 200 {
		err = errors.New("网络错误")
		return
	}

	var resp GetTokenResult
	if err = json.Unmarshal(respBody, &resp); err != nil {
		logger.Error(common.LogTagGetToken, err.Error())
		return
	}

	if resp.Code != 1000 {
		err = errors.New(resp.Msg)
		return
	}

	accessToken = resp.Data.AccessToken
	expiredTime = resp.Data.ExpiredTime

	return
}
