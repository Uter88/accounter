package tools

import (
	"fmt"
	"math"
	"time"
)

func SecondsToTime(s int64) (int64, int64, int64) {
	hours := s / 3600
	s %= 3600
	min := s / 60
	seconds := s % 60

	return hours, min, seconds
}

func FormatDuration(d time.Duration, withSeconds bool) string {
	var (
		h, m, s = SecondsToTime(int64(d.Seconds()))
		format  = "%.2d:%.2d"
		args    = []any{h, m}
	)

	if withSeconds {
		format += ":%.2d"
		args = append(args, s)
	}

	return fmt.Sprintf(format, args...)
}

// round float64
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed[T float32 | float64](num T, precision int) T {
	output := math.Pow(10, float64(precision))

	return T(float64(round(float64(num)*output)) / output)
}
