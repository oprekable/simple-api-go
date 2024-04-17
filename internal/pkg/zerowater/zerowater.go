package zerowater

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/rs/zerolog"
)

//https://github.com/alexdrl/zerowater/tree/main

type ZerologLoggerAdapter struct {
	ctx    context.Context
	logger zerolog.Logger
}

// Logs an error message.
func (l *ZerologLoggerAdapter) Error(msg string, err error, fields watermill.LogFields) {
	event := l.logger.Err(err).Ctx(l.ctx)

	if fields != nil {
		addWatermillFieldsData(event, fields)
	}

	event.Msg(msg)
}

// Info Logs an info message.
func (l *ZerologLoggerAdapter) Info(msg string, fields watermill.LogFields) {
	event := l.logger.Info().Ctx(l.ctx)

	if fields != nil {
		addWatermillFieldsData(event, fields)
	}

	event.Msg(msg)
}

// Debug Logs a debug message.
func (l *ZerologLoggerAdapter) Debug(msg string, fields watermill.LogFields) {
	event := l.logger.Debug().Ctx(l.ctx)

	if fields != nil {
		addWatermillFieldsData(event, fields)
	}

	event.Msg(msg)
}

// Trace Logs a trace.
func (l *ZerologLoggerAdapter) Trace(msg string, fields watermill.LogFields) {
	event := l.logger.Trace().Ctx(l.ctx)

	if fields != nil {
		addWatermillFieldsData(event, fields)
	}

	event.Msg(msg)
}

// With Creates new adapter with the input fields as context.
func (l *ZerologLoggerAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	if fields == nil {
		return l
	}

	subLog := l.logger.With()

	for i, v := range fields {
		subLog = subLog.Interface(i, v)
	}

	return &ZerologLoggerAdapter{
		ctx:    l.ctx,
		logger: subLog.Logger(),
	}
}

// NewZerologLoggerAdapter Gets a new zerolog adapter for use in the watermill context.
func NewZerologLoggerAdapter(ctx context.Context, logger zerolog.Logger) *ZerologLoggerAdapter {
	return &ZerologLoggerAdapter{
		ctx:    ctx,
		logger: logger,
	}
}

func addWatermillFieldsData(event *zerolog.Event, fields watermill.LogFields) {
	for i, v := range fields {
		event.Interface(i, v)
	}
}
