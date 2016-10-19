package main

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
    letterIdxBits = 6 
    letterIdxMask = 1<<letterIdxBits - 1
    letterIdxMax  = 63 / letterIdxBits
)
var src = rand.NewSource(time.Now().UnixNano())

// http://stackoverflow.com/a/31832326
// Random string of 20 characters
func RandomString() string {
    b := make([]byte, 20)
    for i, cache, remain := 20-1, src.Int63(), letterIdxMax; i >= 0; {
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
    return string(b)
}