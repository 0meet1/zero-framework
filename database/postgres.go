package database

import (
	"fmt"
	"time"

	"github.com/0meet1/zero-framework/global"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DATABASE_POSTGRES = "zero.database.postgres"
)

func InitPostgresDatabase() {
	dbURI := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s TimeZone=Asia/Shanghai",
		global.StringValue("zero.postgres.hostname"),
		global.IntValue("zero.postgres.hostport"),
		global.StringValue("zero.postgres.username"),
		global.StringValue("zero.postgres.dbname"),
		global.StringValue("zero.postgres.password"))
	dialector := postgres.New(postgres.Config{
		DSN:                  dbURI,
		PreferSimpleProtocol: false,
	})
	database, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dbPool, err := database.DB()
	if err != nil {
		panic(err)
	}

	dbPool.SetMaxIdleConns(global.IntValue("zero.postgres.maxIdleConns"))
	dbPool.SetMaxOpenConns(global.IntValue("zero.postgres.maxOpenConns"))
	dbPool.SetConnMaxLifetime(time.Second * time.Duration(global.IntValue("zero.postgres.maxLifetime")))

	global.Logger().Info(fmt.Sprintf(
		"postgres connect pool init success with %s, maxIdleConns: %d, maxOpenConns: %d, maxLifetime: %d",
		fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
			global.StringValue("zero.myspostgresql.hostname"),
			global.IntValue("zero.postgres.hostport"),
			global.StringValue("zero.postgres.username"),
			global.StringValue("zero.postgres.dbname")),
		global.IntValue("zero.postgres.maxIdleConns"),
		global.IntValue("zero.postgres.maxOpenConns"),
		global.IntValue("zero.postgres.maxLifetime")))
	dataSource := &GormDataSource{}
	dataSource.init(database)
	global.Key(DATABASE_POSTGRES, dataSource)
}

func InitCustomPostgresDatabase(registerName, prefix string) {

	dbURI := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s TimeZone=Asia/Shanghai",
		global.StringValue(fmt.Sprintf("%s.hostname", prefix)),
		global.IntValue(fmt.Sprintf("%s.hostport", prefix)),
		global.StringValue(fmt.Sprintf("%s.username", prefix)),
		global.StringValue(fmt.Sprintf("%s.dbname", prefix)),
		global.StringValue(fmt.Sprintf("%s.password", prefix)))
	dialector := postgres.New(postgres.Config{
		DSN:                  dbURI,
		PreferSimpleProtocol: false,
	})
	database, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dbPool, err := database.DB()
	if err != nil {
		panic(err)
	}

	dbPool.SetMaxIdleConns(global.IntValue(fmt.Sprintf("%s.maxIdleConns", prefix)))
	dbPool.SetMaxOpenConns(global.IntValue(fmt.Sprintf("%s.maxOpenConns", prefix)))
	dbPool.SetConnMaxLifetime(time.Second * time.Duration(global.IntValue(fmt.Sprintf("%s.maxLifetime", prefix))))

	global.Logger().Info(fmt.Sprintf(
		"postgres connect pool init success with %s, maxIdleConns: %d, maxOpenConns: %d, maxLifetime: %d",
		fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
			global.StringValue(fmt.Sprintf("%s.hostname", prefix)),
			global.IntValue(fmt.Sprintf("%s.hostport", prefix)),
			global.StringValue(fmt.Sprintf("%s.username", prefix)),
			global.StringValue(fmt.Sprintf("%s.dbname", prefix))),
		global.IntValue(fmt.Sprintf("%s.maxIdleConns", prefix)),
		global.IntValue(fmt.Sprintf("%s.maxOpenConns", prefix)),
		global.IntValue(fmt.Sprintf("%s.maxLifetime", prefix))))
	dataSource := &GormDataSource{}
	dataSource.init(database)
	global.Key(registerName, dataSource)
}
