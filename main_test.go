package main

import (
	"testing"
)

func Test_e2e_Success(t *testing.T) {
	err := parseFile(inputFile)

	if err != nil {
		t.Errorf("Expected no error, Got: %s", err)
	}
}

// func Test_e2e_Success(t *testing.T) {
// 	Test image creation
// }
func Test_e2e_WrongLengthInput(t *testing.T) {
	err := parseFile("inputs/wrong_length_input.txt")

	expected := "Incorrect Asset ID length found on row 1"

	if err.Error() != expected {
		t.Errorf("Expected %s, Got: %s", expected, err.Error())
	}
}

func Test_e2e_NotIntegerInput(t *testing.T) {
	err := parseFile("inputs/not_int_input.txt")

	expected := `strconv.Atoi: parsing "b": invalid syntax`

	if err.Error() != expected {
		t.Errorf("Expected %s, Got: %s", expected, err.Error())
	}
}

func Test_asset_createPng_Success(t *testing.T) {
	a := asset{
		id:       []int{1, 3, 3, 7},
		checksum: []int{5, 6},
		encoding: map[int]string{
			1: "11010101",
			2: "11110101",
			3: "01000100",
			4: "11010001",
			5: "11010001",
			6: "01000110",
		},
	}

	err := a.createPng()

	if err != nil {
		t.Errorf("Expected no error, Got: %s", err)
	}
}

func Test_asset_encodingPattern(t *testing.T) {
	a := asset{
		id:       []int{1, 3, 3, 7},
		checksum: []int{5, 6},
		encoding: map[int]string{
			1: "11010101",
			2: "11110101",
			3: "01000100",
			4: "11010001",
			5: "11010001",
			6: "01000110",
		},
	}

	pattern := a.encodingPattern()
	expected := "110101011111010101000100110100011101000101000110"

	if pattern != expected {
		t.Errorf("Expected %s, Got: %s", expected, pattern)
	}
}

func Test_asset_generateChecksum(t *testing.T) {
	a := asset{
		id:       []int{1, 3, 3, 7},
		checksum: make([]int, checksumLength),
		encoding: make(map[int]string),
	}

	checksum := a.generateChecksum()
	expected := 56

	if checksum != expected {
		t.Errorf("Expected %d, Got: %d", expected, checksum)
	}
}

func Test_asset_generateChecksum_Alternative(t *testing.T) {
	a := asset{
		id:       []int{2, 6, 7, 4},
		checksum: make([]int, checksumLength),
		encoding: make(map[int]string),
	}

	checksum := a.generateChecksum()
	expected := 9

	if checksum != expected {
		t.Errorf("Expected %d, Got: %d", expected, checksum)
	}
}

func Test_asset_setChecksum(t *testing.T) {
	a := asset{
		id:       []int{1, 3, 3, 7},
		checksum: make([]int, checksumLength),
		encoding: make(map[int]string),
	}

	_ = a.setChecksum()
	expected := []int{5, 6}

	if !checksumEq(a.checksum, expected) {
		t.Errorf("Expected %v, Got: %v", expected, a.checksum)
	}
}

func Test_asset_setChecksum_Alternative(t *testing.T) {
	a := asset{
		id:       []int{2, 6, 7, 4},
		checksum: make([]int, checksumLength),
		encoding: make(map[int]string),
	}

	_ = a.setChecksum()
	expected := []int{0, 9}

	if !checksumEq(a.checksum, expected) {
		t.Errorf("Expected %v, Got: %v", expected, a.checksum)
	}
}

func Test_asset_setEncoding_Success(t *testing.T) {
	a := asset{
		id:       []int{1, 3, 3, 7},
		checksum: []int{5, 6},
		encoding: make(map[int]string),
	}

	_ = a.setEncoding()
	expected := map[int]string{
		1: "11010101",
		2: "11110101",
		3: "01000100",
		4: "11010001",
		5: "11010001",
		6: "01000110",
	}

	if !encodingEq(a.encoding, expected) {
		t.Errorf("Expected %v, Got: %v", expected, a.encoding)
	}
}

func checksumEq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func encodingEq(a, b map[int]string) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if w, ok := b[k]; !ok || v != w {
			return false
		}
	}

	return true
}
