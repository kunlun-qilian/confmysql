package confmysql

import (
	"fmt"
	"github.com/go-courier/sqlx/v2"
	"github.com/go-courier/sqlx/v2/migration"
	"net/url"
	"time"

	"github.com/go-courier/envconf"
	"github.com/go-courier/sqlx/v2/mysqlconnector"

	"github.com/spf13/cobra"
)

type MySQL struct {
	// DBName
	DBName          string           `env:""`
	Host            string           `env:""`
	Port            int              `env:""`
	User            string           `env:""`
	Password        envconf.Password `env:""`
	Extra           string           `env:""`
	PoolSize        int              `env:""`
	ConnMaxLifetime envconf.Duration
	Retry
	Database *sqlx.Database `env:"-"`
	*sqlx.DB `env:"-"`

	commands []*cobra.Command
}

func (m *MySQL) SetDefaults() {
	if m.Port == 0 {
		m.Port = 3306
	}

	if m.PoolSize == 0 {
		m.PoolSize = 10
	}

	if m.ConnMaxLifetime == 0 {
		m.ConnMaxLifetime = envconf.Duration(1 * time.Hour)
	}

	if m.Extra == "" {
		values := url.Values{}
		values.Set("charset", "utf8mb4")
		values.Set("parseTime", "true")
		values.Set("interpolateParams", "true")
		values.Set("autocommit", "true")
		values.Set("loc", "Local")
		m.Extra = values.Encode()
	}
}

func (m *MySQL) URL() string {
	password := m.Password
	if password != "" {
		password = ":" + password
	}
	return fmt.Sprintf("%s%s@tcp(%s:%d)", m.User, password, m.Host, m.Port)
}

func (m *MySQL) Connect() error {
	m.Database.Name = m.DBName
	m.SetDefaults()
	db := m.Database.OpenDB(&mysqlconnector.MysqlConnector{
		Host:  m.URL(),
		Extra: m.Extra,
	})
	db.SetMaxOpenConns(m.PoolSize)
	db.SetMaxIdleConns(m.PoolSize / 2)
	db.SetConnMaxLifetime(time.Duration(m.ConnMaxLifetime))
	m.DB = db
	return nil
}

func (m *MySQL) Init() {

	// migrate
	m.commands = append(m.commands, &cobra.Command{
		Use: "migrate",
		Run: func(cmd *cobra.Command, args []string) {
			if err := migration.Migrate(m.DB, nil); err != nil {
				panic(err)
			}
		},
	})

	if m.DB == nil {
		m.Do(m.Connect)
	}
}

func (m *MySQL) Get() *sqlx.DB {
	if m.DB == nil {
		panic(fmt.Errorf("get db before init"))
	}
	return m.DB
}

func (m *MySQL) Commands() []*cobra.Command {
	return m.commands
}
