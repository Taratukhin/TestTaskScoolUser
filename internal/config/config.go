package config

import (
	"flag"
	"strconv"
	"strings"
)

type Config struct {
	Port           int
	DBConnectURL   string
	DBInitSQLFiles []string
}

var (
	defaultPort         = 8080
	defaultDBConnectURL = "username:userpassword@tcp(127.0.0.1:5439)/schooldb"
	defaultInitSQLs     = []string{
		"./sql/create.sql",
		"./sql/scheme.sql",
		"./sql/truncate.sql",
		"./sql/data.sql",
		"./sql/constraints.sql",
	}
)

func New() *Config {
	currentConf := Config{}

	flag.IntVar(&currentConf.Port, "port", defaultPort, "Port of HTTP server "+strconv.Itoa(defaultPort))

	flag.StringVar(&currentConf.DBConnectURL, "db", defaultDBConnectURL, "Database connect URL")

	var sql arrayFlags

	flag.Var(&sql, "sql", "List of SQL files to execute, default: --sql "+strings.Join(defaultInitSQLs, " --sql "))

	flag.Parse()

	if len(sql) == 0 {
		sql = append(sql, defaultInitSQLs...)
	}

	return &currentConf
}

type arrayFlags []string // I can't find multiple keys in flag. But this project is too small to use Cobra

func (a *arrayFlags) String() string {
	s := strings.Join(*a, ",")

	return s
}

func (a *arrayFlags) Set(value string) error {
	*a = append(*a, value)

	return nil
}
