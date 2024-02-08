package domain

import "time"

type Article struct {
	Id int64
	Title string
	Content string
    Status ArticleStatus
	//Fucker int
	Author
	Ctime  time.Time
	Utime time.Time
}


//截取一定长度做摘要
func (a Article) Abstract()  string{
	// 这样命名直指本质
     str := []rune(a.Content)
	 if len(str) > 128 {
		 str = str[:128]
	 }
	 return string(str)
}

// 面向领域来说
type  ArticleStatus uint8

const (
	// 考虑到序列化问题 ，未知状态
	ArticleStatusUnknown ArticleStatus  = iota
	ArticleStatusUnPublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

type Author struct {
	Id int64
	Name string
}

func (a ArticleStatus) ToUint8() uint8{
	return uint8(a)
}
