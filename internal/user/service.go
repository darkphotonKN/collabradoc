package user

func FindAllUsers() ([]User, error) {

	return QueryAllUsers()
}

func SignUp(name string, email string, password string) (User, error) {

	// TODO: HASH PASSWORD
	return CreateUser(name, email, password)
}
