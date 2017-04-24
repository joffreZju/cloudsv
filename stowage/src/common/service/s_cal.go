package service

import (
	"common/model"

	"github.com/astaxie/beego"
)

func InsertTpl(t *model.CalTemplate) (err error) {
	err = model.InsertTemplate(t)
	if err != nil {
		beego.Error(err)
	}
	return
}

func GetTpl(uid int) (t *model.CalTemplate, err error) {
	t, err = model.GetTemplate(uid)
	if err != nil {
		beego.Error(err)
		err = errcode.ErrNoRecord
	}
	return
}
