package common

const (
	// CacheTokenPrefix Token前缀
	CacheTokenPrefix string = "TOKEN"
	// CacheEmailCode 邮箱验证码
	CacheEmailCode string = "EmailCode"
	// CacheEmailResetToken 邮箱重置令牌
	CacheEmailResetToken string = "EmailResetToken"
	// CacheSmsCode 短信验证码
	CacheSmsCode string = "SmsCode"
	// CacheVerifyCodeErrorNum 验证码错误次数
	CacheVerifyCodeErrorNum string = "VerifyCodeErrorNum"
	// CachePrefixPasswordError 密码错误次数
	CachePrefixPasswordError = "PasswordError"
	// CacheServiceVersion 服务版本
	CacheServiceVersion string = "ServiceVersion"
	// CacheSendEmailMaxResendNum 发送邮件最大重发次数
	CacheSendEmailMaxResendNum string = "SendEmailMaxResendNum"
	// CacheEmailTemplateType 邮箱模板类型
	CacheEmailTemplateType = "EmailTemplateType"
	// CacheSendEmail 发送邮箱
	CacheSendEmail = "SendEmail"
	// CacheHelpCenter 帮助中心
	CacheHelpCenter = "helpCenter"
	// CacheAiToken AiToken
	CacheAiToken = "AiToken"
	// CacheIotToken IotToken
	CacheIotToken = "IotToken"
	// CacheLiveKitServiceStatus LiveKit服务状态
	CacheLiveKitServiceStatus = "LiveKitServiceStatus"
	// CacheShortLink 短链接
	CacheShortLink = "ShortLink"
	// CacheAsrModel AsrModel
	CacheAsrModel = "AsrModel"
	// CacheAnalysisModel AnalysisModel
	CacheAnalysisModel = "AnalysisModel"
)

// 队列
const (
	// QueuePushEmail 推送邮件
	QueuePushEmail = "QUEUE_PUSH_EMAIL"
	// QueuePushSms 推送短信
	QueuePushSms = "QUEUE_PUSH_SMS"
	// QueueWsServer 推送socket
	QueueWsServer = "QUEUE_WS_SERVER"
	// QueueSSE SSE数据
	QueueSSE = "QUEUE_SSE"
)

// 平台
const (
	// SysClient 客户端
	SysClient = "client"
	// SysBackend 运营后台
	SysBackend = "backend"
	// SysWeb 官网
	SysWeb = "web"
	// SysLiveKit 多人会议
	SysLiveKit = "liveKit"
)

const (
	SMSVerifyCodeTemplate   = "SMS_238152581" //验证码${code}，请在10分钟内按页面提示提交验证码，切勿将验证码泄露于他人。
	SMSJoinNoticeTemplate   = "SMS_207105027" //欢迎加入嘟嘟Talk，您的账户名为${account}，初始密码为${password}
	SMSSendCodeTemplate     = "SMS_206539891" //验证码${code}，该验证码5分钟内有效，请勿泄漏于他人
	SMSSubmitNoticeTemplate = "SMS_206549274" //尊敬的客户，感谢您提交试用申请，我司会尽快与您取得联系，请保持手机畅通
	SMSLoginCodeTemplate    = "SMS_205616179" //验证码：${code}。此验证码只用于短信登录，15分钟内有效
)

const (
	// EmailLoginTemplate 登录模块邮箱模板
	EmailLoginTemplate = "login"
)

// 日志标签
const (
	// LogTagInterfaceRequest 接口请求
	LogTagInterfaceRequest = "接口请求"
	// LogTagRedis Redis缓存
	LogTagRedis = "缓存"
	// LogTagMqtt MQTT
	LogTagMqtt = "MQTT"
	// LogTagDb 数据库
	LogTagDb = "数据库"
	// LogTagJob 推送任务
	LogTagJob = "下发任务"
	// LogTagEtl 同步任务
	LogTagEtl = "同步任务"
	// LogTagCheckToken 同步任务
	LogTagCheckToken = "同步CRM_TOKEN检测"
	// LogTagPushContactRecord 推送联系记录
	LogTagPushContactRecord = "推送联系记录"
	// LogTagApiSystemError API系统错误
	LogTagApiSystemError string = "API系统错误"
	// LogTagSendSms 发送短信
	LogTagSendSms string = "发送短信"
	// LogTagAwsError aws错误
	LogTagAwsError string = "AWS错误"
	// LogTagGoogleAPIGetToken googleAPI获取Token
	LogTagGoogleAPIGetToken string = "googleAPI获取Token"
	// LogTagGoogleAPIGetUser googleAPI获取用户信息
	LogTagGoogleAPIGetUser string = "googleAPI获取用户信息"
	// LogTagFacebookAPIGetToken facebookAPI获取Token
	LogTagFacebookAPIGetToken string = "facebookAPI获取Token"
	// LogTagFacebookAPIGetUser facebookAPI获取用户信息
	LogTagFacebookAPIGetUser string = "facebookAPI获取用户信息"
	// LogTagAsyncSendMailAt 异步发送邮件
	LogTagAsyncSendMailAt string = "异步发送邮件"
	// LogTagAsyncSendSmsAt 异步发送短信
	LogTagAsyncSendSmsAt string = "异步发送短信"
	// LogTagOssError oss错误
	LogTagOssError string = "OSS错误"
	// LogTagCosError  cos错误
	LogTagCosError string = "COS错误"
	// LogTagPushEmail 推送邮件
	LogTagPushEmail string = "推送邮件"
	// LogTagPushSms 推送短信
	LogTagPushSms string = "推送短信"
	// LogTagAiAnalysis AI分析
	LogTagAiAnalysis string = "AI分析"
	// LogTagCustomAiAnalysis 自定义AI分析
	LogTagCustomAiAnalysis string = "自定义AI分析"
	// LogTagReceiveAiAnalysisResult 接收AI分析结果
	LogTagReceiveAiAnalysisResult string = "接收AI分析结果"
	// LogTagReceiveCustomAiAnalysisResult 接收自定义AI分析结果
	LogTagReceiveCustomAiAnalysisResult string = "接收自定义AI分析结果"
	// LogTagAiASR ASR分析
	LogTagAiASR string = "ASR分析"
	// LogTagReceiveAiASRResult 接收ASR分析结果
	LogTagReceiveAiASRResult string = "接收ASR分析结果"
	// LogTagGetToken 获取Token
	LogTagGetToken string = "获取Token"
	// LogTagAudioCallback 录音回调
	LogTagAudioCallback string = "录音回调"
	// LogTagDeviceCallback 设备回调
	LogTagDeviceCallback string = "设备回调"
	// LogTagGetIotTOKEN 获取IotToken
	LogTagGetIotTOKEN string = "获取IotToken"
	// LogTagSetThirdPartyID 设置第三方ID
	LogTagSetThirdPartyID string = "设置第三方ID"
	// LogTagCheckTeamStatus 检查团队状态
	LogTagCheckTeamStatus string = "检查团队状态"
	// LogTagCheckTask 检查任务状态
	LogTagCheckTask string = "检查任务状态"
	// LogTagLiveKitError liveKit错误
	LogTagLiveKitError string = "liveKit错误"
	// LogTagLiveKitCallback 接收LiveKit 回调
	LogTagLiveKitCallback string = "liveKit回调"
	// LogTagWsServer tcp服务
	LogTagWsServer string = "ws服务"
	// LogGetAnalysisModels 获取分析模型
	LogGetAnalysisModels string = "获取分析模型"
	// LogSSE LogSSE
	LogSSE string = "SSE"
)

// 任务状态
const (
	// TaskStatusRunning 任务运行中
	TaskStatusRunning = "RUNNING"
	// TaskStatusDone 任务完成
	TaskStatusDone = "DONE"
)

// 会员类型
const (
	// MemberTypeTrial 试用版
	MemberTypeTrial = "TRIAL"
	// MemberTypeAdvanced 高级版
	MemberTypeAdvanced = "ADVANCED"
)

// oss上传路径
const (
	// ContactImport 导入联系人
	ContactImport = "weihu/import_contact"
	// TeamImport 导入团队
	TeamImport = "weihu/import_team"
	// TeamExport 导出团队
	TeamExport = "weihu/export_team"
	// TeamUserExport 导出团队人员
	TeamUserExport = "weihu/export_team_user"
	// ExportTrialApply 导出试用申请
	ExportTrialApply = "weihu/export_trial_apply"
	// OrderImport 导入订单
	OrderImport = "weihu/import_order"
	// EmailAttachment 邮件附件
	EmailAttachment = "weihu/email_attachment"
	// HelpCenter  帮助中心
	HelpCenter = "weihu/help_center"
	// ASRExport ASR导出
	ASRExport = "weihu/export_asr"
)

const (
	SendEmailMaxResendNum = 3
)

// 授权第三方平台
const (
	Facebook = "facebook"
	Google   = "google"
)

const (
	HelpCenterFileName = "help_center.md"
)
