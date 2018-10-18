package util

import (
	"strings"

	"github.com/shopspring/decimal"
)

// Round a float, it will process decimal well like 0.00065123
func Round(v decimal.Decimal, places int32) decimal.Decimal {
	s := v.String()
	i := strings.Index(s, ".")

	if i == -1 {
		return v
	}

	if v.Cmp(decimal.NewFromFloat(1)) >= 0 || v.Cmp(decimal.NewFromFloat(-1)) <= 0 {
		return v.Round(places)
	}

	z := findNoneZero(s)

	// all zero?
	if z == -1 {
		return v.Round(places)
	}

	return v.Round(places + int32(z))
}

func findNoneZero(s string) int {
	l := len(s)
	b := -1

	for i := 0; i < l; i++ {
		if s[i:i+1] == "." {
			b = 0
		}
		if s[i:i+1] == "0" && b >= 0 {
			b++
		}
		if s[i:i+1] != "." && s[i:i+1] != "0" && b >= 0 {
			return b
		}
	}

	return b
}
