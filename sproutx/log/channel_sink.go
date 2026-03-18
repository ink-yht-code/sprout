package log

import "go.uber.org/zap/zapcore"

type Event struct {
	Entry  zapcore.Entry
	Fields []zapcore.Field
}

type ChannelSink struct {
	ch chan Event
}

func NewChannelSink(buffer int) *ChannelSink {
	if buffer <= 0 {
		buffer = 1024
	}
	return &ChannelSink{ch: make(chan Event, buffer)}
}

func (s *ChannelSink) C() <-chan Event {
	return s.ch
}

func (s *ChannelSink) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	select {
	case s.ch <- Event{Entry: entry, Fields: fields}:
	default:
	}
	return nil
}

func (s *ChannelSink) Sync() error {
	return nil
}
