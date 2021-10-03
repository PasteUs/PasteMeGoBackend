package user

import "github.com/PasteUs/PasteMeGoBackend/model/dao"

func init() {
	dao.CreateTable(&User{})
}

type User struct {
	Username string `json:"username" gorm:"type:varchar(32);primaryKey"`
	Password string `json:"password" gorm:"type:varchar(32)"`
	Email    string `json:"email" gorm:"type:varchar(128)"`
}
