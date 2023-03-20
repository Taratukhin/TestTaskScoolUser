package db

import (
	"context"

	"github.com/Taratukhin/TestTaskSchoolUser/core"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var _ core.Db = &Db{}

type Db struct {
	*sqlx.DB
}

func New(databaseConnect *sqlx.DB, requestsSQL []string) (db *Db, err error) {
	db = &Db{
		DB: databaseConnect,
	}

	ctx := context.Background()

	for _, sql := range requestsSQL {
		if err = db.exec(ctx, sql); err != nil {
			logrus.Errorf("%v in SQL: %s", err, sql)

			return nil, err
		}
	}

	return db, nil
}
