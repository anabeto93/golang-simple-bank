package util

import (
	"math/rand"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max -min + 1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var bytes = make([]byte, n)
	for i := range bytes {
		bytes[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(bytes)
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD", "GHS", "NGN"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}