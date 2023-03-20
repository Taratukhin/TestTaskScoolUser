package db

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/Taratukhin/TestTaskSchoolUser/core"
	"github.com/Taratukhin/TestTaskSchoolUser/internal/config"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib" // needed for test
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

const (
	sqlCreateDB    = "../../sql/create.sql"
	sqlDropDB      = "../../sql/drop.sql"
	sqlScheme      = "../../sql/scheme.sql"
	sqlData        = "../../sql/data.sql"
	sqlConstraints = "../../sql/constraints.sql"
)

var initSQLFilesSequence = []string{sqlCreateDB, sqlDropDB, sqlScheme, sqlData, sqlConstraints}

type DBTestSuite struct {
	suite.Suite
	dbConnect *sqlx.DB
	ctx       context.Context
	db        *Db
}

func (suite *DBTestSuite) SetupSuite() {
	var err error

	suite.ctx = context.Background()

	conf := config.New()

	suite.dbConnect, err = sqlx.Connect("mysql", conf.DBConnectURL+"?multiStatements=true")
	suite.Require().NoErrorf(err, "error connection to DB, %w", err)

	initSQLs := make([]string, 0, len(initSQLFilesSequence))

	for _, sqlFileName := range initSQLFilesSequence {
		sql, err := os.ReadFile(sqlFileName)
		suite.Require().NoError(err)

		initSQLs = append(initSQLs, string(sql))
	}

	suite.db, err = New(suite.dbConnect, initSQLs)
	suite.Require().NoErrorf(err, "error create db, %w", err)
}

func (suite *DBTestSuite) TestSelectIDByApyKey() {
	id, err := suite.db.SelectIDByAPIKey(suite.ctx, "www-dfq92-sqfwf")
	suite.Require().NoError(err)
	suite.Require().Equal(uint(1), id)

	id, err = suite.db.SelectIDByAPIKey(suite.ctx, "not correct api key")
	suite.Require().ErrorIs(err, sql.ErrNoRows)
	suite.Require().Equal(uint(0), id)
}

func (suite *DBTestSuite) TestExec() {
	// checking the possibility of executing multiple SQL statements in one query
	suite.Require().NoError(suite.db.exec(suite.ctx, "SELECT 1; SELECT 2;"))
}

func (suite *DBTestSuite) TestSelectAboutUser() {
	info, err := suite.db.SelectAboutUser(suite.ctx, "test")
	suite.Require().NoError(err)
	suite.ElementsMatch(generateCorrectAboutUserTestAnswer(), info)

	info, err = suite.db.SelectAboutUser(suite.ctx, "unknown username")
	suite.Require().NoError(err)
	suite.Equal([]*core.AboutUser{}, info)

	info, err = suite.db.SelectAboutUser(suite.ctx, "")
	suite.Require().NoError(err)
	suite.Require().ElementsMatch(generateCorrectAboutUserAnswer(), info)
}

func TestDBSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}

func generateCorrectAboutUserTestAnswer() []*core.AboutUser {
	return []*core.AboutUser{
		{
			ID:        1,
			UserName:  "test",
			FirstName: "Александр",
			LastName:  "Школьный",
			City:      "Киев",
			School:    "гімназія №179 міста Києва",
		},
	}
}

func generateCorrectAboutUserAnswer() []*core.AboutUser {
	return []*core.AboutUser{
		{
			ID:        1,
			UserName:  "test",
			FirstName: "Александр",
			LastName:  "Школьный",
			City:      "Киев",
			School:    "гімназія №179 міста Києва",
		},
		{
			ID:        2,
			UserName:  "admin",
			FirstName: "Дмитрий",
			LastName:  "Арбузов",
			City:      "Харьков",
			School:    "ліцей №227",
		},
		{
			ID:        3,
			UserName:  "guest",
			FirstName: "Василий",
			LastName:  "Шпак",
			City:      "Житомир",
			School:    "Медична гімназія №33 міста Києва",
		},
	}
}
