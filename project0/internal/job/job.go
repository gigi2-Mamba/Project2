package job

/*
User: society-programmer on
Date: 2024/2/21  周三
Time: 16:18
*/

type Job interface {
	Name() string
	Run() error
}
