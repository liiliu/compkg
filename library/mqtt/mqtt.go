package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"weihu_server/library/common"
	"weihu_server/library/config"
	"weihu_server/library/logger"
	"weihu_server/library/util"
)

// ClientWrapper 封装 MQTT 客户端
type ClientWrapper struct {
	client mqtt.Client
}

type SendMsgTask struct {
	Topic string      `json:"topic"`
	Msg   interface{} `json:"msg"`
}

var mqttClient *ClientWrapper

// Initial 初始化 MQTT 客户端
func Initial() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", config.GetString("mqtt.broker"), config.GetString("mqtt.port")))
	opts.SetClientID(config.GetString("mqtt.client"))
	opts.SetUsername(config.GetString("mqtt.username"))
	opts.SetPassword(config.GetString("mqtt.password"))
	opts.SetDefaultPublishHandler(messageHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.AutoReconnect = true
	opts.CleanSession = false
	// 自动重连
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		logger.Error(common.LogTagMqtt, fmt.Sprintf("Connection lost: %v", err))
	}

	// 创建客户端
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.Error(common.LogTagMqtt, token.Error().Error())
		return
	}

	mqttClient = &ClientWrapper{client: client}

	logger.Info(common.LogTagMqtt, "MQTT client initialized")

	//发送测试消息
	Publish(fmt.Sprintf("%s%s", config.GetString("mqtt.topicPrefix"), config.GetString("mqtt.client")), "Hello, MQTT!", 0, false)
}

// 连接回调
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
	client.Subscribe(fmt.Sprintf("%s%s", config.GetString("mqtt.topicPrefix"), config.GetString("mqtt.client")), 0, messagePubHandler)
}

// 连接丢失回调
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
	//fmt.Println(err.Error())
}

// 订阅消息处理
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Println(msg.Topic())
	payload := string(msg.Payload())
	log.Println(payload)
}

// messageHandler 是一个默认的没有订阅回调消息处理函数
var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

// Subscribe 订阅主题
func Subscribe(topic string, qos byte) error {
	token := mqttClient.client.Subscribe(topic, qos, nil)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Printf("Subscribed to topic: %s", topic)
	return nil
}

// Publish 发送消息到指定主题 qos:服务质量 0-最多一次 1-至少一次 2-仅一次，retained:是否保留消息
func Publish(topic string, payload interface{}, qos byte, retained bool) error {
	msgTask := new(SendMsgTask)
	msgTask.Topic = topic
	msgTask.Msg = payload
	content := util.JsonToString(msgTask)

	token := mqttClient.client.Publish(topic, qos, retained, content)
	token.Wait() // 等待消息发送
	log.Printf("Message published to topic: %s", topic)
	return nil
}

// Disconnect 断开客户端连接
func Disconnect() {
	mqttClient.client.Disconnect(250) // 等待250毫秒，以确保消息被发送
	log.Println("Client disconnected")
}
