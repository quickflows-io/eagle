package orm

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	otelgorm "github.com/1024casts/gorm-opentelemetry"

	// MySQL driver.
	"gorm.io/driver/mysql"
	// GORM MySQL
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config mysql config
type Config struct {
	Name            string
	Addr            string
	UserName        string
	Password        string
	ShowLog         bool
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxLifeTime time.Duration
	SlowThreshold   time.Duration // Slow query duration, default 500ms
}

// NewMySQL Slow query duration, default 500ms
func NewMySQL(c *Config) (db *gorm.DB) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		c.UserName,
		c.Password,
		c.Addr,
		c.Name,
		true,
		//"Asia/Shanghai"),
		"Local")

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panicf("open mysql failed. database name: %s, err: %+v", c.Name, err)
	}
	// set for db connection
	// It is used to set the maximum number of open connections. The default value is 0, which means no limit.
	//Setting the maximum number of connections can avoid the error of too many connections when connecting to mysql
	//due to too high concurrency.
	sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	// It is used to set the number of idle connections. When the number of idle connections is set, when an open
	//connection is used, it can be placed in the pool for the next use.
	sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(c.ConnMaxLifeTime)

	db, err = gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), gormConfig(c))
	if err != nil {
		log.Panicf("database connection failed. database name: %s, err: %+v", c.Name, err)
	}
	db.Set("gorm:table_options", "CHARSET=utf8mb4")

	// Initialize otel plugin with options
	plugin := otelgorm.NewPlugin(
	// include any options here
	)

	// set trace
	err = db.Use(plugin)
	if err != nil {
		log.Panicf("using gorm opentelemetry, err: %+v", err)
	}

	return db
}

// gormConfig Decide whether to enable logging according to the configuration
func gormConfig(c *Config) *gorm.Config {
	config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true} // 禁止外键约束, 生产环境不建议使用外键约束
	// print all SQL
	if c.ShowLog {
		config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}
	// Only print slow queries
	if c.SlowThreshold > 0 {
		config.Logger = logger.New(
			//get stdout as Writer
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				//Set the slow query time threshold
				SlowThreshold: c.SlowThreshold, // nolint: golint
				Colorful:      true,
				//Set the log level, only the slow query log will be output above the specified level
				LogLevel: logger.Warn,
			},
		)
	}
	return config
}
