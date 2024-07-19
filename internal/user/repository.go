package user

import "github.com/darkphotonKN/collabradoc/internal/db"

func FindAllUsers() ([]User, error) {
	// dummy data
	usersData := []User{
		User{
			ID:       "1",
			Name:     "Bob",
			Email:    "bob@test.com",
			Password: "123456",
		},
		User{
			ID:       "2",
			Name:     "Nick",
			Email:    "nick@test.com",
			Password: "123456",
		},
	}

	return usersData, nil
}

func CreateUser(name string, email string, password string) (User, error) {

	db := db.DBCon

	newUser := User{
		Name:  name,
		Email: email,
		// TODO: HASH PASSWORD
		Password: password,
	}
	result := db.Create(&newUser)

	if result.Error != nil {
		return newUser, result.Error
	}

	return newUser, nil
}
