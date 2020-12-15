package xorm

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tang-go/go-dog/plugins"
	"xorm.io/xorm"
)

//Mysql mysql
type Mysql struct {
	read  *xorm.Engine
	write *xorm.Engine
}

//NewMysql 初始化mysql
func NewMysql(cfg plugins.Cfg) *Mysql {
	mysql := new(Mysql)
	readurl := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local`,
		cfg.GetReadMysql().DbUser,
		cfg.GetReadMysql().DbPWd,
		cfg.GetReadMysql().DbIP,
		cfg.GetReadMysql().DbName)
	read, err := xorm.NewEngine("mysql", readurl)
	if err != nil {
		panic("connect to mysql error:" + err.Error())
	}
	//设置最大空闲连接数
	read.SetMaxIdleConns(cfg.GetReadMysql().MaxIdleConns)
	//设置数据库最大打开连接数
	read.SetMaxOpenConns(cfg.GetReadMysql().MaxOpenConns)
	//设置链接可重用时间
	read.SetConnMaxLifetime(time.Duration(cfg.GetReadMysql().ConnMaxLifetime) * time.Second)
	//设置日志
	read.ShowSQL(cfg.GetReadMysql().OpenLog)
	mysql.read = read
	//初始化写数据库
	writeurl := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local`,
		cfg.GetWriteMysql().DbUser,
		cfg.GetWriteMysql().DbPWd,
		cfg.GetWriteMysql().DbIP,
		cfg.GetWriteMysql().DbName)
	write, err := xorm.NewEngine("mysql", writeurl)
	if err != nil {
		panic("connect to mysql error:" + err.Error())
	}
	//设置最大空闲连接数
	write.SetMaxIdleConns(cfg.GetWriteMysql().MaxIdleConns)
	//设置数据库最大打开连接数
	write.SetMaxOpenConns(cfg.GetWriteMysql().MaxOpenConns)
	//设置链接可重用时间
	write.SetConnMaxLifetime(time.Duration(cfg.GetWriteMysql().ConnMaxLifetime) * time.Second)
	//设置日志
	write.ShowSQL(cfg.GetWriteMysql().OpenLog)
	//不为表增加s
	mysql.write = write
	return mysql
}

//GetReadEngine 获取读Mysql
func (m *Mysql) GetReadEngine() *xorm.Engine {
	return m.read
}

//GetWriteEngine 获取写Mysql
func (m *Mysql) GetWriteEngine() *xorm.Engine {
	return m.write
}
