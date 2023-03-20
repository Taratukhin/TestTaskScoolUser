package app

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Taratukhin/TestTaskSchoolUser/core"
	"github.com/Taratukhin/TestTaskSchoolUser/internal/api"
	"github.com/Taratukhin/TestTaskSchoolUser/internal/config"
	"github.com/Taratukhin/TestTaskSchoolUser/internal/db"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type App struct {
	*sqlx.DB
	conf       *config.Config
	db         core.Db
	httpServer *http.Server
}

func New() (*App, error) {
	conf := config.New()

	databaseConnect, err := sqlx.Connect("mysql", conf.DBConnectURL+"?multiStatements=true")
	if err != nil {
		return nil, err
	}

	initSQLs := make([]string, 0, len(conf.DBInitSQLFiles))

	for _, sqlFileName := range conf.DBInitSQLFiles {
		sql, err := os.ReadFile(sqlFileName)
		if err != nil {
			return nil, err
		}

		initSQLs = append(initSQLs, string(sql))
	}

	db, err := db.New(databaseConnect, initSQLs)
	if err != nil {
		return nil, err
	}

	return &App{
		DB:   databaseConnect,
		conf: conf,
		db:   db,
		httpServer: &http.Server{
			Addr:    ":" + strconv.Itoa(conf.Port),
			Handler: api.NewMux(db),
		},
	}, nil
}

func (app *App) Run(ctx context.Context) error {
	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-gctx.Done()

		waitContext, waitCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer waitCancel()

		err := app.httpServer.Shutdown(waitContext)
		if err != nil {
			logrus.Errorf("graceful http server shutdown error:%v", err)

			return err
		}

		logrus.Info("graceful shutdown http server")

		return nil
	})

	pid := os.Getpid()

	g.Go(func() error {
		logrus.Infof("schoolserver (Pid=%v) started at http://localhost:%d", pid, app.conf.Port)
		err := app.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logrus.Errorf("error start http server:%v", err)
		}

		return err
	})

	err := g.Wait()

	logrus.Infof("schoolserver stopped")

	_ = app.DB.Close()

	return err
}
