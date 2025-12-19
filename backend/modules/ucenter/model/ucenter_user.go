package model

import (
	"time"

	"github.com/cool-team-official/cool-admin-go/cool"
)

const TableNameUcenterUser = "ucenter_user"

// UcenterUser mapped from table <ucenter_user>
type UcenterUser struct {
	*cool.Model
	Name       string    `gorm:"column:name;not null;comment:名称" json:"name"`
	Token      string    `gorm:"column:token;not null;comment:Token" json:"token"`
	Permis     string    `gorm:"column:permis;not null;comment:权限，可以访问的服务，逗号分隔" json:"permis"`
	Email      string    `gorm:"column:email;not null;comment:用户邮箱" json:"email"`
	Status     bool      `gorm:"column:status;not null;comment:状态;default:0" json:"status"`
	ExpireTime time.Time `gorm:"column:expire_time;not null;comment:过期时间" json:"expire_time"`
}

// TableName UcenterUser's table name
func (*UcenterUser) TableName() string {
	return TableNameUcenterUser
}

// GroupName UcenterUser's table group
func (*UcenterUser) GroupName() string {
	return "default"
}

// NewUcenterUser create a new UcenterUser
func NewUcenterUser() *UcenterUser {
	return &UcenterUser{
		Model: cool.NewModel(),
	}
}

// init 创建表
func init() {
	cool.CreateTable(&UcenterUser{})
}
