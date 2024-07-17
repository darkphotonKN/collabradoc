package user

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
