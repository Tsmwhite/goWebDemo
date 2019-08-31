package model

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const (
	DB_DRIVER = "root:196914@tcp(127.0.0.1:3306)/spider?charset=utf8"
)

var DB *gorm.DB

func init(){
	var err error
	DB , err = gorm.Open("mysql",DB_DRIVER)
	if(err != nil){
		panic(err)
	}
	DB.SingularTable(true)
}

type Member struct {
	Id 				int 	`json:"id"`
	Mobile 			string  `json:"mobile"`
	Password 		string  `json:"password"`
	Name 			string  `json:"name"`
	Salt			string  `json:"salt"`
	Headimg			string  `json:"headimg"`
	Balance			float64 `json:"balance"`
	Login_time		int64  	`json:"login_time"`
	Register_time	int64	`json:"register_time"`
	Login_ip		string	`json:"login_ip"`
	Register_ip		string	`json:"register_ip"`
	Is_hid			int 	`json:"is_hid"`
	Is_del			int     `json:"is_del"`
	Role_id			int		`json:"role_id"`
}
func (member Member) TableName() string {
	return "z3_member"
}

func (mem *Member) Insert() error{
	return DB.Create(mem).Error
}

func GetMemberInfo (field string,value interface{}) (*Member, error) {
	var member Member
	err := DB.First(&member, field+" = ?", value).Error
	return &member, err
}

//更新数据
func Update(table string,UpdateData map[string]interface{}) error{
	return DB.Table(table).Where("id = ?",UpdateData["id"]).Updates(UpdateData).Error
}