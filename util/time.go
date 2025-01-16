package util

import "strconv"

func FormatTime(time int64) string {
	msSecond := time % 1000
	second := time / 1000
	minute := second / 60
	hour := minute / 60
	minute %= 60
	second %= 60
	ret := strconv.FormatInt(hour, 10) + ":" + strconv.FormatInt(minute, 10) + ":" + strconv.FormatInt(second, 10) + ":" + strconv.FormatInt(msSecond, 10)
	return ret
}

func FormatHumanTime(time int64) string {
	second := time / 1000
	minute := second / 60
	hour := minute / 60
	minute %= 60
	second %= 60
	ret := strconv.FormatInt(hour, 10) + ":" + strconv.FormatInt(minute, 10) + ":" + strconv.FormatInt(second, 10)
	return ret
}

func FormatHumanDisplayTime(time int64) string {
	second := time / 1000
	minute := second / 60
	hour := minute / 60
	minute %= 60
	second %= 60
	ret := ""
	if hour > 0 {
		ret += strconv.FormatInt(hour, 10) + "时"
	}
	if minute > 0 {
		ret += strconv.FormatInt(minute, 10) + "分"
	}
	ret += strconv.FormatInt(second, 10) + "秒"
	return ret
}
