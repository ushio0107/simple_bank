package util

import (
	"math/rand"
	"strings"
	"time"
)

var r *rand.Rand

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}
func RandomInt(min, max int64) int64 {
	return min + r.Int63n(max-min+1)
}

// RandomString returns a random string with a fixed length.
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[r.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner returns a random owner name.
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney returns a random money amount.
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency returns a random currency string, included EUR, USD and CAD.
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[r.Intn(n)]
}

// RandomEmail returns a random valid email string.
func RandomEmail() string {
	return RandomString(6) + "@gmail.com"
}
