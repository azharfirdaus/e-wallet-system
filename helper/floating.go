package helper

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const maxInt64 = 1<<63 - 1

func ParseAmountMinorUnit(amount string) (int64, error) {
	amount = strings.TrimSpace(amount)
	if amount == "" {
		return 0, errors.New("amount is required")
	}

	parts := strings.Split(amount, ".")
	if len(parts) != 2 {
		return 0, errors.New("amount must have exactly two decimal digits")
	}

	wholePart := parts[0]
	if wholePart == "" {
		return 0, errors.New("amount must be numerical")
	}
	for _, digit := range wholePart {
		if digit < '0' || digit > '9' {
			return 0, errors.New("amount must be numerical")
		}
	}

	fractionPart := parts[1]
	if len(fractionPart) != 2 {
		return 0, errors.New("amount must have exactly two decimal digits")
	}
	for _, digit := range fractionPart {
		if digit < '0' || digit > '9' {
			return 0, errors.New("amount must be numerical")
		}
	}

	wholeAmount, err := strconv.ParseInt(wholePart, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("amount is too large: %w", err)
	}

	fractionAmount, err := strconv.ParseInt(fractionPart, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("amount is too large: %w", err)
	}

	if wholeAmount > (maxInt64-fractionAmount)/100 {
		return 0, errors.New("amount is too large")
	}

	return wholeAmount*100 + fractionAmount, nil
}
