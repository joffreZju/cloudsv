package model

import (
	"runtime/debug"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func LoadRelations(m interface{}, fileds ...string) (err error) {
	o := NewOrm(ReadOnly)
	for _, f := range fileds {
		func() {
			defer func() {
				e := recover()
				if e != nil {
					beego.Critical("panic:", e, string(debug.Stack()))
				}
			}()
			_, err = o.LoadRelated(m, f)
			if err != nil && err != orm.ErrNoRows {
				beego.Error("load related error:", err)
			} else {
				err = nil
			}
		}()
	}
	return
}

func LoadRelationsFromMainDB(m interface{}, fileds ...string) (err error) {
	o := NewOrm()
	for _, f := range fileds {
		func() {
			defer func() {
				e := recover()
				if e != nil {
					beego.Critical("panic:", e, string(debug.Stack()))
				}
			}()
			_, err = o.LoadRelated(m, f)
			if err != nil && err != orm.ErrNoRows {
				beego.Error("load related error:", err)
			} else {
				err = nil
			}
		}()
	}
	return
}
