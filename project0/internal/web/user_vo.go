package web



type SignUpReq struct { // 怎么把证
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type UserProfileReq struct {
	Nickname     string `json:"nickname"`
	Gender       string `json:"gender"`
	Introduction string `json:"introduction"`
	Birthday     string `json:"birthday"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginSmsReq struct {
	Phone string `json:"phone"`
	// 因为大写有问题
	Code string `json:"code"`
}

type SendLoginCodeReq struct {
	Phone string `json:"phone"`
}