package service

import (
	"github.com/intelligentfish/dcn/app"
	"github.com/intelligentfish/dcn/config"
	"github.com/intelligentfish/dcn/define"
	"github.com/intelligentfish/dcn/log"
	"github.com/intelligentfish/dcn/serviceGroup"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
	"time"
)

var (
	dbSrvInst *DBSrv    // db service instance
	dbSrvOnce sync.Once // db service once
)

// DB
type DBSrv struct {
	*BaseSrv
	db *gorm.DB
}

// NewDBSrv factory method
func NewDBSrv() *DBSrv {
	object := &DBSrv{
		BaseSrv: NewSrvBase("DBSrv",
			define.ServiceTypeDB,
			define.StartupPriorityDB,
			define.ShutdownPriorityDB),
	}
	gorm.DefaultTableNameHandler = func(_ *gorm.DB, defaultTableName string) string {
		return "tb_" + defaultTableName
	}
	return object
}

// Start start the service
func (object *DBSrv) Start() (err error) {
	if err = object.BaseSrv.Start(); nil != err {
		return
	}
	for i := 0; i < config.Inst().GetDBRetryTimes() && !app.Inst().IsStopped(); i++ {
		if 0 < len(config.Inst().GetMySQLConnection()) {
			object.db, err = gorm.Open("mysql", config.Inst().GetMySQLConnection())
			if nil != err {
				log.Inst().Error(err.Error())
			} else {
				//model.Init(DBSrvInst().db)
			}
		}
		if nil != err {
			log.Inst().Error(err.Error())
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	return
}

// Stop stop the service
func (object *DBSrv) Stop() {
	var err error
	object.BaseSrv.Stop()
	if nil != object.db {
		if err = object.db.Close(); nil != err {
			log.Inst().Error(err.Error())
		}
	}
}

// WithDB use db resource
func (object *DBSrv) WithDB(callback func(db *gorm.DB)) (err error) {
	callback(object.db)
	err = object.db.Error
	return
}

// WithDBTx use db resource transaction
func (object *DBSrv) WithDBTx(callback func(db *gorm.DB) (err error)) (err error) {
	db := object.db.Begin()
	err = callback(db)
	if nil != err {
		log.Inst().Error(err.Error())
		err = db.Rollback().Error
	} else {
		err = db.Commit().Error
	}
	return
}

// Inst singleton
func DBSrvInst() *DBSrv {
	dbSrvOnce.Do(func() {
		dbSrvInst = NewDBSrv()
	})
	return dbSrvInst
}

// init method
func init() {
	serviceGroup.Inst().AddSrv(DBSrvInst())
}
