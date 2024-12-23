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

type GetDeviceParams struct {
	//设备号
	DeviceNo string `json:"deviceNo"`
	//页码
	Page int `json:"page"`
	//每页数量
	Size int `json:"size"`
}

type GetDeviceResult struct {
	Result int64      `json:"result"`
	Msg    string     `json:"msg"`
	Code   int64      `json:"code"`
	Data   DeviceData `json:"data"`
}

type DeviceData struct {
	Rows    []DeviceRow `json:"rows"`
	Total   int         `json:"total"`
	Pages   int         `json:"pages"`
	PageNum int         `json:"pageNum"`
}

type DeviceRow struct {
	DeviceNo        string `json:"deviceNo"`        //设备号
	DeviceType      string `json:"deviceType"`      //设备类型
	OnlineStatus    int    `json:"onlineStatus"`    //设备状态 0-离线 1-空闲 2-忙碌
	LastOnlineTime  string `json:"lastOnlineTime"`  //最后在线时间
	LastOfflineTime string `json:"lastOfflineTime"` //最后离线时间
	JoinTime        string `json:"joinTime"`        //设备创建时间
	Sim1Status      int    `json:"sim1Status"`      //sim1状态 -1未插卡，0未注册,当前未搜索要注册的运营商，1已注册,归属地网络，2未注册,正在搜索要注册的运营商，3注册被拒绝，4未知状态，5已注册，漫游网络
	Sim2Status      int    `json:"sim2Status"`      // sim2状态 -1未插卡，0未注册,当前未搜索要注册的运营商，1已注册,归属地网络，2未注册,正在搜索要注册的运营商，3注册被拒绝，4未知状态，5已注册，漫游网络
	BleStatus       int    `json:"bleStatus"`       //蓝牙状态 0-未连接 1-已连接
	PhoneStatus     int    `json:"phoneStatus"`     //电话状态 0空闲，1去电中，2来电中，3通话中，4结束
	FirmwareVersion string `json:"firmwareVersion"` //固件版本号
	RemainPower     int    `json:"remainPower"`     //电量
	PendingFileNum  int    `json:"pendingFileNum"`  //待上传文件数
	MsgType         int    `json:"msgType"`         //0心跳，1设备在线，2设备离线，4录音上传，5录音上传结束，6SIM1插卡，7SIM1拔卡，8SIM2插卡，9SIM2拔卡，10蓝牙链接，11蓝牙断开，12录音开始，13录音结束，14通话开始，15通话结束，16通话录音保存开关开启，17通话录音保存开关关闭
	ReportTime      string `json:"reportTime"`      //上报时间
	ID1             string `json:"id1"`
	ID2             string `json:"id2"`
	ID3             string `json:"id3"`
}

func GetDevice(params GetDeviceParams, token string) (data DeviceData, err error) {
	jsonStr := util.JsonToString(params)
	headers := make(map[string]string)
	headers["Authorization"] = token

	u := config.GetString("iot.getDevice")
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

	var resp GetDeviceResult
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
