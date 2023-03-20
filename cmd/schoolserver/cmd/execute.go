package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Taratukhin/TestTaskSchoolUser/internal/app"
	"github.com/sirupsen/logrus"
)

func Execute(ctx context.Context, app *app.App) {
	ctx, done := context.WithCancel(ctx)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-signalChannel:
			logrus.Infof("received signal: %s\n", sig)
			done()
		case <-ctx.Done():
		}
	}()

	if err := app.Run(ctx); err != nil {
		signal.Stop(signalChannel)
		os.Exit(1)
	}

	signal.Stop(signalChannel)
	done()
}
