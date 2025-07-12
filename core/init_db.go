package core

import (
	"Blog_server/conf"
	"Blog_server/global"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

func InitDB() *gorm.DB {

	dc := global.Config.DB   //读库
	dc1 := global.Config.DB1 //写库

	db, err := gorm.Open(mysql.Open(conf.DB{}.DSN()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, //不生成外键约束
	})
	if err != nil {
		logrus.Errorf("数据库连接失败")
		return nil
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	logrus.Infof("数据库连接成功")

	if dc1.Empty() == false {
		//说明需要读写分离  这样我就设置
		//读写分离
		err := db.Use(dbresolver.Register(dbresolver.Config{
			// use `db2` as sources, `db3`, `db4` as replicas
			Sources:  []gorm.Dialector{mysql.Open(dc1.DSN())}, //写
			Replicas: []gorm.Dialector{mysql.Open(dc.DSN())},  //读
			// sources/replicas load balancing policy
			Policy: dbresolver.RandomPolicy{},
		}))
		if err != nil {
			logrus.Errorf("读写配置失败")
		}
	}

	return db
}
