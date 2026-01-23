package utils

import "time"

// SecondsToTime convert UNIXTIME to hours, minutes and seconds
func SecondsToTime(s int64) (int64, int64, int64) {
	hours := s / 3600
	s %= 3600
	min := s / 60
	seconds := s % 60

	return hours, min, seconds
}

// SetDate returns new Time with year, month and day from src and hour, minute, second and location from dst
func SetDate(src, dst time.Time) time.Time {
	return time.Date(src.Year(), src.Month(), src.Day(), dst.Hour(), dst.Minute(), dst.Second(), 0, dst.Location())
}

// SetTime returns new Time with year, month and day from dst and hour, minute, second and location from src
func SetTime(src, dst time.Time) time.Time {
	return SetDate(dst, src)
}

// SetDateTs see SetDate
func SetDateTs(src, dst int64) int64 {
	return SetDate(time.Unix(src, 0), time.Unix(dst, 0)).Unix()
}

// SetTimeTs see SetTime
func SetTimeTs(src, dst int64) int64 {
	return SetTime(time.Unix(src, 0), time.Unix(dst, 0)).Unix()
}
