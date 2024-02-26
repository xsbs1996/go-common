package mysqldb

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MysqlRwConf struct {
	Name     string `json:"name" required:"true"`
	Ip       string `json:"ip" required:"true"`
	Port     string `json:"port" required:"true"`
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
	DB       string `json:"db" required:"true"`
}

var (
	onceRw  sync.Once
	dbMapRw sync.Map
)

// GetMysqlDBRw 获取默认数据库
func GetMysqlDBRw(c *MysqlRwConf) (*gorm.DB, error) {
	return GetDBRw(c)

}

func GetDBRw(dbInfo *MysqlRwConf) (*gorm.DB, error) {
	db, ok := dbMapRw.Load(dbInfo.Name)
	if !ok {
		return nil, NoneExistDB
	}
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		return nil, NoneExistDB
	}
	return gormDB, nil
}

func InitMysqlDBRw(dbInfoList []*MysqlRwConf, SlowThreshold time.Duration, MaxIdleConns int, MaxOpenConns int) {
	var err error
	onceRw.Do(func() {
		for _, dbInfo := range dbInfoList {
			// 数据库链接
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Asia%%2FShanghai",
				dbInfo.Username, dbInfo.Password, dbInfo.Ip, dbInfo.Port, dbInfo.DB)

			//日志
			dbLog := logger.New(
				logrus.New(),
				logger.Config{
					SlowThreshold:             SlowThreshold, // 慢 SQL 阈值
					LogLevel:                  logger.Info,   // 日志级别
					IgnoreRecordNotFoundError: false,         // 忽略ErrRecordNotFound（记录未找到）错误
					Colorful:                  false,         // 禁用彩色打印
				},
			)

			var database *gorm.DB
			database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: dbLog})
			if err != nil {
				logx.Errorf("Failed to open a db connection, db:%v, err:%v", dbInfo.DB, err)
				panic(err)
			}
			sqlDB, err := database.DB()
			if err != nil {
				logx.Errorf("Failed to get a db, db:%v, err:%v", dbInfo.DB, err)
				panic(err)
			}

			sqlDB.SetMaxIdleConns(MaxIdleConns)
			sqlDB.SetMaxOpenConns(MaxOpenConns)

			dbMapRw.Store(dbInfo.Name, database)
		}
	})

}
