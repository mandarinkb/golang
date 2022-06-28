package repository

type User struct {
	ID       uint   `gorm:"column:USER_ID" json:"userId"`
	Username string `gorm:"column:USERNAME;type:varchar(50)" json:"username"`
	Password string `gorm:"column:PASSWORD;type:varchar(250)" json:"password"`
	UserRole string `gorm:"column:USER_ROLE;type:varchar(20)" json:"role"`
}

// กรณีต้องการเปลี่ยนชื่อ table ใหม่
func (User) TableName() string {
	return "USERS"
}

type UserRepository interface {
	Authenticate(username string) (*User, error)
	Read() ([]User, error)
	ReadById(id int) (*User, error)
	Create(user User) error
	Update(user User) error
	Delete(id int) error
}
