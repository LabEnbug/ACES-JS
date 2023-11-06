package algorithm

import (
	"backend/database/mysql"
	"backend/model"
	"encoding/json"
	"fmt"
	"sort"
	"math/rand"
)

// 函数将JSON字符串解析回JsonData结构体
func jsonStringToJSON(jsonString string) (model.RecommendMatrix, error) {
    var data model.RecommendMatrix
    err := json.Unmarshal([]byte(jsonString), &data)
    if err != nil {
        return model.RecommendMatrix{}, err
    }
    return data, nil
}

// 将 map[string]int 转换为按值从大到小排序的键值对切片
func sortMapByValue(myMap map[uint]int) []model.Pair {
	// 将 map 转换为切片
	pairs := make([]model.Pair, 0, len(myMap))
	for key, value := range myMap {
		pairs = append(pairs, model.Pair{key, value})
	}

	// 使用 sort.Slice 方法对切片进行排序
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})

	return pairs
}

func InitRecommendationModel() {
	userNum := 32
	for i := 2; i <= userNum; i++ {
		userId := uint(i)
		recommendMatrix := computeRecommendMatrixByUserId(userId)
		fmt.Println(string(recommendMatrix))

		ok := mysql.SetUserRecommendMatrix(userId, recommendMatrix)
		if !ok {
			errorMsg := "Unknown error."
			fmt.Println(errorMsg)
		} else {
			fmt.Println("Init user id =", userId, ", done")
		}
	}
	
}

func computeRecommendMatrixByUserId(userId uint) []byte {

	queryType := int(0)
	queryLimit := int(1000)  // get last 1000 videoes
	queryStart := int(0)

	// get userId's recently watched video list history, recently means limit 10 (for instance)
	watchedVideoList := mysql.GetVideoList(queryType, 0, "watched", queryLimit, queryStart, userId)
	// get userId's recently liked video list history, recently means limit 10 (for instance)
	likedVideoList := mysql.GetVideoList(queryType, userId, "liked", queryLimit, queryStart, userId)
	// get userId's recently favorit video list history, recently means limit 10 (for instance)
	favoriteVideoList := mysql.GetVideoList(queryType, userId, "favorite", queryLimit, queryStart, userId)

	if len(watchedVideoList) == 0 || len(likedVideoList) == 0 || len(favoriteVideoList) == 0 {
		fmt.Println("No video found.")
		
		nilLData := model.RecommendMatrix{
			TypeProbability: make([]float64, 1),
			UpCountMap: make(map[uint]int),
			PopularProbability: make([]float64, 1),
		}
		nilJsonData, _ := json.Marshal(nilLData)
		return nilJsonData
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

	// watched + liked + favorite + comment + forwarded == 100% . According to the degree of popularity, as part of the Recommendation Matrix
	popular_probability := []float64{0.1, 0.4, 0.15, 0.2, 0.15}
	// fmt.Println(popular_probability)

	matrixData := model.RecommendMatrix{
		TypeProbability: type_probability,
		UpCountMap: upCountMap,
		PopularProbability: popular_probability,
	}
	jsonData, _ := json.Marshal(matrixData)

	// fmt.Println(string(jsonData))

	return jsonData
}

func randIntByProb(p []float64) int {
    // Generate a random number in the range [0.0, 1.0)
    r := rand.Float64()
    
    // Accumulate the probabilities
    sum := 0.0
    for i, prob := range p {
        sum += prob
        // Check if the random number is less than the accumulated probability
        if r < sum {
            return i + 1 // Return the index + 1, because we want 1 to 5
        }
    }
    
    // In case of floating point arithmetic issues, return the last index
    return len(p)
}

func getRecommendMatrixByUserId(userId uint) model.RecommendMatrix {
	errorMsg := ""
	// get user info
	user, ok, errNo := mysql.GetRecommendMatrixByUserId(userId)
	if !ok {
		if errNo == 1 { // user not found
			errorMsg = "User not found."
		} else {
			errorMsg = "Unknown error."
		}
		fmt.Println(errorMsg)
	}
	matrixData := string(user.RecommendMatrix)

	recommendMatrix, err := jsonStringToJSON(matrixData)
    if err != nil {
        fmt.Println("Error parsing JSON: %s", err)
    }
    // fmt.Println(recommendMatrix.TypeProbability)
	// fmt.Println(recommendMatrix.UpCountMap)
	// fmt.Println(recommendMatrix.PopularProbability)
	return recommendMatrix
}

func GetRecommendVideoList(userId uint, requiredVideo int, queryStart int) []model.Video {
	var videoList []model.Video
	var oneVideo []model.Video
	queryLimit := 1

	recommendMatrix := getRecommendMatrixByUserId(userId)

	type_probability := recommendMatrix.TypeProbability
	upCountMap := recommendMatrix.UpCountMap
	popular_probability := recommendMatrix.TypeProbability
	
	// type of recommendation: base on video type, or up, or popularity, or traffic pool
	var recomType_probability []float64
	var upCountSortedPairsCut []model.Pair
	var up_probability []float64
	var recomTypeChoice int

	upCountSortedPairs := sortMapByValue(upCountMap)

	if len(upCountSortedPairs) > 5 {
		recomType_probability = []float64{0.4, 0.1, 0.2, 0.1, 0.2}

		upCountSortedPairsCut = upCountSortedPairs[1: 6]
		up_probability = make([]float64, len(upCountSortedPairsCut))
		countSum := 0
		for i, p := range upCountSortedPairsCut {
			up_probability[i] = float64(p.Value) 
			countSum += p.Value
		}
		for i := 0; i < len(up_probability); i++ {
			up_probability[i] = up_probability[i] / float64(countSum)
		}
	} else {
		// history is too short, so ignore video type and up
		recomType_probability = []float64{0, 0, 0.6, 0.2, 0.2}
	}
	
	for len(videoList) < requiredVideo {
		recomTypeChoice = randIntByProb(recomType_probability)
		switch recomTypeChoice {
		case 1:
			choice := randIntByProb(type_probability[1:10])
			oneVideo = mysql.GetOneRecommendVideoByProbabilityMatrix(1, choice, queryLimit, queryStart)
		case 2:
			choice := randIntByProb(up_probability)
			upId := int(upCountSortedPairsCut[choice-1].Key)
			oneVideo = mysql.GetOneRecommendVideoByProbabilityMatrix(2, upId, queryLimit, queryStart)
		case 3:
			choice := randIntByProb(popular_probability)
			oneVideo = mysql.GetOneRecommendVideoByProbabilityMatrix(3, choice, queryLimit, queryStart)
		case 4:
			// traffic pool
			oneVideo = mysql.GetOneRecommendVideoByProbabilityMatrix(4, 0, queryLimit, queryStart)
		case 5:
			oneVideo = mysql.GetOneRecommendVideoByProbabilityMatrix(5, 0, queryLimit, queryStart)
		}
		videoList = append(videoList, oneVideo...)
	}

	// fmt.Println(len(videoList))
	return videoList
}
