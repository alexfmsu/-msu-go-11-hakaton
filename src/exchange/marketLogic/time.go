package marketLogic

import (
	"strconv"
	"time"
)

var (
	StartTime int64
	StartTs   int64
)

func init() {
	StartTime = time.Now().Unix()
}

func getTime(t string) (time.Time, error) {
	tmp, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return time.Now(), err
	}
	tmp -= StartTs
	seconds := tmp % 100
	tmp /= 100
	minutes := tmp % 100
	tmp /= 100
	hours := tmp % 100
	seconds = seconds + 60*minutes + 3600*hours
	return time.Unix(StartTime+seconds, 0), nil
}
