package confmysql

import (
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type AutoMigrateConfig struct {
	// models
	Models []interface{}
	// generate model path
	ModelPath string
	// generate query path
	QueryPath string
}

type MySQL struct {
	DSN               string `env:""`
	AutoMigrateConfig *AutoMigrateConfig

	db *gorm.DB

	commands []*cobra.Command
}

func (c *MySQL) Init() {

	conf := &gorm.Config{}
	conf.NamingStrategy = schema.NamingStrategy{
		TablePrefix:   "t_",
		SingularTable: true,
	}

	db, err := gorm.Open(mysql.Open(c.DSN), conf)
	if err != nil {
		panic(err)
	}
	c.db = db

	if c.AutoMigrateConfig != nil {
		c.commands = make([]*cobra.Command, 0)
		// migrate model to db
		c.commands = append(c.commands, &cobra.Command{
			Use: "migrate",
			Run: func(cmd *cobra.Command, args []string) {
				if err := c.db.AutoMigrate(c.AutoMigrateConfig.Models...); err != nil {
					panic(err)
				}
			},
		})

		// generate model from db table
		c.commands = append(c.commands, &cobra.Command{
			Use: "gen-model",
			Run: func(cmd *cobra.Command, args []string) {
				g := gen.NewGenerator(gen.Config{
					OutPath: c.AutoMigrateConfig.ModelPath,
				})
				g.UseDB(c.db)
				g.GenerateAllTable()
				g.Execute()
			},
		})
		// generate query from model
		c.commands = append(c.commands, &cobra.Command{
			Use: "gen-query",
			Run: func(cmd *cobra.Command, args []string) {
				g := gen.NewGenerator(gen.Config{
					OutPath: c.AutoMigrateConfig.QueryPath,
				})
				g.UseDB(c.db)
				g.ApplyBasic(c.AutoMigrateConfig.Models...)
				g.Execute()
			},
		})
	}
}

func (c *MySQL) Commands() []*cobra.Command {
	return c.commands
}

func (c *MySQL) DB() *gorm.DB {
	if c.db == nil {
		panic("get db before init")
	}
	return c.db
}
