package domain

import (
	"time"
)

// 领域包是抽象，面对服务
type User struct {
	Id       int64
	Email    string
	Password string
	Phone    string
	Ctime    time.Time

	WechatInfo WechatInfo
	//Addr Address

}

type UserProfile struct {
	Id           int64
	Gender       string
	NickName     string
	Introduction string
	BirthDate    time.Time
}

//type Address struct {
//	Province string
//	Region   string
//}
