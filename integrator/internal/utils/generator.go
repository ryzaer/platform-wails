package utils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	alphabetCode      = "JA2BC3DE4FG5HKLMNQR6TU7VW8XY9ZP"
	defaultCodeLength = 7
)

func GenerateElsanaCode() (string, error) {

	return randomString(
		alphabetCode,
		defaultCodeLength,

		MaxRepeat(2),
		MaxSequential(4),
		MinLetter(2),
		MinNumber(2),
	)
}

func GenerateElsanaTenantCode() (string, error) {
	return randomString(
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		3,
		MaxRepeat(2),
	)
}

func GenerateCode(length int) (string, error) {

	if length <= 3 {
		length = defaultCodeLength
	}

	return randomString(
		alphabetCode,
		length,
		MaxRepeat(2),
		MaxSequential(4),
		MinLetter(2),
		MinNumber(2),
	)
}

func randomString(alphabet string, length int, rules ...Rule) (string, error) {

	max := big.NewInt(int64(len(alphabet)))

	for {

		var b strings.Builder
		b.Grow(length)

		for i := 0; i < length; i++ {

			n, err := rand.Int(rand.Reader, max)
			if err != nil {
				return "", err
			}

			b.WriteByte(alphabet[n.Int64()])
		}

		code := b.String()

		if validate(code, rules...) {
			return code, nil
		}
	}
}

// ===============>
// generator rules
// ===============>

type Rule func(string) bool

func validate(code string, rules ...Rule) bool {

	for _, rule := range rules {

		if !rule(code) {
			return false
		}
	}

	return true
}

func MinLetter(min int) Rule {

	return func(code string) bool {

		count := 0

		for _, c := range code {

			if c >= 'A' && c <= 'Z' {
				count++
			}
		}

		return count >= min
	}
}

func MaxRepeat(max int) Rule {
	return func(code string) bool {
		if len(code) == 0 {
			return true
		}

		run := 1

		for i := 1; i < len(code); i++ {
			if code[i] == code[i-1] {
				run++
			} else {
				run = 1
			}

			if run > max {
				return false
			}
		}

		return true
	}
}

func MinNumber(min int) Rule {
	return func(code string) bool {
		count := 0
		for _, c := range code {
			if c >= '0' && c <= '9' {
				count++
			}
		}

		return count >= min
	}
}

func MaxSequential(max int) Rule {

	return func(code string) bool {
		run := 1

		for i := 1; i < len(code); i++ {
			if code[i] == code[i-1]+1 {
				run++
			} else {
				run = 1
			}

			if run > max {
				return false
			}
		}

		return true
	}
}

var ElsanaRules = []Rule{
	MaxRepeat(2),
	MaxSequential(4),
	MinLetter(2),
	MinNumber(2),
}
