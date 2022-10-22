package utils

import (
	"fmt"
)

func FormatRank(rank int) (output string) {
	var suffix string
	var integerValue int
	var decimalValue int
	if rank >= 1000 {
		suffix = "k"
		integerValue = rank / 1000
		decimalValue = (rank - (integerValue * 1000)) / 100
	} else if rank >= 1_000_000 {
		suffix = "m"
		integerValue = rank / 1_000_000
		decimalValue = (rank - (integerValue * 1_000_000)) / 100_000
	} else {
		return fmt.Sprintf("#%d", rank)
	}
	if rank >= 10_000 {
		decimalValue = 0
	}
	if decimalValue == 0 {
		return fmt.Sprintf("#%d%s", integerValue, suffix)
	} else {
		return fmt.Sprintf("#%d.%d%s", integerValue, decimalValue, suffix)
	}
}
