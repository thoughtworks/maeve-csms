package ocpp

import (
	"fmt"
	"regexp"
	"strings"
)

var emaidPattern = regexp.MustCompile("([A-Za-z]{2})(-?)([A-Za-z0-9]{3})(-?)([A-Za-z0-9]{9})(-?)([A-Za-z0-9])?")

func NormalizeEmaid(emaid string) (string, error) {
	norm, err := normalizeEmaid(emaid)
	if err != nil {
		return "", err
	}
	digit := calculateEmaidCheckDigit(norm)
	if len(norm) > 14 {
		if digit != rune(norm[14]) {
			return "", fmt.Errorf("emaid check digit %c expected value %c", norm[14], digit)
		}
		return norm, nil
	}

	return fmt.Sprintf("%s%c", norm, digit), nil
}

func normalizeEmaid(emaid string) (string, error) {
	parts := emaidPattern.FindStringSubmatch(emaid)
	if len(parts) > 0 {
		str := fmt.Sprintf("%s%s%s%s", parts[1], parts[3], parts[5], parts[7])
		return strings.ToUpper(str), nil
	}
	return "", fmt.Errorf("emaid %s is invalid", emaid)
}

var reverseLookup = "0I9R3LCU6OFX----1JAS4MDV7PGY----2KBT5NEW8QHZ"

var lookupTable = [][]int{
	{0, 0, 0, 0}, // 0
	{0, 0, 0, 1}, // 1
	{0, 0, 0, 2}, // 2
	{0, 0, 1, 0}, // 3
	{0, 0, 1, 1}, // 4
	{0, 0, 1, 2}, // 5
	{0, 0, 2, 0}, // 6
	{0, 0, 2, 1}, // 7
	{0, 0, 2, 2}, // 8
	{0, 1, 0, 0}, // 9
	{},
	{},
	{},
	{},
	{},
	{},
	{},
	{0, 1, 0, 1}, // A
	{0, 1, 0, 2}, // B
	{0, 1, 1, 0}, // C
	{0, 1, 1, 1}, // D
	{0, 1, 1, 2}, // e
	{0, 1, 2, 0}, // F
	{0, 1, 2, 1}, // G
	{0, 1, 2, 2}, // H
	{1, 0, 0, 0}, // I
	{1, 0, 0, 1}, // J
	{1, 0, 0, 2}, // K
	{1, 0, 1, 0}, // L
	{1, 0, 1, 1}, // M
	{1, 0, 1, 2}, // N
	{1, 0, 2, 0}, // O
	{1, 0, 2, 1}, // P
	{1, 0, 2, 2}, // Q
	{1, 1, 0, 0}, // R
	{1, 1, 0, 1}, // S
	{1, 1, 0, 2}, // T
	{1, 1, 1, 0}, // U
	{1, 1, 1, 1}, // V
	{1, 1, 1, 2}, // W
	{1, 1, 2, 0}, // X
	{1, 1, 2, 1}, // Y
	{1, 1, 2, 2}, // Z
}

var p1 = [][]int{
	{0, 1, 1, 1},
	{1, 1, 1, 0},
	{1, 0, 0, 1},
	{0, 1, 1, 1},
	{1, 1, 1, 0},
	{1, 0, 0, 1},
	{0, 1, 1, 1},
	{1, 1, 1, 0},
	{1, 0, 0, 1},
	{0, 1, 1, 1},
	{1, 1, 1, 0},
	{1, 0, 0, 1},
	{0, 1, 1, 1},
	{1, 1, 1, 0},
	{1, 0, 0, 1},
}

var p2 = [][]int{
	{0, 1, 1, 2},
	{1, 2, 2, 2},
	{2, 2, 2, 0},
	{2, 0, 0, 2},
	{0, 2, 2, 1},
	{2, 1, 1, 1},
	{1, 1, 1, 0},
	{1, 0, 0, 1},
	{0, 1, 1, 2},
	{1, 2, 2, 2},
	{2, 2, 2, 0},
	{2, 0, 0, 2},
	{0, 2, 2, 1},
	{2, 1, 1, 1},
	{1, 1, 1, 0},
}

// calculateEmaidCheckDigit takes a normalized eMAID and calculates the expected check digit.
// The algorithm is described in:
// http://www.ochp.eu/wp-content/uploads/2014/02/E-Mobility-IDs_EVCOID_Check-Digit-Calculation_Explanation.pdf
func calculateEmaidCheckDigit(emaid string) rune {
	var c1, c2, c3, c4 int
	for i := 0; i < 14; i++ {
		matrix := lookupTable[emaid[i]-'0']
		c1 += matrix[0]*p1[i][0] + matrix[1]*p1[i][2]
		c2 += matrix[0]*p1[i][1] + matrix[1]*p1[i][3]
		c3 += matrix[2]*p2[i][0] + matrix[3]*p2[i][2]
		c4 += matrix[2]*p2[i][1] + matrix[3]*p2[i][3]
	}

	c1 = c1 % 2
	c2 = c2 % 2
	c3 = c3 % 3
	c4 = c4 % 3

	q1 := c1
	q2 := c2
	var r1, r2 int

	switch c4 {
	case 0:
		r1 = 0
	case 1:
		r1 = 2
	case 2:
		r1 = 1
	}

	switch c3 + r1 {
	case 0:
		r2 = 0
	case 1:
		r2 = 2
	case 2:
		r2 = 1
	case 3:
		r2 = 0
	case 4:
		r2 = 2
	}

	return rune(reverseLookup[q1+q2*2+r1*4+r2*16])
}
