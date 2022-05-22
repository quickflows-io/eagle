package model

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/go-eagle/eagle/pkg/config"
	"github.com/go-eagle/eagle/pkg/storage/orm"
)

// DB database global variables
var DB *gorm.DB

// Init Initialize the database
func Init() *gorm.DB {
	cfg, err := loadConf()
	if err != nil {
		panic(fmt.Sprintf("load orm conf err: %v", err))
	}

	DB = orm.NewMySQL(cfg)
	return DB
}

// GetDB Return to the default database
func GetDB() *gorm.DB {
	return DB
}

// loadConf load gorm config
func loadConf() (ret *orm.Config, err error) {
	var cfg orm.Config
	if err := config.Load("database", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
