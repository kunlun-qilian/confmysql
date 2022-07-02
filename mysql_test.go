package confmysql

import (
	"github.com/go-courier/sqlx/v2/migration"
	"github.com/kunlun-qilian/confmysql/tests"
	"testing"
)

func TestMySQL_Connect(t *testing.T) {
	tests.DB.Name = "example"
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
