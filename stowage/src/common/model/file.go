package model

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type File struct {
	Id         string    `orm:"pk;column(id)" json:"id"`
	Userid     int       `json:"userid"`
	Name       string    `json:"name"`
	Size       int       `json:"size"`
	Md5        string    `json:"md5,omitempty"`
	CreateTime time.Time `orm:"type(datetime)" json:"create_time,omitempty"`
	Data       string    `orm:"type(text)" json:"-"`
}

type Document struct {
	Id       int64  `orm:"auto;pk;column(id)" json:"id"`
	DocType  int    `json:"docTypei,omitempty"`
	Uploader int64  `json:"uploader,omitempty"`
	Fileid   string `json:"fileid,omitempty"`
	Desc     string `orm:type(text) json:"desc,omitempty"`
	Status   string `json:status,omitempty`
}

//插入文件
func CreateFile(f *File) (err error) {
	_, err = orm.NewOrm().Insert(f)
	return
}

//删除文件
func DeleteFile(f *File) (err error) {
	_, err = orm.NewOrm().Delete(f)
	return
}

//获取文件按文件id
func GetFile(id string) (f *File, err error) {
	f = &File{Id: id}
	err = orm.NewOrm().Read(f)
	return
}

//获取文件按用户/// 此处以后可能有分页
func GetFilesByUser(userid int) (fs []*File, err error) {
	num, err := orm.NewOrm().QueryTable("File").Filter("userid", userid).All(&fs)
	fmt.Println("----------hebl filenum::", num)
	return
}

//-------------------------------Document

//创建文档
func CreateDocument(doc *Document) (err error) {
	id, err := orm.NewOrm().Insert(doc)
	if err != nil {
		return
	}
	doc.Id = id
	return
}

//更新文档
func UpdateDocument(doc *Document) (err error) {
	_, err = orm.NewOrm().Update(doc)
	return
}

//获取文档
func GetDocument(id int64) (doc *Document, err error) {
	doc = &Document{Id: id}
	err = orm.NewOrm().Read(doc)
	return
}

//根据userid获取文档
func GetDocumentByUserId(userid int64) (docs []*Document, err error) {
	num, err := orm.NewOrm().QueryTable("Document").Filter("Uploader", userid).All(&docs)
	fmt.Println("---------------hebl get document::", num)
	return
}

//删除文档
func DeleteDocument(doc *Document) (err error) {
	_, err = orm.NewOrm().Delete(doc)
	return
}
