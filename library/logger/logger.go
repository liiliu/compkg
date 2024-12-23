package logger

import (
	"encoding/json"
	"fmt"
	"github.com/Arvintian/loki-client-go/loki"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gogo/protobuf/proto"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/tidwall/gjson"
	"runtime"
	"time"
	"weihu_server/library/config"
	"weihu_server/library/loghub"
	"weihu_server/library/util"
)

var ConsoleLogger *slog.Record

var LokiLogger *loki.Client

func Init() {
	slog.SetLogLevel(slog.TraceLevel)
	//h1 := handler.MustFileHandler("./log/error.txt", handler.WithLogLevels(slog.DangerLevels))
	//h1.SetFormatter(slog.NewJSONFormatter())
	//
	//h2 := handler.MustFileHandler("./log/normal.txt", handler.WithLogLevels(slog.Levels{slog.InfoLevel, slog.DebugLevel}))
	//h2.SetFormatter(slog.NewJSONFormatter())
	//
	//h3 := handler.MustFileHandler("./log/api.txt", handler.WithLogLevels(slog.Levels{slog.NoticeLevel, slog.TraceLevel}))
	//h3.SetFormatter(slog.NewJSONFormatter())

	h4 := handler.NewConsoleHandler(slog.Levels{slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel, slog.NoticeLevel, slog.InfoLevel, slog.DebugLevel})
	h4.SetFormatter(slog.NewTextFormatter())

	//consoleLogger := slog.NewWithHandlers(h1, h2, h3, h4)
	consoleLogger := slog.NewWithHandlers(h4)
	if config.GetBool("server.debug") {
		h5 := handler.NewConsoleHandler(slog.Levels{slog.TraceLevel})
		h5.SetFormatter(slog.NewTextFormatter())
		consoleLogger = slog.NewWithHandlers(h4, h5)
	}

	consoleLogger.ReportCaller = true
	consoleLogger.CallerFlag = 1
	consoleLogger.CallerSkip = consoleLogger.CallerSkip + 1

	ConsoleLogger = consoleLogger.WithFields(slog.M{
		"name": config.GetString("server.name"),
		"env":  config.GetString("server.env"),
	})

	if config.GetBool("server.loki_push_switch") {
		LokiLogger, _ = loki.NewWithDefault(config.GetString("server.loki_push_url"))
	}
}

func Error(tag string, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	lokiLabels := map[string]string{
		"application": config.GetString("server.application"),
		"server":      config.GetString("server.name"),
		"env":         config.GetString("server.env"),
		"level":       "ERROR",
		"tag":         tag,
	}
	//if labels != nil {
	//	for key, value := range labels {
	//		lokiLabels[key] = value
	//	}
	//}

	_, file, line, _ := runtime.Caller(1)
	jsonObject := make(map[string]interface{})
	jsonObject["file"] = fmt.Sprintf("%s:%d", file, line)
	if gjson.Valid(msg) {
		msgMap := util.JsonToMap(msg)
		jsonObject["msg"] = msgMap
	} else {
		jsonObject["msg"] = msg
	}
	for key, value := range lokiLabels {
		jsonObject[key] = value
	}
	jsonPost, _ := json.Marshal(jsonObject)
	if config.GetBool("server.loki_push_switch") {
		LokiLogger.Handle(lokiLabels, time.Now(), string(jsonPost))
	}

	if config.GetBool("logHub.pushSwitch") {
		projectName := config.GetString("logHub.projectName")
		logStoreName := config.GetString("logHub.logStoreName")

		logs := make([]*sls.Log, 0)
		contents := make([]*sls.LogContent, 0)

		contents = append(contents, &sls.LogContent{
			Key:   proto.String("env"),
			Value: proto.String(config.GetString("server.env")),
		})

		contents = append(contents, &sls.LogContent{
			Key:   proto.String("level"),
			Value: proto.String("ERROR"),
		})

		contents = append(contents, &sls.LogContent{
			Key:   proto.String("tag"),
			Value: proto.String(tag),
		})

		contents = append(contents, &sls.LogContent{
			Key:   proto.String("content"),
			Value: proto.String(string(jsonPost)),
		})

		log := new(sls.Log)
		log.Time = proto.Uint32(uint32(time.Now().Unix()))
		log.Contents = contents

		logs = append(logs, log)
		logGroup := new(sls.LogGroup)
		logGroup.Logs = logs
		// 推送日志到阿里云日志服务器
		loghub.PutLogs(projectName, logStoreName, logGroup)
	}

	ConsoleLogger.WithTime(time.Now()).Error(jsonObject)
}

func Debug(tag string, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	lokiLabels := map[string]string{
		"application": config.GetString("server.application"),
		"server":      config.GetString("server.name"),
		"env":         config.GetString("server.env"),
		"level":       "DEBUG",
		"tag":         tag,
	}
	//if labels != nil {
	//	for key, value := range labels {
	//		lokiLabels[key] = value
	//	}
	//}

	_, file, line, _ := runtime.Caller(1)
	jsonObject := make(map[string]interface{})
	jsonObject["file"] = fmt.Sprintf("%s:%d", file, line)
	if gjson.Valid(msg) {
		msgMap := util.JsonToMap(msg)
		jsonObject["msg"] = msgMap
	} else {
		jsonObject["msg"] = msg
	}
	for key, value := range lokiLabels {
		jsonObject[key] = value
	}
	jsonPost, _ := json.Marshal(jsonObject)
	if config.GetBool("server.loki_push_switch") {
		LokiLogger.Handle(lokiLabels, time.Now(), string(jsonPost))
	}
	if config.GetBool("logHub.pushSwitch") {
		projectName := config.GetString("logHub.projectName")
		logStoreName := config.GetString("logHub.logStoreName")

		logs := make([]*sls.Log, 0)
		contents := make([]*sls.LogContent, 0)

		contents = append(contents, &sls.LogContent{
			Key:   proto.String("env"),
			Value: proto.String(config.GetString("server.env")),
		})

		contents = append(contents, &sls.LogContent{
			Key:   proto.String("level"),
			Value: proto.String("DEBUG"),
		})

		contents = append(contents, &sls.LogContent{
			Key:   proto.String("tag"),
			Value: proto.String(tag),
		})

		contents = append(contents, &sls.LogContent{
			Key:   proto.String("content"),
			Value: proto.String(string(jsonPost)),
		})

		log := new(sls.Log)
		log.Time = proto.Uint32(uint32(time.Now().Unix()))
		log.Contents = contents

		logs = append(logs, log)
		logGroup := new(sls.LogGroup)
		logGroup.Logs = logs
		// 推送日志到阿里云日志服务器
		loghub.PutLogs(projectName, logStoreName, logGroup)

	}
	ConsoleLogger.WithTime(time.Now()).Error(jsonObject)
}

func Info(tag string, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	lokiLabels := map[string]string{
		"application": config.GetString("server.application"),
		"server":      config.GetString("server.name"),
		"env":         config.GetString("server.env"),
		"level":       "INFO",
		"tag":         tag,
	}
	//if labels != nil {
	//	for key, value := range labels {
	//		lokiLabels[key] = value
	//	}
	//}

	_, file, line, _ := runtime.Caller(1)
	jsonObject := make(map[string]interface{})
	jsonObject["file"] = fmt.Sprintf("%s:%d", file, line)
	if gjson.Valid(msg) {
		msgMap := util.JsonToMap(msg)
		jsonObject["msg"] = msgMap
	} else {
		jsonObject["msg"] = msg
	}
	for key, value := range lokiLabels {
		jsonObject[key] = value
	}
	jsonPost, _ := json.Marshal(jsonObject)
	if config.GetBool("server.loki_push_switch") {
		LokiLogger.Handle(lokiLabels, time.Now(), string(jsonPost))
	}

	if config.GetBool("logHub.pushSwitch") {
		projectName := config.GetString("logHub.projectName")
		logStoreName := config.GetString("logHub.logStoreName")

		logs := make([]*sls.Log, 0)
		contents := make([]*sls.LogContent, 0)

		contents = append(contents, &sls.LogContent{
			Key:   proto.String("env"),
			Value: proto.String(config.GetString("server.env")),
		})

		contents = append(contents, &sls.LogContent{
			Key:   proto.String("level"),
			Value: proto.String("INFO"),
		})

		contents = append(contents, &sls.LogContent{
			Key:   proto.String("tag"),
			Value: proto.String(tag),
		})

		contents = append(contents, &sls.LogContent{
			Key:   proto.String("content"),
			Value: proto.String(string(jsonPost)),
		})

		log := new(sls.Log)
		log.Time = proto.Uint32(uint32(time.Now().Unix()))
		log.Contents = contents

		logs = append(logs, log)
		logGroup := new(sls.LogGroup)
		logGroup.Logs = logs
		// 推送日志到阿里云日志服务器
		loghub.PutLogs(projectName, logStoreName, logGroup)
	}

	ConsoleLogger.WithTime(time.Now()).Info(jsonObject)
}
