package domain

type User struct {
	Login    string
	Password string
	Token    string
}

type AuthorizationData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
