package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	DocUsing = 1
	DocHide  = 2
)

type File struct {
	Id         int       `orm:"pk;auto;column(id)"`
	FileNo     string    `orm:"size(50)"`
	Uid        int       `json:",omitempty"`
	Name       string    `orm:"size(50)" json:"name,omitempty"`
	Mime       string    `orm:"size(250)"`
	Size       int       `json:"size,omitempty"`
	Md5        string    `orm:"size(50)" json:"md5,omitempty"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime)" json:",omitempty"`
	Data       string    `orm:"type(text),null" json:"-"`
}

type Document struct {
	Id       int `orm:"auto;pk;column(id)"`
	DocType  int
	Uploader int    `json:",omitempty"`
	FileNo   string `json:",omitempty"`
	Desc     string `orm:"null" json:",omitempty"`
	Status   int    `json:",omitempty"`
}

//插入文件
func CreateFile(f *File) (err error) {
	_, err = orm.NewOrm().Insert(f)
	return
}

//删除文件
func DeleteFile(no string) (err error) {
	f := &File{
		FileNo: no,
	}
	_, err = orm.NewOrm().Delete(f, "FileNo")
	return
}

//获取文件按文件编号
func GetFile(id string) (f *File, err error) {
	f = &File{FileNo: id}
	err = orm.NewOrm().QueryTable("File").Filter("FileNo", id).One(f)
	return
}

//检查文件是否存在
func CheckFileExist(id string) (exist bool) {
	exist = NewOrm(ReadOnly).QueryTable("File").Filter("FileNo", id).Exist()
	return
}

//获取文件按用户/// 此处以后可能有分页
func GetFilesByUser(userid int) (fs []*File, err error) {
	_, err = NewOrm(ReadOnly).QueryTable("File").Filter("Uid", userid).All(&fs)
	return
}

//-------------------------------Document

//创建文档
func CreateDocument(doc *Document) (err error) {
	id, err := orm.NewOrm().Insert(doc)
	if err != nil {
		return
	}
	doc.Id = int(id)
	return
}

//更新文档
func UpdateDocument(doc *Document, fields ...string) (err error) {
	_, err = orm.NewOrm().Update(doc, fields...)
	return
}

//获取文档
func GetDocument(id int) (doc *Document, err error) {
	doc = &Document{Id: id}
	err = NewOrm(ReadOnly).Read(doc)
	return
}

func SetDocHide(tp int) (err error) {
	_, err = orm.NewOrm().QueryTable("Documnet").Filter("DocType", tp).Update(orm.Params{
		"Status": DocHide,
	})

	return
}

//根据文档类型获取当前有效文档
func GetDocByType(tp int) (doc *Document, err error) {
	err = orm.NewOrm().QueryTable("Document").Filter("DocType", tp).Filter("Status", DocUsing).One(doc)
	return
}

func GetDocListByType(tp int) (docs *[]Document, err error) {
	_, err = orm.NewOrm().QueryTable("Documnet").Filter("DocType", tp).All(&docs)
	return
}

/*
//获取文档列表
func ListDocument(limit int, mark int) (docs []*Document, err error) {
	_, err = NewOrm(ReadOnly).QueryTable("Document").Limit(limit, mark).All(&docs)
	return
}
//根据userid获取文档
func GetDocumentByUserId(userid int64) (docs []*Document, err error) {
	_, err = NewOrm(ReadOnly).QueryTable("Document").Filter("Uploader", userid).All(&docs)
	return
}

//删除文档
func DeleteDocument(doc *Document) (err error) {
	_, err = orm.NewOrm().Delete(doc)
	return
}
*/
