package utils

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// FormatMoney format money value to string representation
//
//	FormatMoney(1545285.25) => "1 545 285.25"
//	FormatMoney(1) => "1.00"
//	FormatMoney(0) => "0.00"
func FormatMoney[T float32 | float64](m T) string {
	var result strings.Builder

	items := strings.Split(fmt.Sprintf("%.2f", m), ".")

	items[0] = ReverseString(items[0])
	items[1] = ReverseString(items[1])

	result.WriteString(items[1])
	result.WriteRune('.')

	for i, item := range items[0] {
		if i > 0 && i%3 == 0 {
			result.WriteRune(' ')
		}

		result.WriteRune(item)
	}

	return ReverseString(result.String())
}

// ReverseString reverse string ordering
func ReverseString(s string) string {
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

// FormatDuration format time.Duration to clock presentation like 00:00 or 00:00:00
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

// ToFixed returns T value with specified precission: ToFixed(10.24353, 1) => 10.2
//
//	ToFixed(10.24535, 2) => 10.2
func ToFixed[T float32 | float64](num T, precision int) T {
	output := math.Pow(10, float64(precision))

	return T(float64(round(float64(num)*output)) / output)
}
