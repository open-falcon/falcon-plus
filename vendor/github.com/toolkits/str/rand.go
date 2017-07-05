package str

import (
	"fmt"
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var digits = []rune("0123456789")

const size = 62

func RandSeq(n int) string {

	seed := time.Now().UnixNano()
	rand.Seed(seed)

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(size)]
	}

	return fmt.Sprintf("%d_%s", seed, string(b))

}

func RandStr(n int) string {

	seed := time.Now().UnixNano()
	rand.Seed(seed)

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(size)]
	}

	return fmt.Sprintf("%s", string(b))
}

func RandDigits(n int) string {

	seed := time.Now().UnixNano()
	rand.Seed(seed)

	b := make([]rune, n)
	for i := range b {
		b[i] = digits[rand.Intn(10)]
	}

	return fmt.Sprintf("%s", string(b))
}
