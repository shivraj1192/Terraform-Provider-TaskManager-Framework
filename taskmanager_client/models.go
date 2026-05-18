package taskmanager_client

// HTTPErrorResponse - godoc
type HTTPErrorResponse struct {
	Status    int64       `json:"status"`
	Message   string      `json:"message"`
	ErrorCode string      `json:"errorCode"`
	Details   interface{} `json:"details"`
}

type UserModel struct {
	User User `json:"user"`
}

type User struct {
	ID       int    `json:"ID"`
	Uname    string `json:"uname"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
