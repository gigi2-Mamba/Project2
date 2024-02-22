package domain

type Interactive struct {
	//榜单模型后加入了bizId
	BizId      int64
	ReadCnt    int64 `json:"read_cnt,omitempty"`
	LikeCnt    int64 `json:"like_cnt,omitempty"`
	CollectCnt int64 `json:"collect_cnt,omitempty"`
	Liked      bool  `json:"liked,omitempty"`
	Collected  bool  `json:"collected,omitempty"`
}
