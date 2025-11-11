package utils

import (
	"regexp"
	"strconv"
)

// ValidateCPF validates if a CPF is valid according to Brazilian rules
func ValidateCPF(cpf string) bool {
	// Remove dots and hyphens
	re := regexp.MustCompile(`[^\d]`)
	cpf = re.ReplaceAllString(cpf, "")

	// Check if CPF has 11 digits
	if len(cpf) != 11 {
		return false
	}

	// Check if all digits are the same (invalid CPFs)
	if cpf == "00000000000" || cpf == "11111111111" || cpf == "22222222222" ||
		cpf == "33333333333" || cpf == "44444444444" || cpf == "55555555555" ||
		cpf == "66666666666" || cpf == "77777777777" || cpf == "88888888888" ||
		cpf == "99999999999" {
		return false
	}

	// Calculate first verification digit
	sum := 0
	for i := 0; i < 9; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (10 - i)
	}
	remainder := sum % 11
	firstDigit := 0
	if remainder >= 2 {
		firstDigit = 11 - remainder
	}

	// Check first verification digit
	digit10, _ := strconv.Atoi(string(cpf[9]))
	if digit10 != firstDigit {
		return false
	}

	// Calculate second verification digit
	sum = 0
	for i := 0; i < 10; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (11 - i)
	}
	remainder = sum % 11
	secondDigit := 0
	if remainder >= 2 {
		secondDigit = 11 - remainder
	}

	// Check second verification digit
	digit11, _ := strconv.Atoi(string(cpf[10]))
	return digit11 == secondDigit
}

// FormatCPF formats a CPF string with dots and hyphen
func FormatCPF(cpf string) string {
	// Remove existing formatting
	re := regexp.MustCompile(`[^\d]`)
	cpf = re.ReplaceAllString(cpf, "")

	if len(cpf) != 11 {
		return cpf
	}

	return cpf[:3] + "." + cpf[3:6] + "." + cpf[6:9] + "-" + cpf[9:]
}

// CleanCPF removes formatting from CPF (dots and hyphens)
func CleanCPF(cpf string) string {
	re := regexp.MustCompile(`[^\d]`)
	return re.ReplaceAllString(cpf, "")
}