package utils

import "time"

func SecondsToTime(s int64) (int64, int64, int64) {
	hours := s / 3600
	s %= 3600
	min := s / 60
	seconds := s % 60

	return hours, min, seconds
}

func SetDate(src, dst time.Time) time.Time {
	return time.Date(src.Year(), src.Month(), src.Day(), dst.Hour(), dst.Minute(), dst.Second(), 0, dst.Location())
}

func SetTime(src, dst time.Time) time.Time {
	return SetDate(dst, src)
}

func SetDateTs(src, dst int64) int64 {
	return SetDate(time.Unix(src, 0), time.Unix(dst, 0)).Unix()
}

func SetTimeTs(src, dst int64) int64 {
	return SetTime(time.Unix(src, 0), time.Unix(dst, 0)).Unix()
}
