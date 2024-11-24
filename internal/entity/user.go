package entity

type User struct {
	ID       uint64 `gorm:"primarykey"`
	Login    string
	Password string
}
