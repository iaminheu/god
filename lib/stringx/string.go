package stringx

import (
	"errors"
	"git.zc0901.com/go/god/lib/lang"
	"math"
	"strconv"
	"strings"
)

var (
	ErrInvalidStartPosition = errors.New("起始位置无效")
	ErrInvalidStopPosition  = errors.New("结束位置无效")
	byteSizes               = []string{"Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
)

func Contains(list []string, str string) bool {
	for _, each := range list {
		if each == str {
			return true
		}
	}

	return false
}

func Filter(s string, filter func(r rune) bool) string {
	var n int
	chars := []rune(s)
	for i, x := range chars {
		if n < i {
			chars[n] = x
		}
		if !filter(x) {
			n++
		}
	}

	return string(chars[:n])
}

func HasEmpty(args ...string) bool {
	for _, arg := range args {
		if len(arg) == 0 {
			return true
		}
	}

	return false
}

func NotEmpty(args ...string) bool {
	return !HasEmpty(args...)
}

func Remove(strings []string, strs ...string) []string {
	out := append([]string(nil), strings...)

	for _, str := range strs {
		var n int
		for _, v := range out {
			if v != str {
				out[n] = v
				n++
			}
		}
		out = out[:n]
	}

	return out
}

func Reverse(s string) string {
	runes := []rune(s)

	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}

	return string(runes)
}

// Substr returns runes between start and stop [start, stop) regardless of the chars are ascii or utf8
func Substr(str string, start int, stop int) (string, error) {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return "", ErrInvalidStartPosition
	}

	if stop < 0 || stop > length {
		return "", ErrInvalidStopPosition
	}

	return string(rs[start:stop]), nil
}

func TakeOne(valid, or string) string {
	if len(valid) > 0 {
		return valid
	} else {
		return or
	}
}

func TakeWithPriority(fns ...func() string) string {
	for _, fn := range fns {
		val := fn()
		if len(val) > 0 {
			return val
		}
	}

	return ""
}

func Union(first, second []string) []string {
	set := make(map[string]lang.PlaceholderType)

	for _, each := range first {
		set[each] = lang.Placeholder
	}
	for _, each := range second {
		set[each] = lang.Placeholder
	}

	merged := make([]string, 0, len(set))
	for k := range set {
		merged = append(merged, k)
	}

	return merged
}

// ToFixed truncates float64 type to a particular precision in string.
func ToFixed(n float64, precision int) string {
	s := strconv.FormatFloat(n, 'f', precision, 64)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

// 字节转带单位的大小
func BytesToSize(bytes uint) (size string) {
	if bytes == 0 {
		return "0"
	}
	i := math.Floor(math.Log(float64(bytes)) / math.Log(1024))
	total := float64(bytes) / math.Pow(1024, i)
	precision := 0
	if total < 10 && i > 0 {
		precision = 1
	}
	return ToFixed(total, precision) + " " + byteSizes[int(i)]
}

func Bytes2Size(bytes uint64) (size string) {
	base := math.Log(float64(bytes)) / math.Log(1024)
	getSize := Round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	getSuffix := byteSizes[int(math.Floor(base))]
	return strconv.FormatFloat(getSize, 'f', -1, 64) + " " + getSuffix
}
