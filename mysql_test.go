package confmysql

import (
	"github.com/go-courier/sqlx/v2/migration"
	"github.com/kunlun-qilian/confmysql/v1/tests"
	"testing"
)

func TestMySQL_Connect(t *testing.T) {
	m := &MySQL{
		Host:     "127.0.0.1",
		User:     "root",
		Port:     33306,
		DBName:   "example",
		Password: "123456",
		Database: tests.DB,
	}

	m.SetDefaults()
	m.Init()

	if err := migration.Migrate(m.DB, nil); err != nil {
		panic(err)
	}
}
