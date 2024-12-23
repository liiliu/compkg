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

type GetAudioParams struct {
	//设备号
	DeviceNo string `json:"deviceNo"`
	//页码
	Page int `json:"page"`
	//每页数量
	Size int `json:"size"`
	//开始时间 yyyy-MM-dd HH:mm:ss
	StartTime string `json:"startTime"`
	//结束时间 yyyy-MM-dd HH:mm:ss
	EndTime string `json:"endTime"`
	//上传状态 0-未上传 1-已上传
	UploadStatus int `json:"uploadStatus"`
}

type GetAudioResult struct {
	Result int64     `json:"result"`
	Msg    string    `json:"msg"`
	Code   int64     `json:"code"`
	Data   AudioData `json:"data"`
}

type AudioData struct {
	Rows    []AudioRow `json:"rows"`
	Total   int        `json:"total"`
	Pages   int        `json:"pages"`
	PageNum int        `json:"pageNum"`
}

type AudioRow struct {
	DeviceNo     string `json:"deviceNo"`     //设备号
	EventId      string `json:"eventId"`      //录音id
	FileName     string `json:"fileName"`     //录音文件名
	FileSize     int64  `json:"fileSize"`     //录音文件大小B
	StartTime    string `json:"startTime"`    //录音开始时间,2020-10-01 09:23:23
	EndTime      string `json:"endTime"`      //录音结束时间,2020-10-01 09:23:23
	Seconds      int64  `json:"seconds"`      //录音时长,秒
	FileUrl      string `json:"fileUrl"`      //录音文件url
	LeftFileURL  string `json:"leftFileUrl"`  //左声道录音文件
	RightFileURL string `json:"rightFileUrl"` //右声道录音文件
	ID1          string `json:"id1"`
	ID2          string `json:"id2"`
	ID3          string `json:"id3"`
}

func GetAudio(params GetAudioParams, token string) (data AudioData, err error) {
	jsonStr := util.JsonToString(params)
	headers := make(map[string]string)
	headers["Authorization"] = token

	u := config.GetString("iot.getAudio")
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
		return
	}
	if respCode != 200 {
		err = errors.New("网络错误")
		return
	}

	var resp GetAudioResult
	if err = json.Unmarshal(respBody, &resp); err != nil {
		logger.Error(common.LogTagSetThirdPartyID, err.Error())
		return
	}

	//回调记录
	record := new(model.CallbackRecord)
	record.CallbackType = model.CallbackTypeGetDeviceInfo
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
		err = Err
		return
	}

	data = resp.Data

	return
}
