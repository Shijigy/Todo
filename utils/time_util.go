package utils

import (
	"fmt"
	"time"
)

// FormatDate 将时间格式化为字符串（yyyy-MM-dd HH:mm:ss）
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// ParseDate 解析时间字符串为 time.Time
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", dateStr)
}

// GetCurrentTime 获取当前时间
func GetCurrentTime() string {
	return FormatDate(time.Now())
}

// AddDaysToCurrentDate 给当前日期加上指定天数
func AddDaysToCurrentDate(days int) string {
	return FormatDate(time.Now().AddDate(0, 0, days))
}

// CalculateTimeDifference 计算两个时间的差值（返回小时、分钟、秒）
func CalculateTimeDifference(startTime, endTime time.Time) string {
	duration := endTime.Sub(startTime)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%d hours, %d minutes, %d seconds", hours, minutes, seconds)
}
