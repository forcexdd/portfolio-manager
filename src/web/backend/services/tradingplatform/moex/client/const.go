package client

import "time"

const (
	baseUrl    = "https://iss.moex.com/iss/"
	language   = "en"
	maxRetries = 10
	waitTime   = 1 * time.Second
)
