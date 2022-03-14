package collisionresolver

import "math"

func myAtoi(s string) int {
	multiplier := 1
	result := 1

	if len(s) == 0 {
		return 0
	}

	firstToken := string(s[0])
	if firstToken == "-" {
		multiplier = -1
		s = s[:1]
	} else if firstToken == "+" {
		multiplier = 1
		s = s[:1]
	}

	for i, c := range s {
		n := 0
		token := string(c)
		if token == "0" {
			n = 0
		} else if token == "1" {
			n = 1
		} else if token == "2" {
			n = 2
		} else if token == "3" {
			n = 3
		} else if token == "4" {
			n = 4
		} else if token == "5" {
			n = 5
		} else if token == "6" {
			n = 6
		} else if token == "7" {
			n = 7
		} else if token == "8" {
			n = 8
		} else if token == "9" {
			n = 9
		}
		result += int(math.Pow(10, float64(len(s)-i-1))) * n
	}
	result *= multiplier
	return result
}
