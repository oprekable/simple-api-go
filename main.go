package main

import (
	"context"
	"embed"
	"errors"
	"os"
	"simple-api-go/cmd"
	"simple-api-go/internal/pkg/shutdown"
	"simple-api-go/internal/pkg/utils/log"
	"simple-api-go/variable"
	"time"

	"golang.org/x/sync/errgroup"
)

//go:embed all:embeds
var embedFS embed.FS

func init() {
	_ = os.Setenv(variable.TZ, variable.TimeZone)
}

func main() {
	loc, _ := time.LoadLocation(variable.TimeZone)
	_, offset1 := time.Now().Zone()
	_, offset2 := time.Now().In(loc).Zone()

	if offset1 != offset2 {
		time.Local = loc
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, log.StartTime, time.Now())

	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)
	sigTrap := shutdown.TermSignalTrap()

	cmd.Execute(
		ctx,
		cancel,
		eg,
		&embedFS,
		variable.TimeZone,
		variable.TimeFormatPostgresString,
		variable.TimePostgresFriendlyFormat,
	)

	eg.Go(func() error {
		return sigTrap.Wait(ctx)
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, shutdown.ErrTermSig) {
		log.MsgOrPanic(context.Background(), "graceful shutdown successfully finished", err)
	}
}
