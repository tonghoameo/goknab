package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		ch := alphabet[rand.Intn(k)] // get random index character
		sb.WriteByte(ch)
	}
	return sb.String()
}

// random owner generates a random owner name

func RandomOwner() string {
	return RandomString(6)
}
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}
func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD}
	return currencies[rand.Intn(len(currencies))]
}
func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}
