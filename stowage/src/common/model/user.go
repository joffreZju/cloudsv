package model

import (
	"errors"
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type User struct {
	Id         int       `orm:"auto;pk;column(id)" json:"-"` // 用户ID，表内自增
	Uid        string    `orm:"unique"  json:"-"`
	Tel        string    `orm:"unique" json:",omitempty"`
	Password   string    `json:",omitempty"` // 密码
	Icon       string    `json:",omitempty"`
	Desc       string    `json:",omitempty"`
	Gender     int8      `json:",omitempty"`
	Address    string    `json:",omitempty"`
	LoginTime  time.Time `orm:"type(datetime)" json:",omitempty"` //登录时间
	CreateTime time.Time `orm:"type(datetime)" json:",omitempty"` //
	Mail       string    `json:",omitempty"`
	RegisterID string    `json:",omitempty"` // 用于给用户推送消息
	//Groups     []*Group  `orm:"-" json:",omitempty"` // 用户的所在组织
}

type UserAccount struct {
}

func CreateUser(u *User) (err error) {
	id, err := orm.NewOrm().Insert(u)
	if err != nil {
		return
	}
	u.Id = int(id)
	return
}

func UpdateUser(u *User) (err error) {
	_, err = orm.NewOrm().Update(u)
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
	err = NewOrm(ReadOnly).QueryTable("User").Filter("Tel", tel).One(u)
	return
}

type Group struct {
	Id         int    `orm:"auto;pk;"`
	Gid        string `orm:"unique" json:"-"`
	AdminId    int
	Name       string
	Desc       string
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

/*


func (ug *User) GetBindUids() (ids []int) {
	idsStr := strings.Split(ug.BindUids, "|")
	for _, idStr := range idsStr {
		id, _ := strconv.ParseInt(idStr, 10, 64)
		if id > 0 {
			ids = append(ids, int(id))
		}
	}
	return
}

func (ug *User) GetBindTelUid() (id int, err error) {
	for _, uid := range ug.GetBindUids() {
		if uid%10 == LoginTel {
			id = uid
			return
		}
	}
	err = fmt.Errorf("not found")
	return
}

func (ug *User) SetBindUids(ids []int) {
	idsStr := ""
	for _, id := range ids {
		if id > 0 {
			idsStr += strconv.Itoa(id) + "|"
		}
	}
	if len(idsStr) > 0 {
		ug.BindUids = idsStr[:len(idsStr)-1]
	}
	return
}

func GUID(ug *User) int {
	if ug == nil {
		return 0
	}
	no := int(crc32.ChecksumIEEE([]byte(ug.Account)) % 10)
	return ug.InnerId*100 + no*10 + ug.Type
}

func GetUser(id int) (ug *User, err error) {
	ug = &User{Id: id, InnerId: id / 100}
	sql := fmt.Sprintf("select * from %s where `id` = ?", ug.TableName())
	err = NewOrm(ReadOnly).Raw(sql, ug.InnerId).QueryRow(ug)
	return
}

func GetUsers(ids []int) (list []*User, err error) {
	mpList := make(map[int][]interface{})
	for _, id := range ids {
		mpList[id%100] = append(mpList[id%100], id/100)
	}
	o := NewOrm(ReadOnly)
	for tb, ids := range mpList {
		var tList []*User
		ug := &User{Id: tb}
		params := strings.Repeat("?,", len(ids))
		params = params[:len(params)-1]
		sql := fmt.Sprintf("select * from %s where id in (%s)", ug.TableName(), params)
		o.Raw(sql, ids...).QueryRows(&tList)
		list = append(list, tList...)
	}
	return
}

func GetUserByAccount(typ int, account string) (ug *User, err error) {
	ug = &User{Type: typ, Account: account}
	sql := fmt.Sprintf("select * from %s where `type` = ? and account = ?", ug.TableName())
	err = NewOrm(ReadOnly).Raw(sql, typ, account).QueryRow(ug)
	return
}

func CreateUser(ug *User) (err error) {
	values := []interface{}{ug.Type, ug.Account, ug.LoginTime}
	sql := fmt.Sprintf("insert into %s (`type`, account, login_time) values (?, ?, ?)", ug.TableName())
	if ug.InnerId > 0 {
		ug.Id = GUID(ug)
		sql = fmt.Sprintf("insert into %s (`type`, account, login_time, id, guid) values (?, ?, ?, ?, ?)", ug.TableName())
		values = append(values, ug.InnerId, ug.Id)
	}
	result, err := NewOrm().Raw(sql, values...).Exec()
	if err != nil {
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		return
	}
	ug.InnerId = int(id)
	ug.Id = ug.InnerId*100 + int(crc32.ChecksumIEEE([]byte(ug.Account))%10)*10 + ug.Type
	beego.Debug(ug)
	return UpdateUser(ug)
}

func UpdateUser(ug *User, fields ...string) (err error) {
	if ug.TableName() == emptyTable {
		return ErrEmptyTable
	}
	if ug.InnerId == 0 {
		ug.InnerId = ug.Id / 100
	}
	if len(fields) == 0 {
		fields = append(fields, "Id", "BindUids", "Name", "Nick", "Icon",
			"Sex", "Desc", "DefaultCarId", "Status", "LoginTime",
			"Tel", "UserName", "Password", "RYUserId", "RYToken", "Session", "DeviceId")
	}
	sql := fmt.Sprintf("update %s set PARAMS where id = ?", ug.TableName())

	params, values := "", []interface{}{}
	for _, f := range fields {
		switch f {
		case "Id":
			params += " `guid` = ? ,"
			values = append(values, ug.Id)
		case "BindUids":
			params += " `bind_uids`= ? ,"
			values = append(values, ug.BindUids)
		case "Name":
			params += " `name`= ? ,"
			values = append(values, ug.Name)
		case "Nick":
			params += " `nick`= ? ,"
			values = append(values, ug.Nick)
		case "Icon":
			params += " `icon`= ? ,"
			values = append(values, ug.Icon)
		case "Sex":
			params += " `sex`= ? ,"
			values = append(values, ug.Sex)
		case "Desc":
			params += "`desc`= ? ,"
			values = append(values, ug.Desc)
		case "DefaultCarId":
			params += " `default_car_id`= ? ,"
			values = append(values, ug.DefaultCarId)
		case "Status":
			params += " `status`= ? ,"
			values = append(values, ug.Status)
		case "LoginTime":
			params += " `login_time`= ? ,"
			values = append(values, ug.LoginTime.Format(TimeFormat))
		case "Tel":
			params += " `tel`= ? ,"
			values = append(values, ug.Tel)
		case "UserName":
			params += " `user_name`= ? ,"
			values = append(values, ug.UserName)
		case "Password":
			params += " `password`= ? ,"
			values = append(values, ug.Password)
		case "RYUserId":
			params += " `r_y_user_id`= ? ,"
			values = append(values, ug.RYUserId)
		case "RYToken":
			params += " `r_y_token`= ? ,"
			values = append(values, ug.RYToken)
		case "Session":
			params += " `session`= ? ,"
			values = append(values, ug.Session)
		case "DeviceId":
			params += " `device_id`= ? ,"
			values = append(values, ug.DeviceId)
		}
	}
	values = append(values, ug.InnerId)
	if len(params) > 1 {
		params = params[:len(params)-1]
	}
	sql = strings.Replace(sql, "PARAMS", params, 1)
	result, err := NewOrm().Raw(sql, values...).Exec()
	if err != nil {
		return
	}
	_, err = result.RowsAffected()
	return
}

func DeleteUser(id int) (err error) {
	ug := &User{Id: id, InnerId: id / 100, Status: 1}
	err = UpdateUser(ug, "Status")
	return
}
*/
