package squirrel

import (
	"github.com/nats-io/nats.go"
	"time"
)

type EventBus struct {
	conn *nats.Conn
}

// NewEventBus 创建一个新的 EventBus 实例并连接到 NATS 服务器
func NewEventBus(url string) (*EventBus, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &EventBus{conn: nc}, nil
}

// Publish 发布消息到指定主题
func (e *EventBus) Publish(subject string, data []byte) error {
	return e.conn.Publish(subject, data)
}

// Subscribe 订阅指定主题，处理接收到的消息
func (e *EventBus) Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	return e.conn.Subscribe(subject, handler)
}

// Request 发送请求并等待响应
func (e *EventBus) Request(subject string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	return e.conn.Request(subject, data, timeout)
}

// Close 关闭与 NATS 服务器的连接
func (e *EventBus) Close() {
	e.conn.Close()
}
