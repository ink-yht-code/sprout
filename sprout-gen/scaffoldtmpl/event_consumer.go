package scaffoldtmpl

// EventConsumerTmpl 事件消费者模板
var EventConsumerTmpl = `package event

import (
	"context"
	"encoding/json"
)

type Consumer interface {
	Consume(ctx context.Context) ([]byte, error)
}

// {{.NameUpper}}EventConsumer {{.Name}} 事件消费者
type {{.NameUpper}}EventConsumer struct {
	consumer Consumer
	// TODO: 添加需要的服务依赖
}

// New{{.NameUpper}}EventConsumer 创建事件消费者
func New{{.NameUpper}}EventConsumer(consumer Consumer) *{{.NameUpper}}EventConsumer {
	return &{{.NameUpper}}EventConsumer{
		consumer: consumer,
	}
}

// Start 启动消费者
func (c *{{.NameUpper}}EventConsumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := c.consumer.Consume(ctx)
			if err != nil {
				continue
			}
			
			// 解析事件
			var evt {{.NameUpper}}Event
			if err := json.Unmarshal(msg, &evt); err != nil {
				continue
			}
			
			// 处理事件
			if err := c.handleEvent(ctx, evt); err != nil {
				// TODO: 错误处理
				continue
			}
		}
	}
}

// handleEvent 处理事件
func (c *{{.NameUpper}}EventConsumer) handleEvent(ctx context.Context, evt {{.NameUpper}}Event) error {
	// TODO: 实现事件处理逻辑
	return nil
}
`
