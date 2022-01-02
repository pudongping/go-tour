package timer

import (
	"time"
)

// GetNowTime 获取本地时间的 Time 对象
func GetNowTime() time.Time {
	// 设置时区
	location, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(location)
}

// GetCalculateTime 时间推算
func GetCalculateTime(currentTimer time.Time, d string) (time.Time, error) {
	duration, err := time.ParseDuration(d)
	if err != nil {
		return time.Time{}, err
	}

	return currentTimer.Add(duration), nil
}
