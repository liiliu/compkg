package fiber_logger

import (
	"encoding/json"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gofiber/fiber/v2"
	"github.com/gogo/protobuf/proto"
	"github.com/tidwall/gjson"
	"net/url"
	"time"
	"weihu_server/library/common"
	"weihu_server/library/config"
	"weihu_server/library/logger"
	"weihu_server/library/loghub"
	"weihu_server/library/util"

	"github.com/google/uuid"

	"github.com/gookit/slog"
)

type Config struct {
	DefaultLevel     slog.Level
	ClientErrorLevel slog.Level
	ServerErrorLevel slog.Level

	WithRequestID bool
}

// New returns a fiber.Handler (middleware) that logs requests using slog.
//
// Requests with errors are logged using slog.Error().
// Requests without errors are logged using slog.Info().
func New() fiber.Handler {
	return NewWithConfig(Config{
		DefaultLevel:     slog.InfoLevel,
		ClientErrorLevel: slog.WarnLevel,
		ServerErrorLevel: slog.ErrorLevel,
		WithRequestID:    true,
	})
}

// NewWithConfig returns a fiber.Handler (middleware) that logs requests using slog.
func NewWithConfig(cnf Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Path()
		start := time.Now().UnixMilli()
		path := c.Path()

		requestID := uuid.New().String()
		if cnf.WithRequestID {
			c.Context().SetUserValue("request-id", requestID)
			c.Set("X-Request-ID", requestID)
		}

		err := c.Next()

		end := time.Now().UnixMilli()
		latency := end - start

		if string(c.Context().Method()) == "POST" {
			if logger.ConsoleLogger != nil {
				lokiLabels := map[string]string{
					"application": config.GetString("server.application"),
					"server":      config.GetString("server.name"),
					"env":         config.GetString("server.env"),
					"level":       "INFO",
					"tag":         common.LogTagInterfaceRequest,
				}

				reqBody := c.Body()
				params := make(map[string]interface{})
				if gjson.ValidBytes(reqBody) {
					_ = c.BodyParser(&params)
				} else {
					values, _ := url.ParseQuery(string(c.Body()))
					for k, v := range values {
						if len(v) > 0 {
							params[k] = v[0]
						}
					}
				}

				jsonObject := make(map[string]interface{})
				jsonObject["status"] = c.Response().StatusCode()
				jsonObject["method"] = string(c.Context().Method())
				jsonObject["path"] = path
				jsonObject["body"] = params
				respBody := c.Response().Body()
				if gjson.ValidBytes(respBody) {
					respBodyMap := util.ByteToMap(respBody)
					jsonObject["res"] = respBodyMap
				} else {
					jsonObject["res"] = string(respBody)
				}
				//jsonObject["res"] = string(c.Response().Body())
				//jsonObject["ip"] = c.Context().RemoteIP().String()
				jsonObject["ms"] = latency
				jsonObject["time"] = end

				for key, value := range lokiLabels {
					jsonObject[key] = value
				}
				jsonPost, _ := json.Marshal(jsonObject)

				if config.GetBool("server.loki_push_switch") {
					logger.LokiLogger.Handle(lokiLabels, time.Now(), string(jsonPost))
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
						Value: proto.String(common.LogTagInterfaceRequest),
					})

					contents = append(contents, &sls.LogContent{
						Key:   proto.String("path"),
						Value: proto.String(path),
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
			}
		}
		return err
	}
}

// GetRequestID returns the request identifier
func GetRequestID(c *fiber.Ctx) string {
	requestID, ok := c.Context().UserValue("request-id").(string)
	if !ok {
		return ""
	}

	return requestID
}
