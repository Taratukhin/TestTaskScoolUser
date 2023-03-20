package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Taratukhin/TestTaskSchoolUser/core"
	"github.com/Taratukhin/TestTaskSchoolUser/internal/config"
	"github.com/Taratukhin/TestTaskSchoolUser/internal/db"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib" // needed for test
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
	dbConnect *sqlx.DB
	ctx       context.Context
	db        *db.Db
	srv       *httptest.Server
}

func (suite *ServerTestSuite) SetupSuite() {
	var err error

	suite.ctx = context.Background()

	conf := config.New()

	suite.dbConnect, err = sqlx.Connect("mysql", conf.DBConnectURL+"?multiStatements=true")
	suite.Require().NoErrorf(err, "error connection to DB, %w", err)

	initSQLs := make([]string, 0, len(conf.DBInitSQLFiles))

	for _, sqlFileName := range conf.DBInitSQLFiles {
		sql, err := os.ReadFile(sqlFileName)
		suite.Require().NoError(err)

		initSQLs = append(initSQLs, string(sql))
	}

	suite.db, err = db.New(suite.dbConnect, initSQLs)
	suite.Require().NoErrorf(err, "error create db, %w", err)

	suite.srv = httptest.NewServer(NewMux(suite.db))
}

func (suite *ServerTestSuite) TearDownTest() {
	suite.srv.Close()
}

func (suite *ServerTestSuite) TestRoutingGetUserProfile() {
	req, err := http.NewRequestWithContext(suite.ctx, http.MethodGet, suite.srv.URL+pathProfile, nil)
	suite.Require().NoError(err)

	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)
	suite.Require().Equal(http.StatusForbidden, res.StatusCode)
	_, _ = io.ReadAll(res.Body)
	res.Body.Close()

	req.Header.Set("Api-key", "www-dfq92-sqfwf")
	res, err = http.DefaultClient.Do(req)
	suite.Require().NoError(err)
	suite.Require().Equal(http.StatusOK, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	suite.Require().NoError(err)
	res.Body.Close()

	var myAnswers []*core.AboutUser

	suite.Require().NoError(json.Unmarshal(body, &myAnswers))
	suite.Require().ElementsMatch(generateCorrectAboutUserAnswer(), myAnswers)

	req, err = http.NewRequestWithContext(suite.ctx, http.MethodGet, suite.srv.URL+pathProfile+"?username=test", nil)
	suite.Require().NoError(err)

	req.Header.Set("Api-key", "www-dfq92-sqfwf")
	res, err = http.DefaultClient.Do(req)
	suite.Require().NoError(err)
	suite.Require().Equal(http.StatusOK, res.StatusCode)

	body, err = io.ReadAll(res.Body)
	suite.Require().NoError(err)
	res.Body.Close()

	var aboutUserTest *core.AboutUser

	suite.Require().NoError(json.Unmarshal(body, &aboutUserTest))
	suite.Require().Equal(generateCorrectAboutUserTestAnswer(), aboutUserTest)
}

func TestDBSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func generateCorrectAboutUserTestAnswer() *core.AboutUser {
	return &core.AboutUser{
		ID:        1,
		UserName:  "test",
		FirstName: "Александр",
		LastName:  "Школьный",
		City:      "Киев",
		School:    "гімназія №179 міста Києва",
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
