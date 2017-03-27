package main

import (
	//"common/filter"
	"common/model"
	"common/service"
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	//"common/lib/cache"
	"time"

	"github.com/astaxie/beego"
	"github.com/ysqi/tokenauth2beego"
)

var (
	key = "54ef2c87604c32c482dd81f1f63f518e"
)

func Init() (err error) {
	rand.Seed(time.Now().UnixNano())

	// init log
	err = InitLog()
	if err != nil {
		beego.Error("init log failed : ", err)
		return
	}

	// init tokenauth
	tokenauth2beego.Init(key)

	// init pgsql
	err = model.InitPgSQL(key)
	if err != nil {
		beego.Error("init database failed : ", err)
		return
	}

	// init service
	err = service.Init()
	if err != nil {
		beego.Error("init service failed : ", err)
		return
	}
	beego.Info("init success, start server ...")

	// 生成pprof文件
	go Printpprof()

	return
}

func InitLog() (err error) {
	//filter.LoadLogFilter()
	typ := beego.AppConfig.String("log::type")
	cons := beego.AppConfig.String("log::params")
	return beego.SetLogger(typ, cons)
}

func Printpprof() {
	if beego.BConfig.RunMode == "prod" {
		return
	}
	cpuprofile := "/tmp/cpuprofile"
	for i := 0; i < 60; i++ {
		func() {
			defer recover()
			filename := fmt.Sprintf("%s_%d.pprof", cpuprofile, i)
			f, err := os.Create(filename)
			if err != nil {
				beego.Error("create cpu profile error:", err)
				return
			}
			defer f.Close()
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
			time.Sleep(time.Minute)
		}()
	}
}