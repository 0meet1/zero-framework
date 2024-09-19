package database

import (
	"fmt"
	"time"

	"github.com/0meet1/zero-framework/global"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DATABASE_MYSQL = "zero.database.mysql"
)

func InitMYSQLDatabase() {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
		global.StringValue("zero.mysql.username"),
		global.StringValue("zero.mysql.password"),
		global.StringValue("zero.mysql.hostname"),
		global.IntValue("zero.mysql.hostport"),
		global.StringValue("zero.mysql.dbname"))
	dialector := mysql.New(mysql.Config{
		DSN:                       dbURI,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	})
	database, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dbPool, err := database.DB()
	if err != nil {
		panic(err)
	}

	dbPool.SetMaxIdleConns(global.IntValue("zero.mysql.maxIdleConns"))
	dbPool.SetMaxOpenConns(global.IntValue("zero.mysql.maxOpenConns"))
	dbPool.SetConnMaxLifetime(time.Second * time.Duration(global.IntValue("zero.mysql.maxLifetime")))

	global.Logger().Info(fmt.Sprintf(
		"mysql connect pool init success with %s, maxIdleConns: %d, maxOpenConns: %d, maxLifetime: %d",
		fmt.Sprintf("username:password@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
			global.StringValue("zero.mysql.hostname"),
			global.IntValue("zero.mysql.hostport"),
			global.StringValue("zero.mysql.dbname")),
		global.IntValue("zero.mysql.maxIdleConns"),
		global.IntValue("zero.mysql.maxOpenConns"),
		global.IntValue("zero.mysql.maxLifetime")))
	dataSource := &GormDataSource{}
	dataSource.init(database)
	global.Key(DATABASE_MYSQL, dataSource)
}

func InitCustomMysqlDatabase(registerName, prefix string) {

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
		global.StringValue(fmt.Sprintf("%s.username", prefix)),
		global.StringValue(fmt.Sprintf("%s.password", prefix)),
		global.StringValue(fmt.Sprintf("%s.hostname", prefix)),
		global.IntValue(fmt.Sprintf("%s.hostport", prefix)),
		global.StringValue(fmt.Sprintf("%s.dbname", prefix)))
	dialector := mysql.New(mysql.Config{
		DSN:                       dbURI,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
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
		"mysql connect pool init success with %s, maxIdleConns: %d, maxOpenConns: %d, maxLifetime: %d",
		fmt.Sprintf("username:password@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
			global.StringValue(fmt.Sprintf("%s.hostname", prefix)),
			global.IntValue(fmt.Sprintf("%s.hostport", prefix)),
			global.StringValue(fmt.Sprintf("%s.dbname", prefix))),
		global.IntValue(fmt.Sprintf("%s.maxIdleConns", prefix)),
		global.IntValue(fmt.Sprintf("%s.maxOpenConns", prefix)),
		global.IntValue(fmt.Sprintf("%s.maxLifetime", prefix))))
	dataSource := &GormDataSource{}
	dataSource.init(database)
	global.Key(registerName, dataSource)
}
