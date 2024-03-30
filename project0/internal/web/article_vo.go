package web

type ArticleEdit struct {
	// 没加json完蛋?
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ArticlePublishReq struct {
	// 没加json完蛋?
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ArticleVo struct {
	// 这里在ming version是驼峰json，我的是下划线
	Id         int64  `json:"id,omitempty"`
	Title      string `json:"title,omitempty"`
	Content    string `json:"content,omitempty"`
	Abstract   string `json:"abstract,omitempty"`
	Status     uint8  `json:"status,omitempty"`
	AuthorId   int64  `json:"author_id,omitempty"`
	AuthorName string `json:"author_name,omitempty"`
	Ctime      string `json:"ctime,omitempty"`
	Utime      string `json:"utime,omitempty"`

	ReadCnt    int64 `json:"read_cnt,omitempty"`
	LikeCnt    int64 `json:"like_cnt,omitempty"`
	CollectCnt int64 `json:"collect_cnt,omitempty"`
	Liked      bool  `json:"liked,omitempty"`
	Collected  bool  `json:"collected,omitempty"`
}

type LikeReq struct {
	Id   int64 `json:"id"`
	Like bool  `json:"like"`
}

type CollectReq struct {
	Id  int64 `json:"id"`
	Cid int64 `json:"uid"`
}
