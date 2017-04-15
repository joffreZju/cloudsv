package model

import (
	//"errors"
	//"fmt"

	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
	DateFormat = "2006-01-02"
)

type User struct {
	Id         int       `orm:"auto;pk;column(id)" ` // 用户ID，表内自增
	Unum       string    `orm:"null"  json:"-"`
	Tel        string    `orm:"unique;size(15)" json:",omitempty"`
	Password   string    `json:"-"`                              // 密码
	UserName   string    `orm:"null;size(64)" json:",omitempty"` // 用户名
	Icon       string    `orm:"null;size(64)" json:",omitempty"`
	Descp      string    `orm:"null" json:",omitempty"`
	Gender     int8      `orm:"null" json:",omitempty"`
	Address    string    `orm:"null;size(64)" json:",omitempty"`
	LoginTime  time.Time `orm:"type(datetime);null" json:",omitempty"`         //登录时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime)" json:",omitempty"` //
	Mail       string    `orm:"null;size(64)" json:",omitempty"`
	UserType   int       `orm:"default(1)" json:",omitempty"` //1 普通用户,2 代理商
	Referer    string    `orm:"null;size(16)" json:",omitempty"`
	RegisterID string    `orm:"null;size(32)" json:",omitempty"` // 用于给用户推送消息
	//Groups     []*Group  `orm:"-" json:",omitempty"` // 用户的所在组织
}

func (u *User) TableName() string {
	return "allsum_user"
}

func CreateUser(u *User) (err error) {
	id, err := orm.NewOrm().Insert(u)
	if err != nil {
		return
	}
	u.Id = int(id)
	return
}
func CreateUserIfNotExist(u *User) (err error) {
	u.CreateTime = time.Now()
	_, id, err := orm.NewOrm().ReadOrCreate(u, "tel")
	if err != nil {
		return
	}
	u.Id = int(id)
	return
}

func GetUsersByReferer(tel string) (list []*User, err error) {
	_, err = NewOrm(ReadOnly).QueryTable("allsum_user").Filter("Referer", tel).All(&list)
	return
}

func GetUserCountsByReferer(tel string) (c int64, err error) {
	c, err = NewOrm(ReadOnly).QueryTable("allsum_user").Filter("Referer", tel).Count()
	return
}

func UpdateUser(u *User, fields ...string) (err error) {
	if len(fields) == 0 {
		fields = append(fields, "Id", "Icon",
			"Gender", "Descp", "Address", "LoginTime",
			"Tel", "UserName", "Password", "Mail", "Referer")
	}
	sql := fmt.Sprintf("update allsum_user set PARAMS where id = ?")

	params, values := "", []interface{}{}
	for _, f := range fields {
		switch f {
		case "UserName":
			params += " user_name= ? ,"
			values = append(values, u.UserName)
		case "Descp":
			params += " descp= ? ,"
			values = append(values, u.Descp)
		case "Gender":
			params += " gender= ? ,"
			values = append(values, u.Gender)
		case "Address":
			params += " address= ? ,"
			values = append(values, u.Address)
		case "LoginTime":
			params += " login_time= ? ,"
			values = append(values, time.Now().Format(TimeFormat))
		case "Mail":
			params += " mail= ? ,"
			values = append(values, u.Mail)
		case "Password":
			params += " password= ? ,"
			values = append(values, u.Password)

		}
	}
	values = append(values, u.Id)
	if len(params) > 1 {
		params = params[:len(params)-1]
	}
	sql = strings.Replace(sql, "PARAMS", params, 1)
	fmt.Printf("---sql:%s\n--%s\n", sql, params)
	result, err := NewOrm().Raw(sql, values...).Exec()
	if err != nil {
		return
	}
	_, err = result.RowsAffected()
	return
}

func GetUser(id int) (u *User, err error) {
	u = &User{Id: id}
	err = NewOrm(ReadOnly).Read(u)
	return
}

func DeleteUser(u *User) (err error) {
	_, err = orm.NewOrm().Delete(u)
	return
}

func GetUserByTel(tel string) (u *User, err error) {
	u = new(User)
	err = NewOrm(ReadOnly).QueryTable("allsum_user").Filter("Tel", tel).One(u)
	return
}

/*
func GetUsers(ids []int) (list []*User, err error) {

}*/

/*
type Group struct {
	Id         int    `orm:"auto;pk;"`
	Gid        string `orm:"unique" json:"-"`
	AdminId    int
	Name       string
	Descp      string
	Users      []*GroupUser `orm:"rel(m2m)" json:",omitempty"`
	CreateTime time.Time    `orm:"type(datetime)" json:",omitempty"` //
}

func GetGroup(gid int) (g *Group, err error) {
	g = new(Group)
	g.Id = gid
	o := NewOrm(ReadOnly)
	if err = o.Read(g); err != nil {
		return nil, err
	}
	return
}

func DeleteGroup(gid int) (err error) {
	g := Group{Id: gid}
	_, err = orm.NewOrm().Delete(&g)
	return
}

func InsertGroup(g *Group) (err error) {
	id, err := orm.NewOrm().Insert(g)
	if err != nil {
		return
	}
	g.Id = int(id)
	return
}

func UpdateGroup(g *Group, fileds ...string) (err error) {
	_, err = orm.NewOrm().Update(g, fileds...)
	return
}

func AddGroupUser(gid int, uid int) (err error) {
	g, err := GetGroup(gid)
	if err != nil {
		return
	}
	u, err := GetUser(uid)
	if err != nil {
		return
	}
	_, err = orm.NewOrm().QueryM2M(g, "Users").Add(u)
	return
}

func DeleteGroupUser(gid int, uid int) (err error) {
	g, err := GetGroup(gid)
	if err != nil {
		return
	}
	u, err := GetUser(uid)
	if err != nil {
		return
	}
	_, err = orm.NewOrm().QueryM2M(g, "Users").Remove(u)
	return
}

type GroupUser struct {
	Id           int          `orm:"auto;pk;" json:",omitempty"`
	Group        *Group       `orm:"rel(fk)" json:",omitempty"`
	User         *User        `orm:"rel(fk)" json:",omitempty"`
	DepartmentId int          `json:",omitempty"` // 部门ID
	Role         int          // 用户的角色
	Temp         string       `orm:"column(permissions)" json:"Permission"` // 1|2|3 特殊权限的表
	Permissions  map[int]bool `orm:"-" json:"-"`                            // 对应的权限列表
	Status       int          `json:"-"`                                    // 用户状态 0 正常 1 删除
}
*/
