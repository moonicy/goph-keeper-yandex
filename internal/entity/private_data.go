package entity

type PrivateData struct {
	ID     uint `gorm:"primarykey"`
	UserID uint
	Data   []byte
}
