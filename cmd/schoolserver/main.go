package main

import (
	"context"
	"os"

	"github.com/Taratukhin/TestTaskSchoolUser/cmd/schoolserver/cmd"
	"github.com/Taratukhin/TestTaskSchoolUser/internal/app"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/sirupsen/logrus"
)

func main() {
	app, err := app.New()
	if err != nil {
		logrus.Errorf("%v", err)
		os.Exit(1)
	}

	cmd.Execute(context.Background(), app)
}
