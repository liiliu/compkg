package refresh

import (
	"fmt"
	"gorm.io/gorm"
	"weihu_server/library/cache"
	"weihu_server/library/common"
	"weihu_server/library/database"
	"weihu_server/library/logger"
	"weihu_server/library/util"
	"weihu_server/model"
)

func BaseData() {
	AllDictionary()
	SendEmailMaxResendNum()
	LiveKitServiceStatus()
	AsrModel()
	AnalysisModel()
}

// AllDictionary 获取数据字典数据存入缓存
func AllDictionary() {
	// 查询所有数据
	dictionaries := make([]*model.Dictionary, 0)
	if err := database.BackendReadDb.Where("status = ? ", model.StatusNormal).Order("sort asc").Find(&dictionaries).Error; err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return
	}

	// 封装数据
	dictionaryMap := make(map[string][]*model.Dictionary)
	for _, dictionary := range dictionaries {
		dictionaryMap[dictionary.Category] = append(dictionaryMap[dictionary.Category], dictionary)
	}

	// 服务版本
	if val, ok := dictionaryMap[model.DicServiceVersion]; ok {
		cache.Save(common.CacheServiceVersion, val, 0)
	}

	// 邮件模板类型
	if val, ok := dictionaryMap[model.DicEmailTemplateType]; ok {
		cache.Save(common.CacheEmailTemplateType, val, 0)
	}

	// 发送邮箱
	if val, ok := dictionaryMap[model.DicSendEmail]; ok {
		cache.Save(common.CacheSendEmail, val, 0)
	}
}

// SendEmailMaxResendNum 发送邮件最大重发次数
func SendEmailMaxResendNum() {
	dic := new(model.Dictionary)
	if err := database.BackendReadDb.Where("category = ? and status = ? ", model.DicSendEmailMaxResendNum, model.StatusNormal).First(dic).Error; err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(common.LogTagRedis, err.Error())
		return
	}
	if dic.ID == 0 {
		dic = new(model.Dictionary)
		dic.Category = model.DicSendEmailMaxResendNum
		dic.CategoryName = "发送邮件最大重发次数"
		dic.DataName = "发送邮件最大重发次数"
		dic.DataDisplay = "发送邮件最大重发次数"
		dic.DataValue = fmt.Sprintf("%d", common.SendEmailMaxResendNum)
		dic.Remark = "发送邮件最大重发次数"
		dic.Status = model.StatusNormal
		if err := database.MasterDb.Create(dic).Error; err != nil {
			logger.Error(common.LogTagRedis, err.Error())
			return
		}
	}

	cache.Save(common.CacheSendEmailMaxResendNum, common.SendEmailMaxResendNum, 0)
}

// GetEmailTemplateType 获取邮件模板类型
func GetEmailTemplateType() {
	emailTemplateTypes := make([]*model.Dictionary, 0)
	if err := database.BackendReadDb.Where("category = ? and status = ? ", model.DicEmailTemplateType, model.StatusNormal).Order("sort asc").Find(&emailTemplateTypes).Error; err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return
	}

	cache.Save(common.CacheEmailTemplateType, emailTemplateTypes, 0)
}

// GetSendEmail 获取发送邮箱
func GetSendEmail() {
	sendEmails := make([]*model.Dictionary, 0)
	if err := database.BackendReadDb.Where("category = ? and status = ? ", model.DicSendEmail, model.StatusNormal).Order("sort asc").Find(&sendEmails).Error; err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return
	}

	cache.Save(common.CacheSendEmail, sendEmails, 0)
}

// LiveKitServiceStatus LiveKit服务状态
func LiveKitServiceStatus() {
	dic := new(model.Dictionary)
	if err := database.BackendReadDb.Where("category = ? and status = ? ", model.DicLiveKitServiceStatus, model.StatusNormal).First(dic).Error; err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(common.LogTagRedis, err.Error())
		return
	}
	if dic.ID == 0 {
		dic = new(model.Dictionary)
		dic.Category = model.DicLiveKitServiceStatus
		dic.CategoryName = "LiveKit服务状态"
		dic.DataName = "LiveKit服务状态"
		dic.DataDisplay = "启用"
		dic.Remark = "服务状态(启用,停用)"
		dic.Status = model.StatusNormal
		if err := database.MasterDb.Create(dic).Error; err != nil {
			logger.Error(common.LogTagRedis, err.Error())
			return
		}
	}

	cache.SaveString(common.CacheLiveKitServiceStatus, dic.DataDisplay, 0)
}

func AsrModel() {
	dic := new(model.Dictionary)
	if err := database.BackendReadDb.Where("category = ? and status = ? ", model.DicAsrModel, model.StatusNormal).First(dic).Error; err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(common.LogTagRedis, err.Error())
		return
	}
	if dic.ID == 0 {
		dic = new(model.Dictionary)
		dic.Category = model.DicAsrModel
		dic.CategoryName = "Asr模型"
		dic.DataName = "Asr模型"
		dic.DataDisplay = "腾讯"
		dic.DataValue = fmt.Sprintf("%d", model.AsrModelTencent)
		dic.Remark = "火山-1,fun-2,腾讯-3"
		dic.Status = model.StatusNormal
		if err := database.MasterDb.Create(dic).Error; err != nil {
			logger.Error(common.LogTagRedis, err.Error())
			return
		}
	}

	cache.Save(common.CacheAsrModel, util.StringToInt64(dic.DataValue), 0)
}

func AnalysisModel() {
	dic := new(model.Dictionary)
	if err := database.BackendReadDb.Where("category = ? and status = ? ", model.DicAnalysisModel, model.StatusNormal).First(dic).Error; err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(common.LogTagRedis, err.Error())
		return
	}
	if dic.ID == 0 {
		dic = new(model.Dictionary)
		dic.Category = model.DicAnalysisModel
		dic.CategoryName = "Analysis模型"
		dic.DataName = "Analysis模型"
		dic.DataDisplay = ""
		dic.DataValue = "qwen:long"
		dic.Remark = ""
		dic.Status = model.StatusNormal
		if err := database.MasterDb.Create(dic).Error; err != nil {
			logger.Error(common.LogTagRedis, err.Error())
			return
		}
	}

	cache.SaveString(common.CacheAnalysisModel, dic.DataValue, 0)
}
