package oauth2

import "time"

// GetNowInEpochTime returns the epoch time normally used in JWT's
func GetNowInEpochTime() int64 {
	return time.Now().UnixNano() / 1000000000
}
