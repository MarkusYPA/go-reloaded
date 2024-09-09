// Package piscinefuncs includes functions created during the Go piscine at Grit:Lab
package piscinefuncs

// StrLen returns the length of string s
func StrLen(s string) int {
	length := 0
	for range s {
		length++
	}
	return length
}

// ToLower puts letters in a string to lowercase
func ToLower(s string) string {
	runes := []rune(s)
	for i := 0; i < StrLen(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			runes[i] = rune(s[i]) + 32
		}
	}
	return string(runes)
}

// ToUpper puts letters in a string to uppercase
func ToUpper(s string) string {
	runes := []rune(s)
	for i := 0; i < StrLen(s); i++ {
		if s[i] >= 'a' && s[i] <= 'z' {
			runes[i] = rune(s[i]) - 32
		}
	}
	return string(runes)
}

// Capitalize makes the first letters of words in a string uppercase and the subsequent letters lowercase
func Capitalize(s string) string {
	rns := []rune(s)
	isLetter := false
	isSmall := false
	isCapital := false
	prevIsLetter := false

	for i := 0; i < len(rns); i++ {
		if (rns[i] >= '0' && rns[i] <= '9') || (rns[i] >= 'A' && rns[i] <= 'Z') || (rns[i] >= 'a' && rns[i] <= 'z') {
			isLetter = true

			if rns[i] >= 'a' && rns[i] <= 'z' {
				isSmall = true
			}

			if rns[i] >= 'A' && rns[i] <= 'Z' {
				isCapital = true
			}
		}

		if !prevIsLetter && isSmall {
			rns[i] = rns[i] - 32
		}

		if prevIsLetter && isCapital {
			rns[i] = rns[i] + 32
		}

		prevIsLetter = isLetter
		isLetter = false
		isCapital = false
		isSmall = false
	}

	return string(rns)
}

// ToDec converts string s in base b to int
func ToDec(s, b string) int {
	num := 0
	mult := 1

	if s[0] == '-' || s[0] == '+' {
		if s[0] == '-' {
			mult = -1
		}
		s = s[1:]
	}

	for i := len(s) - 1; i >= 0; i-- {
		for j := 0; j < len(b); j++ {
			if s[i] == b[j] {
				toAdd := pawa(len(b), len(s)-1-i) * j
				num += toAdd
			}
		}
	}

	return num * mult
}

// ToBase converts a number n to base b
func ToBase(n int, b string) string {
	result := ""
	mult := ""
	minInt := false

	if n == -9223372036854775808 {
		n = 9223372036854775807
		mult = "-"
		minInt = true
	}

	if n < 0 {
		mult = "-"
		n *= -1
	}

	if n == 0 {
		return string(b[0])
	}

	for n > 0 {
		index := n % len(b)
		if minInt {
			index++
			if len(b) == index {
				index = 0
			}
		}
		result = string(b[index]) + result
		n /= len(b)
	}

	return mult + result
}

// pawa returns the number n to the power p
func pawa(n, p int) int {
	res := 1

	for i := 0; i < p; i++ {
		res = res * n
	}

	return res
}
