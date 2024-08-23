package moex_api_client

import (
	"time"
)

type MoexApiClient struct {
	BaseUrl string
}

func NewMoexApiClient() *MoexApiClient {
	return &MoexApiClient{BaseUrl: getBaseUrl()}
}

func isInside[T comparable](value T, array []T) bool {
	for _, item := range array {
		if item == value {
			return true
		}
	}

	return false
}

func toString(value interface{}) string {
	str, isString := value.(string)
	if isString {
		return str
	}
	return ""
}

func toFloat64(value interface{}) float64 {
	num, isFloat := value.(float64)
	if isFloat {
		return num
	}
	return 0
}

func getBaseUrl() string {
	return "https://iss.moex.com/iss/"
}

func getCurrentTime() string {
	currTime := time.Now()
	return currTime.Format("2006-01-02")
}
