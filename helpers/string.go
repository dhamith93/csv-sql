package helpers

import (
	"math/rand"
	"strconv"
	"time"
	"unsafe"
)

var src = rand.NewSource(time.Now().UnixNano())

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func RandSeq(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func IsIntegral(val string) bool {
	if !isNumeric(val) {
		return false
	}
	num, _ := strconv.ParseFloat(val, 64)
	return isIntegral(num)
}

func ConvertToIntString(val string) string {
	num, _ := strconv.ParseFloat(val, 64)
	intNum := int(num)
	return strconv.Itoa(intNum)
}

func isIntegral(val float64) bool {
	return val == float64(int(val))
}

func isNumeric(val string) bool {
	_, err := strconv.ParseFloat(val, 64)
	return err == nil
}
