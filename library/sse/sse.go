package sse

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"strings"
	"sync"
	"time"
	"weihu_server/api/view/desk"
	"weihu_server/library/cache"
	"weihu_server/library/common"
	"weihu_server/library/config"
	"weihu_server/library/logger"
	"weihu_server/library/util"
)

type SSEClient struct {
	ID         string
	Connection chan string
}

var (
	clients    = make(map[string]*SSEClient)
	clientsMux sync.Mutex
)

// Add a new client
func addClient(id string) *SSEClient {
	client := &SSEClient{
		ID:         id,
		Connection: make(chan string),
	}
	clientsMux.Lock()
	clients[id] = client
	clientsMux.Unlock()
	return client
}

// Remove a client
func removeClient(id string) {
	clientsMux.Lock()
	if client, exists := clients[id]; exists {
		close(client.Connection)
		delete(clients, id)
	}
	clientsMux.Unlock()
}

// Broadcast a message to a specific client
func sendToClient(id string, message string) {
	clientsMux.Lock()
	if client, exists := clients[id]; exists {
		client.Connection <- message
	}
	clientsMux.Unlock()
}

func Middleware(c *fiber.Ctx) error {
	id := c.Params("id")
	client := addClient(id)
	fmt.Println("sse id:", id)

	//go func() {
	//	for {
	//		time.Sleep(5 * time.Second)
	//		sendToClient(id, fmt.Sprintf("Hello %s! Time: %s", id, time.Now().Format(time.RFC3339)))
	//	}
	//}()
	go GetSseMessage()

	// Set the proper headers for SSE
	c.Set("Content-Type", "text/event-stream")
	c.Set("X-Accel-Buffering", "no")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")
	c.Set("Access-Control-Allow-Origin", "*")

	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		ticker := time.NewTicker(15 * time.Second) // 每 15 秒发送一次心跳包
		defer ticker.Stop()

		for {
			select {
			case msg := <-client.Connection:
				fmt.Fprintf(w, "data: Message: %s\n\n", msg)
				fmt.Println(msg)
				err := w.Flush()
				if err != nil {
					fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
					removeClient(id)
					return
				}
			case <-ticker.C: // 定时发送心跳包
				// 发送心跳包（保持连接活跃）
				fmt.Fprintf(w, ": keep-alive\n\n")
				err := w.Flush()
				if err != nil {
					fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
					removeClient(id)
					return
				}
			}
		}
	}))

	return nil
}

// SetSseMessage 设置sse消息
func SetSseMessage(id, msg string) {
	queueName := strings.ToUpper(fmt.Sprintf("%s_%s_%s", common.QueueSSE, config.GetString("server.name"), config.GetString("server.env")))
	sseDataStruct := desk.SSEData{
		ID:      id,
		Message: msg,
	}
	cache.Rpush(queueName, util.JsonToString(sseDataStruct))
}

// GetSseMessage 获取sse消息
func GetSseMessage() {
	for {
		queueName := strings.ToUpper(fmt.Sprintf("%s_%s_%s", common.QueueSSE, config.GetString("server.name"), config.GetString("server.env")))
		count, err := cache.Llen(queueName)
		if err != nil {
			logger.Error(common.LogSSE, "GetSseData error: %v", err)
			return
		}

		for i := 0; i < int(count); i++ {
			reply, err := cache.Lpop(queueName)
			if err != nil {
				logger.Error(common.LogSSE, "GetSseData error: %v", err)
				break
			}

			sseDataStruct := new(desk.SSEData)
			err = json.Unmarshal(reply, &sseDataStruct)
			if err != nil {
				logger.Error(common.LogSSE, "GetSseData error: %v", err)
				continue
			}
			logger.Info(common.LogSSE, "GetSseData: %+v", sseDataStruct)

			sendToClient(sseDataStruct.ID, sseDataStruct.Message)
		}
		time.Sleep(1 * time.Second)
	}
}
