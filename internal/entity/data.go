package entity

type Data struct {
	ID     uint64 `gorm:"primarykey"`
	UserID uint64
	Data   []byte
}
