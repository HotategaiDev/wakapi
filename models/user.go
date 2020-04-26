package models

type User struct {
	ID       string `json:"id" gorm:"primary_key"`
	ApiKey   string `json:"api_key" gorm:"unique"`
	Password string `json:"-"`
}
