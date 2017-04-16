package model

import (
	"common/lib/keycrypt"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
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

func InitPgSQL(key string) (err error) {
	username := beego.AppConfig.String("pgsql::username")
	password := beego.AppConfig.String("pgsql::password")
	addr := beego.AppConfig.String("pgsql::addr")
	port := beego.AppConfig.String("pgsql::port")
	addr_ro := beego.AppConfig.String("pgsql::addr_ro")
	dbname := beego.AppConfig.String("pgsql::dbname")

	if len(key) > 0 {
		password, err = keycrypt.Decode(key, password)
		if err != nil {
			return
		}
	}

	err = orm.RegisterDriver("postgres", orm.DRPostgres)
	if err != nil {
		return
	}
	beego.Debug(username, password, addr, dbname)
	err = orm.RegisterDataBase("default", "postgres",
		fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			username, password, addr, port, dbname))
	if err != nil {
		return
	}
	if len(addr_ro) > 0 {
		err = orm.RegisterDataBase(readOnlyDBName, "postgres",
			fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
				username, password, addr_ro, port, dbname))
		if err != nil {
			return
		}
		hasReadOnly = true
	}

	orm.RegisterModel(new(User))
	orm.RegisterModel(new(File))
	orm.RegisterModel(new(Document))
	orm.RegisterModel(new(Agent))
	orm.RegisterModel(new(Account))
	orm.RegisterModel(new(Bill))
	orm.RegisterModel(new(Order))
	orm.RegisterModel(new(Coupon))

	err = orm.RunSyncdb("default", false, true)

	if beego.BConfig.RunMode == "prod" {
		orm.Debug = false
	} else {
		orm.Debug = true
		//orm.DebugLog = orm.NewLog()
	}

	return
}
