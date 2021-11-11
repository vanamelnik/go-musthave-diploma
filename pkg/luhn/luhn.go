package luhn

import (
	"errors"
	"fmt"
	"strconv"
)

// Validate performs check the provided number by Luhn's algorithm.
// Returns false if the number is empty, the check is failed or the number contains any characters except digits.
func Validate(number string) bool {
	if len(number) == 0 {
		return false
	}
	digits := number[:(len(number) - 1)]
	checksum, err := strconv.Atoi(string(number[len(number)-1]))
	if err != nil {
		return false
	}
	calculatedSum, err := Checksum(digits)
	if err != nil {
		return false
	}

	return checksum == calculatedSum
}

// Calculate adds to the provided number a checksum digit, calculated by Luhn's algorithm.
// Returns an error if the number is empty or contains any characters except digits.
func Calculate(number string) (string, error) {
	checksum, err := Checksum(number)
	if err != nil {
		return "", err
	}

	return number + fmt.Sprint(checksum), nil
}

// Checksum calculates a checksum digit by Luhn's algorithm for the number provided.
// Returns an error if the number is empty or contains any characters except digits.
func Checksum(number string) (int, error) {
	if number == "" {
		return -1, errors.New("luhn: number is empty") // -1 is an extra protection if somebody won't check the error...
	}
	sum := 0
	for i, pos := len(number)-1, 0; i >= 0; i-- {
		pos++
		digit, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return -1, fmt.Errorf("luhn: %w", err)
		}
		if pos%2 == 1 {
			d2 := digit * 2
			digit = d2/10 + d2%10
		}
		sum += digit
	}
	return (sum * 9) % 10, nil
}
