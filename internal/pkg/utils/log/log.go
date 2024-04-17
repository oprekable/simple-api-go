package log

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type KeyType string

const (
	// StartTime ...
	StartTime KeyType = "start_time"
)

type UptimeHook struct{}

func (u UptimeHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	startTime := ctx.Value(StartTime)
	uptime := ""

	if startTime != nil {
		uptime = time.Since(startTime.(time.Time)).String()
	}

	e.Str("uptime", uptime)
}

type Zero struct {
	LogCtx context.Context
	LogZ   zerolog.Logger
}

func (l *Zero) Printf(ctx context.Context, format string, v ...interface{}) {
	startTime := ctx.Value(StartTime)
	uptime := ""

	if startTime != nil {
		uptime = time.Since(startTime.(time.Time)).String()
	}

	ctx = l.LogZ.WithContext(l.LogCtx)
	zerolog.Ctx(ctx).
		Info().
		Str("uptime", uptime).
		Msg(fmt.Sprintf(format, v...))
}

type Caller struct {
	File     string
	Function string
	Line     int
}

func New() Caller {
	pc, f, l, _ := runtime.Caller(4)
	return Caller{
		File:     f,
		Line:     l,
		Function: runtime.FuncForPC(pc).Name(),
	}
}

func (lc Caller) MarshalZerologObject(e *zerolog.Event) {
	e.Str("file", lc.File).
		Int("line", lc.Line).
		Str("function", lc.Function)
}

func AddStrWithKey(ctx context.Context, key, msg string) {
	zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str(key, msg)
	})
}

func AddStr(ctx context.Context, msg string) {
	zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str(strconv.FormatInt(time.Now().UnixNano(), 10), msg)
	})
}

func AddErr(ctx context.Context, err error) {
	if err == nil {
		return
	}

	zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Object(strconv.FormatInt(time.Now().UnixNano(), 10), New()).
			AnErr(strconv.FormatInt(time.Now().UnixNano(), 10), err)
	})
}

func AddErrAndStr(ctx context.Context, msg string, err error) {
	zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Object(strconv.FormatInt(time.Now().UnixNano(), 10), New()).
			AnErr(strconv.FormatInt(time.Now().UnixNano(), 10), err).
			Str(strconv.FormatInt(time.Now().UnixNano(), 10), msg)
	})
}

func AddStrOrPanic(ctx context.Context, err error, strFatal string, strInfo string) {
	if err != nil {
		zerolog.Ctx(ctx).
			Panic().
			Err(err).
			Ctx(ctx).
			Caller(2).
			Msg(strFatal)

		return
	}

	zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str(strconv.FormatInt(time.Now().UnixNano(), 10), strInfo)
	})
}

func AddStrOrAddErr(ctx context.Context, err error, strError string, strInfo string) {
	if err != nil {
		zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Object(strconv.FormatInt(time.Now().UnixNano(), 10), New()).
				AnErr(strconv.FormatInt(time.Now().UnixNano(), 10), err).
				Str(strconv.FormatInt(time.Now().UnixNano(), 10), strError)
		})

		return
	}

	zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str(strconv.FormatInt(time.Now().UnixNano(), 10), strInfo)
	})
}

func Panic(ctx context.Context, msg string, err error) {
	zerolog.Ctx(ctx).
		Panic().
		Err(err).
		Ctx(ctx).
		Caller(2).
		Msg(msg)
}

func MsgOrPanic(ctx context.Context, msg string, err error) {
	if err != nil {
		zerolog.Ctx(ctx).
			Panic().
			Err(err).
			Ctx(ctx).
			Caller(2).
			Msg(msg)

		return
	}

	zerolog.Ctx(ctx).
		Info().
		Ctx(ctx).
		Msg(msg)
}

func MsgOrErr(ctx context.Context, msg string, err error) {
	if err != nil {
		zerolog.Ctx(ctx).
			Err(err).
			Ctx(ctx).
			Caller(2).
			Msg(msg)

		return
	}

	zerolog.Ctx(ctx).
		Info().
		Ctx(ctx).
		Msg(msg)
}

func Msg(ctx context.Context, msg string) {
	zerolog.Ctx(ctx).
		Info().
		Ctx(ctx).
		Msg(msg)
}

func Err(ctx context.Context, msg string, err error) {
	zerolog.Ctx(ctx).
		Err(err).
		Ctx(ctx).
		Caller(2).
		Msg(msg)
}
