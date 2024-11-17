package entity

type User struct {
	ID       uint `gorm:"primarykey"`
	Login    string
	Password string
}
