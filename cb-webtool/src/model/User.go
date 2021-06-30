package model

// 사용하지 않는것으로 보임
type User struct {
	//user(username, password, email)
	UserID   string `json:"userid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Users []User
