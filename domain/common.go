package domain

import (
	"fmt"
	"strings"
	"time"
)

// CustomTime 自定义时间类型，用于处理时间字符串的解析
type CustomTime time.Time

// UnmarshalJSON 实现json.Unmarshaler接口
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	// 移除引号
	s := strings.Trim(string(data), "\"")
	if s == "null" || s == "" {
		*ct = CustomTime(time.Time{})
		return nil
	}

	// 尝试解析多种时间格式
	t, err := parseTime(s)
	if err != nil {
		return err
	}

	*ct = CustomTime(t)
	return nil
}

// MarshalJSON 实现json.Marshaler接口
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	t := time.Time(ct)
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", t.Format(time.RFC3339))), nil
}

// String 实现fmt.Stringer接口
func (ct CustomTime) String() string {
	return time.Time(ct).String()
}

// Time 返回time.Time类型
func (ct CustomTime) Time() time.Time {
	return time.Time(ct)
}

func (ct *CustomTime) UnmarshalParam(param string) error {
	if param == "" {
		*ct = CustomTime(time.Time{})
		return nil
	}

	t, err := parseTime(param)
	if err != nil {
		return err
	}

	*ct = CustomTime(t)
	return nil
}

// parseTime 辅助函数，尝试多种时间格式解析
func parseTime(timeStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05-0700",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	var err error
	loc := time.Local
	for _, format := range formats {
		t, parseErr := time.ParseInLocation(format, timeStr, loc)
		if parseErr == nil {
			return t, nil
		}
		err = parseErr
	}

	return time.Time{}, fmt.Errorf("无法解析时间字符串 '%s': %v", timeStr, err)
}
