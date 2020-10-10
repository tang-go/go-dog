package mysql

import (
	"fmt"
	"go-dog/plugins"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//Mysql 数据库
type Mysql struct {
	read  *gorm.DB
	write *gorm.DB
}

//NewMysql 初始化数据库
func NewMysql(cfg plugins.Cfg) *Mysql {
	mysql := new(Mysql)
	//初始化读数据库
	readurl := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local`,
		cfg.GetReadMysql().DbUser,
		cfg.GetReadMysql().DbPWd,
		cfg.GetReadMysql().DbIP,
		cfg.GetReadMysql().DbName)
	read, err := gorm.Open("mysql", readurl)
	if err != nil {
		panic("connect to mysql error:" + err.Error())
	}
	//设置最大空闲连接数
	read.DB().SetMaxIdleConns(cfg.GetReadMysql().MaxIdleConns)
	//设置数据库最大打开连接数
	read.DB().SetMaxOpenConns(cfg.GetReadMysql().MaxOpenConns)
	//设置链接可重用时间
	read.DB().SetConnMaxLifetime(time.Duration(cfg.GetReadMysql().ConnMaxLifetime) * time.Second)
	//设置日志
	read.LogMode(cfg.GetReadMysql().OpenLog)
	//不为表增加s
	read.SingularTable(true)
	//初始化写数据库
	writeurl := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local`,
		cfg.GetWriteMysql().DbUser,
		cfg.GetWriteMysql().DbPWd,
		cfg.GetWriteMysql().DbIP,
		cfg.GetWriteMysql().DbName)
	write, err := gorm.Open("mysql", writeurl)
	if err != nil {
		panic("connect to mysql error:" + err.Error())
	}
	//设置最大空闲连接数
	write.DB().SetMaxIdleConns(cfg.GetWriteMysql().MaxIdleConns)
	//设置数据库最大打开连接数
	write.DB().SetMaxOpenConns(cfg.GetWriteMysql().MaxOpenConns)
	//设置链接可重用时间
	write.DB().SetConnMaxLifetime(time.Duration(cfg.GetWriteMysql().ConnMaxLifetime) * time.Second)
	//设置日志
	write.LogMode(cfg.GetWriteMysql().OpenLog)
	//不为表增加s
	write.SingularTable(true)
	mysql.write = write
	mysql.read = read
	return mysql
}

//GetReadEngine 获取读Mysql
func (m *Mysql) GetReadEngine() *gorm.DB {
	return m.read
}

//GetWriteEngine 获取写Mysql
func (m *Mysql) GetWriteEngine() *gorm.DB {
	return m.write
}
