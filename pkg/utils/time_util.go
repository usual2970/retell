package utils

import (
	"fmt"
	"time"
)

// GetUTCRangeFromStartOfLastMonthToNowInLocation 获取“指定时区”中从上月第一天零点到当前时间的 UTC 范围
// offsetHours：时区偏移（如 8 表示 UTC+8）
func GetUTCRangeFromStartOfLastMonthToNowInLocation(offsetHours int) (startUTC, endUTC int64) {
	// 构造时区
	loc := time.FixedZone(fmt.Sprintf("UTC%+d", offsetHours), offsetHours*3600)

	// 当前时间（该时区）
	now := time.Now().UTC().In(loc)

	// 计算“上个月”的年份与月份
	year, month := now.Year(), now.Month()
	if month == time.January {
		year -= 1
		month = time.December
	} else {
		month -= 1
	}

	// 得到该时区下的“上个月第一天零点”
	startLocal := time.Date(year, month, 1, 0, 0, 0, 0, loc)

	// UTC 时间戳
	startUTC = startLocal.Unix()
	endUTC = now.UTC().Unix()

	return startUTC, endUTC
}

// GetUTCTodayRangeByOffset 根据时区偏移获取“今天”的 UTC 范围（单位：秒）
func GetUTCTodayRangeByOffset(offsetHours int) (int64, int64) {
	loc := time.FixedZone(fmt.Sprintf("UTC%+d", offsetHours), offsetHours*3600)
	// 当前时刻（按指定时区）
	now := time.Now().UTC().In(loc)
	// 自然日起点（本地时间）
	startLocal := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	return startLocal.Unix(), now.UTC().Unix()
}

// GetUTCYesterdayRangeByOffset 根据时区偏移获取“昨天”的 UTC 范围(单位: 秒)
func GetUTCYesterdayRangeByOffset(offsetHours int) (startUTC, endUTC int64) {
	loc := time.FixedZone(fmt.Sprintf("UTC%+d", offsetHours), offsetHours*3600)
	// 当前时间转换为指定时区
	now := time.Now().UTC().In(loc)

	// 构造“昨天”的本地时间段（00:00:00 ~ 23:59:59）
	yesterday := now.AddDate(0, 0, -1)
	startLocal := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, loc)
	endLocal := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, loc)

	// 转换为 UTC 时间戳
	return startLocal.Unix(), endLocal.Unix()
}

// GetUTCThisMonthRangeByOffset 根据时间偏移获取“本月”的 UTC 范围
func GetUTCThisMonthRangeByOffset(offsetHours int) (int64, int64) {
	loc := getLocationByOffset(offsetHours)
	now := time.Now().UTC().In(loc)

	// 本月第一天
	startLocal := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)

	// 下月第一天 -1 秒
	firstOfNextMonth := startLocal.AddDate(0, 1, 0)
	endLocal := firstOfNextMonth.Add(-time.Second)

	return startLocal.Unix(), endLocal.Unix()
}

// GetUTCLastMonthRangeByOffset 根据时间偏移获取“上月”的 UTC 范围
func GetUTCLastMonthRangeByOffset(offsetHours int) (int64, int64) {
	loc := getLocationByOffset(offsetHours)
	now := time.Now().UTC().In(loc)
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)

	lastMonthStart := currentMonthStart.AddDate(0, -1, 0)
	lastMonthEnd := currentMonthStart.Add(-time.Second) // One second before current month starts

	return lastMonthStart.Unix(), lastMonthEnd.Unix()
}

// getLocationByOffset returns a *time.Location based on hour offset from UTC.
func getLocationByOffset(offsetHours int) *time.Location {
	return time.FixedZone(fmt.Sprintf("UTC%+d", offsetHours), offsetHours*3600)
}

// GetLocalTimeString 根据时间戳和时区偏移，返回该时区下的本地时间字符串（原始时间精度保持不变）
func GetLocalTimeString(timestamp int64, offsetHours int) string {
	loc := time.FixedZone(fmt.Sprintf("UTC%+d", offsetHours), offsetHours*3600)
	localTime := time.Unix(timestamp, 0).In(loc)
	return localTime.Format("2006-01-02 15:04:05")
}

// GetLocalTimeStringMonth 根据时间戳和时区偏移，返回该时区下的本地时间字符串，返回月（原始时间精度保持不变）
func GetLocalTimeStringMonth(timestamp int64, offsetHours int) string {
	loc := time.FixedZone(fmt.Sprintf("UTC%+d", offsetHours), offsetHours*3600)
	localTime := time.Unix(timestamp, 0).In(loc)
	return localTime.Format("2006-01")
}

// AdjustTimeRangeToNaturalDay 查询的是UTC的时间范围
func AdjustTimeRangeToNaturalDay(utcStart, utcEnd int64, offset int) (int64, int64) {
	loc := time.FixedZone(fmt.Sprintf("UTC%+d", offset), offset*3600)
	startLocal := time.Unix(utcStart, 0).In(loc)
	endLocal := time.Unix(utcEnd, 0).In(loc)

	startDay := time.Date(startLocal.Year(), startLocal.Month(), startLocal.Day(), 0, 0, 0, 0, loc)
	endDay := time.Date(endLocal.Year(), endLocal.Month(), endLocal.Day(), 23, 59, 59, 0, loc)

	return startDay.UTC().Unix(), endDay.UTC().Unix()

}

// GetUTCTimeRangeUntilNow 计算从指定时间（按 offset 时区，减去年/月/日）到当前时刻的 UTC 时间戳范围
func GetUTCTimeRangeUntilNow(offsetHours, years, months, days int) (int64, int64) {
	loc := time.FixedZone(fmt.Sprintf("UTC%+d", offsetHours), offsetHours*3600)

	// 当前时间（以指定时区为本地）
	now := time.Now().UTC().In(loc)

	// 开始时间（往前推）
	startLocal := now.AddDate(years, months, days)

	// 自然日起点（如当天 00:00:00）
	startOfDay := time.Date(startLocal.Year(), startLocal.Month(), startLocal.Day(), 0, 0, 0, 0, loc)

	// 结束时间（当前时刻）
	return startOfDay.Unix(), now.UTC().Unix()
}

// GetUTCWeekStart 获取一周的起始时间。
func GetUTCWeekStart(offset int) int64 {
	loc := time.FixedZone(fmt.Sprintf("UTC%+d", offset), offset*3600)
	now := time.Now().In(loc)
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, loc)
	return weekStart.UTC().Unix()
}

// GetUTCThisWeekRangeByOffset 根据时间偏移获取“本周”的 UTC 范围（周一 00:00:00 到 周日 23:59:59）
func GetUTCThisWeekRangeByOffset(offsetHours int) (int64, int64) {
	loc := getLocationByOffset(offsetHours)
	now := time.Now().UTC().In(loc)

	// 获取本地今天是周几（周一=1，周日=7）
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	// 本周一 00:00:00
	weekStartLocal := time.Date(now.Year(), now.Month(), now.Day()-(weekday-1), 0, 0, 0, 0, loc)
	// 本周日 23:59:59
	weekEndLocal := weekStartLocal.AddDate(0, 0, 7).Add(-time.Second)

	return weekStartLocal.UTC().Unix(), weekEndLocal.UTC().Unix()
}

// GetUTCLastWeekRangeByOffset 根据时间偏移获取“上周”的 UTC 范围（上周一 00:00:00 到 上周日 23:59:59）
func GetUTCLastWeekRangeByOffset(offsetHours int) (int64, int64) {
	loc := getLocationByOffset(offsetHours)
	now := time.Now().UTC().In(loc)

	// 获取本地今天是周几（周一=1，周日=7）
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	// 本周一 00:00:00
	thisWeekStartLocal := time.Date(now.Year(), now.Month(), now.Day()-(weekday-1), 0, 0, 0, 0, loc)
	// 上周一 00:00:00
	lastWeekStartLocal := thisWeekStartLocal.AddDate(0, 0, -7)
	// 上周日 23:59:59
	lastWeekEndLocal := thisWeekStartLocal.Add(-time.Second)

	return lastWeekStartLocal.UTC().Unix(), lastWeekEndLocal.UTC().Unix()
}
