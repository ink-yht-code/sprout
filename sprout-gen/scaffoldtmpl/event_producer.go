package scaffoldtmpl

// EventProducerTmpl Event 生产者模板
var EventProducerTmpl = `package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

// {{.NameUpper}}EventProducer {{.Name}} 事件生产者
type {{.NameUpper}}EventProducer struct {
	writer *kafka.Writer
}

// New{{.NameUpper}}EventProducer 创建事件生产者
func New{{.NameUpper}}EventProducer(writer *kafka.Writer) *{{.NameUpper}}EventProducer {
	return &{{.NameUpper}}EventProducer{writer: writer}
}

// Produce{{.NameUpper}}Created 发布{{.Name}}创建事件
func (p *{{.NameUpper}}EventProducer) Produce{{.NameUpper}}Created(ctx context.Context, {{.Name}}ID string) error {
	evt := {{.NameUpper}}Event{
		Type:      {{.NameUpper}}Created,
		{{.NameUpper}}ID:  {{.Name}}ID,
		Timestamp: time.Now().Unix(),
	}
	return p.produce(ctx, evt)
}

// Produce{{.NameUpper}}Updated 发布{{.Name}}更新事件
func (p *{{.NameUpper}}EventProducer) Produce{{.NameUpper}}Updated(ctx context.Context, {{.Name}}ID string) error {
	evt := {{.NameUpper}}Event{
		Type:      {{.NameUpper}}Updated,
		{{.NameUpper}}ID:  {{.Name}}ID,
		Timestamp: time.Now().Unix(),
	}
	return p.produce(ctx, evt)
}

// Produce{{.NameUpper}}Deleted 发布{{.Name}}删除事件
func (p *{{.NameUpper}}EventProducer) Produce{{.NameUpper}}Deleted(ctx context.Context, {{.Name}}ID string) error {
	evt := {{.NameUpper}}Event{
		Type:      {{.NameUpper}}Deleted,
		{{.NameUpper}}ID:  {{.Name}}ID,
		Timestamp: time.Now().Unix(),
	}
	return p.produce(ctx, evt)
}

func (p *{{.NameUpper}}EventProducer) produce(ctx context.Context, evt {{.NameUpper}}Event) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Topic: "{{.Name}}-events",
		Key:   []byte(evt.Type),
		Value: data,
	}

	return p.writer.WriteMessages(ctx, msg)
}
`
