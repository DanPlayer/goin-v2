package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// RandomInt64 根据区间产生随机数
func RandomInt64(min, max int64) int64 {
	rand.Seed(time.Now().Unix())
	return rand.Int63n(max-min) + min
}

// RandomFloat64 根据区间产生随机数
func RandomFloat64(min, max float64) float64 {
	rand.Seed(time.Now().Unix())
	return min + rand.Float64() * (max - min)
}

// Decimal 保留float64两位小数
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}