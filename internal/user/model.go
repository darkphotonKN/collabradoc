package user

type User struct {
	ID       string `gorm:"primarykey" json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}
