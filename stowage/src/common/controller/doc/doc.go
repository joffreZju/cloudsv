package doc

import (
	"common/controller/base"
	"common/service"
	"io/ioutil"
	"s4s/common/lib/errcode"
	"stowage/src/common/lib/util"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
}

func (c *Controller) AddDocument() {
	fno := c.GetString("fno")
	desc := c.GetString("desc")
	dtp, _ := c.GetInt("dtp")
	uid := c.UserId
	d := model.Documnet{
		DocType:  dtp,
		Uploader: uid,
		FileNo:   fno,
		Desc:     desc,
		Status:   model.DocUsing,
	}
	err = service.NewDocument(&d)
	if err != nil {
		beego.Error("create documnet failed:", err)
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(d)
}

func (c *Controller) SetDocStatus() {
	dtp, _ := c.GetInt("dtp")
	id, _ := c.GetInt("id")
	status, _ := c.GetInt("status")
	err := service.SetDocStatus(status, dtp, id)
	if err != nil {
		beego.Error("set documnet status failed:", err)
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(d)
}

func (c *Controller) GetDocList() {
	dtp, _ := c.GetInt("dtp")
	list, err := model.GetDocListByType(dtp)
	if err != nil {
		beego.Error("get documnet list failed:", err)
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(list)
}

func (c *Controller) GetDocUsing() {
	dtp, _ := c.GetInt("dtp")
	doc, err := model.GetDocByType(dtp)
	if err != nil {
		beego.Error("get documnet list failed:", err)
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(*doc)
}

func (c *Controller) AddFile() {
	uid, _ := c.GetInt("id")
	f, h, err := c.GetFile("doc")
	if err != nil {
		beego.Error("User.UploadDoc error: ", err)
		c.ReplyErr(errcode.ErrUserUploadDocFailed)
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	filename := strings.Replace(time.Now().Format("060102150405.999"), ".", "", 1)
	if pos := strings.Index(h.Filename, "."); pos >= 0 {
		filename += h.Filename[pos:]
	}
	mime := h.Header.Get("Content-Type")
	beego.Debug("file content-type:", mime)
	beego.Debug("data info :", filename)
	md5 := util.Md5Cal2String(data)
	sz := len(data) / 1024
	no := md5 + util.UniqueRandom()
	file := File{
		Uid:    uid,
		Name:   filename,
		Size:   sz,
		Md5:    md5,
		Data:   data,
		FileNo: no,
		Mime:   mime,
	}
	err = service.NewFile(&file)
	if err != nil {
		beego.Error("create file failed:", err)
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(file)
}

func (c *Controller) FileDownload() {
	no := c.GetString("fno")
	f, err := model.GetFile(no)
	if err != nil {
		beego.Error("get file failed:", err)
		c.ReplyErr(err)
		return
	}
	c.ReplyFile(f.Mime, f.Name, []byte(f.Data))

}