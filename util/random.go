package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const ab = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of n length
func RandomString(n int) string {
	var sb strings.Builder
	k := len(ab)

	for i := 0; i < n; i++ {
		c := ab[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomCurrency() string {
	c := []string{USD, EUR, GBP}
	return c[RandomInt(0, 2)]
}

func RandomBalance() int64 {
	return RandomInt(0, 1000)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
