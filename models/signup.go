package models

type SignupData struct {
	Username string `json:"username"  validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

type LoggedInUserdata struct {
	Username     string
	Email	     string
	Hashpassword string
}
