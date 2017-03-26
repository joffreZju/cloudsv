package model

import (
	"fmt"
	"s4s/common/lib/keycrypt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

const ReadOnly = 1

var (
	hasReadOnly    = false
	readOnlyDBName = "stowage-ro"
)

func NewOrm(readOnly ...int) orm.Ormer {
	o := orm.NewOrm()
	if hasReadOnly && len(readOnly) > 0 && readOnly[0] == ReadOnly {
		err := o.Using(readOnlyDBName)
		if err != nil {
			o = orm.NewOrm()
		}
	}
	return o
}

func InitMySQL(key string) (err error) {
	username := beego.AppConfig.String("mysql::username")
	password := beego.AppConfig.String("mysql::password")
	addr := beego.AppConfig.String("mysql::addr")
	addr_ro := beego.AppConfig.String("mysql::addr_ro")
	dbname := beego.AppConfig.String("mysql::dbname")

	if len(key) > 0 {
		password, err = keycrypt.Decode(key, password)
		if err != nil {
			return
		}
	}

	err = orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		return
	}
	err = orm.RegisterDataBase("default", "mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
			username, password, addr, dbname))
	if err != nil {
		return
	}
	if len(addr_ro) > 0 {
		err = orm.RegisterDataBase(readOnlyDBName, "mysql",
			fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
				username, password, addr_ro, dbname))
		if err != nil {
			return
		}
		hasReadOnly = true
	}

	orm.RegisterModel(new(User))
	err = orm.RunSyncdb("default", false, true)

	if beego.BConfig.RunMode != "prod" {
		orm.Debug = false
	}

	return
}
