package users

type UserDTO struct {
	Id       string
	Email    string
	Secret   string
	Verified bool
}

type CreateUserParams struct {
	Email    string
	Password string
}

type SessionDTO struct {
	UserId string
}

type LoginResult struct {
	User  UserDTO
	Token string
}
