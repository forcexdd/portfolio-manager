package client

import (
	"errors"
	"net/http"
	"time"
)

type MoexApiClient struct {
	BaseUrl string
}

func NewMoexApiClient() *MoexApiClient {
	return &MoexApiClient{BaseUrl: baseUrl}
}

func sendGETRequest(url string) (*http.Response, error) {
	wasResponse := false
	var err error
	var response *http.Response
	var i int
	for i = 0; wasResponse == false && i < maxRetries; i++ {
		response, err = http.Get(url)
		if err != nil {
			return nil, err
		}

		wasResponse = true
		err = checkHTTPResponse(response)
		if err != nil {
			if errors.Is(err, errBadHTTPResponse) {
				err = nil
				wasResponse = false
			} else {
				return nil, err
			}
		}

		if wasResponse == false {
			err = response.Body.Close() // Since we are in a loop we should close without defer
			if err != nil {
				return nil, err
			}
			time.Sleep(waitTime)
		}
	}

	if i == maxRetries {
		return nil, errors.New("max retries exceeded")
	}

	return response, nil
}

func checkHTTPResponse(response *http.Response) error {
	badResponseStatusCodes := []int{
		http.StatusInternalServerError,
		http.StatusNotImplemented,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	}

	if isInside(response.StatusCode, badResponseStatusCodes) {
		return errBadHTTPResponse
	}

	return nil
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
