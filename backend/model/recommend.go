package model



// 定义一个键值对的结构体
type Pair struct {
	Key   uint
	Value int
}

// 用户的推荐概率矩阵，包含2个float64的概率切片和一个观看视频上传者(up主)的视频数量的字典
type RecommendMatrix struct {
	TypeProbability     []float64         `json:"typeProbability"`
	UpCountMap           map[uint]int      `json:"upCountMap"`
	PopularProbability  []float64         `json:"popularProbability"`
}