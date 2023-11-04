package algorithm

import (
	"backend/database/mysql"
	"encoding/json"
	"fmt"
	"sort"
)

// 定义一个键值对的结构体
type pair struct {
	key   uint
	value int
}

type recommendMatrix struct {
	TypeProbability     []float64         `json:"typeProbability"`
	UpCountMap           map[uint]int      `json:"upCountMap"`
	PopularProbability  []float64         `json:"popularProbability"`
}

// 将 map[string]int 转换为按值从大到小排序的键值对切片
func sortMapByValue(myMap map[uint]int) []pair {
	// 将 map 转换为切片
	pairs := make([]pair, 0, len(myMap))
	for key, value := range myMap {
		pairs = append(pairs, pair{key, value})
	}

	// 使用 sort.Slice 方法对切片进行排序
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].value > pairs[j].value
	})

	return pairs
}

func InitRecommendationModel() {
	userId := uint(1)
	recommendMatrix := getRecommendMatrixByUserId(userId)
	fmt.Println(string(recommendMatrix))

	ok := mysql.SetUserRecommendMatrix(userId, recommendMatrix)
	if !ok {
		errorMsg := "Unknown error."
		fmt.Println(errorMsg)
	} else {
		fmt.Println("Init user id =", userId, ", done")
	}
}

func getRecommendMatrixByUserId(userId uint) []byte {

	queryType := int(0)
	queryLimit := int(2)  // get all, no limit
	queryStart := int(0)

	// get userId's recently watched video list history, recently means limit 10 (for instance)
	watchedVideoList := mysql.GetVideoList(queryType, 0, "watched", queryLimit, queryStart, userId)
	// get userId's recently liked video list history, recently means limit 10 (for instance)
	likedVideoList := mysql.GetVideoList(queryType, userId, "liked", queryLimit, queryStart, userId)
	// get userId's recently favorit video list history, recently means limit 10 (for instance)
	favoriteVideoList := mysql.GetVideoList(queryType, userId, "favorite", queryLimit, queryStart, userId)

	if len(watchedVideoList) == 0 && len(likedVideoList) == 0 && len(favoriteVideoList) == 0 {
		fmt.Println("No video found.")
		
		errorData := recommendMatrix{
			TypeProbability: make([]float64, 1),
			UpCountMap: make(map[uint]int),
			PopularProbability: make([]float64, 1),
		}
		errorJsonData, _ := json.Marshal(errorData)
		return errorJsonData
	}

	// get standard data which is required by recommendation matrix
	// type_weight := 0.5
	// up_weight := 0.1
	// popular_weight := 0.3
	// traffic_weight := 0.1

	// According to the Type of video users watched
	type_probability := make([]float64, 40)
	type_probability[10] = float64(len(watchedVideoList))
	type_probability[20] = float64(len(likedVideoList))
	type_probability[30] = float64(len(favoriteVideoList))

	watched_weight := 0.1
	liked_weight := 0.4
	favorit_weight := 0.5

	watched_type_nums := make([]int, 10)
	for i := 0; i < len(watchedVideoList); i++ {
		watched_type_nums[watchedVideoList[i].Type] += 1
	}
	for i := 1; i < 10; i++ {
		type_probability[i] += watched_weight * float64(watched_type_nums[i]) / float64(len(watchedVideoList))
		type_probability[10 + i] = float64(watched_type_nums[i])
	}

	liked_type_nums := make([]int, 10)
	for i := 0; i < len(likedVideoList); i++ {
		liked_type_nums[likedVideoList[i].Type] += 1
	}
	for i := 1; i < 10; i++ {
		type_probability[i] += liked_weight * float64(liked_type_nums[i]) / float64(len(likedVideoList))
		type_probability[20 + i] = float64(liked_type_nums[i])
	}

	favorite_type_nums := make([]int, 10)
	for i := 0; i < len(favoriteVideoList); i++ {
		favorite_type_nums[favoriteVideoList[i].Type] += 1
	}
	for i := 1; i < 10; i++ {
		type_probability[i] += favorit_weight * float64(favorite_type_nums[i]) / float64(len(favoriteVideoList))
		type_probability[30 + i] = float64(favorite_type_nums[i])
	}

	// fmt.Println(len(type_probability), type_probability)

	upCountMap := make(map[uint]int)
	for i := 0; i < len(watchedVideoList); i++ {
		upid := watchedVideoList[i].UserId
		upCountMap[upid] += 1
	}
	for i := 0; i < len(likedVideoList); i++ {
		upid := likedVideoList[i].UserId
		upCountMap[upid] += 1
	}
	for i := 0; i < len(favoriteVideoList); i++ {
		upid := favoriteVideoList[i].UserId
		upCountMap[upid] += 1
	}
	allUpCount := len(watchedVideoList) + len(likedVideoList) + len(favoriteVideoList)
	upCountMap[0] = allUpCount

	// fmt.Println(upCountMap)

	// upCountSortedPairs := sortMapByValue(upCountMap)
	// for _, p := range upCountSortedPairs {
	// 	fmt.Printf("%d:%d ", p.key, p.value)
	// }

	// watched + liked + favorite + comment + forwarded == 100% . According to the degree of popularity, as part of the Recommendation Matrix
	popular_probability := []float64{0.1, 0.4, 0.15, 0.2, 0.15}
	// fmt.Println(popular_probability)

	matrixData := recommendMatrix{
		TypeProbability: type_probability,
		UpCountMap: upCountMap,
		PopularProbability: popular_probability,
	}
	jsonData, _ := json.Marshal(matrixData)

	// fmt.Println(string(jsonData))

	return jsonData
}
