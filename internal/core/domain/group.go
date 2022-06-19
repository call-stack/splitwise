package domain

type Group struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserID []int  `json:"userid"`
}

type GroupEntity struct {
	ID   int `gorm:"column:gid"`
	Name string
}

type UserGroup struct {
	GID int `gorm:"column:gid"`
	UID int `gorm:"column:uid"`
}
