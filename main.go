package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var inputFile string
var defaultInputFile = "inputs/test_input.txt"

// Represents the 4-digit numerical asset IDs
type id []int

// Represents the 2-digit numerical checksum
type checksum []int

// The key represent the character (number) on the lcd display.
// The value represents the encoded bit string
//
// map[int]string{ 1: "11010101", 2: "11010101" }
//
// Prefer storing encoding as a map to prevent need to concat a string continually
// plus provide flexibility for future requirements.
type encoding map[int]string

// An asset represents a "display" as mentioned in the specification
// This struct holds the id read from the input file which is then augmented
// with the generated checksum and encoding(s).
type asset struct {
	id
	checksum
	encoding
}

const (
	idLength       = 4
	checksumLength = 2
	pngWidth       = 256
	pngHeight      = 1
	offset         = 8 // Defines offset (given that bits 0–7 are reserved for other uses)
)

func init() {
	flag.StringVar(&inputFile, "inputFile", defaultInputFile, "Specifies dir of input file")
}

func main() {
	flag.Parse()

	err := parseFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}
}

func parseFile(inputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rowNumber := 0
	for scanner.Scan() {
		rowNumber++

		rowData := strings.Split(scanner.Text(), "")
		err = handleRow(rowData, rowNumber)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	log.Println("Done")
	return nil
}

func handleRow(rowData []string, rowNumber int) error {
	if len(rowData) != 4 {
		return fmt.Errorf("Incorrect Asset ID length found on row %d", rowNumber)
	}

	convertedRow, err := sliceAtoi(rowData)
	if err != nil {
		return err
	}

	a := asset{
		id:       convertedRow,
		checksum: make([]int, checksumLength),
		encoding: make(map[int]string),
	}

	err = a.setChecksum()
	if err != nil {
		return err
	}

	err = a.setEncoding()
	if err != nil {
		return err
	}

	img, err := a.buildImage()
	if err != nil {
		return err
	}

	err = a.persistToFile(img)
	if err != nil {
		return err
	}

	return nil
}

func (a *asset) idStr() string {
	return fmt.Sprintf("%d%d%d%d", a.id[0], a.id[1], a.id[2], a.id[3])
}

func (a *asset) setChecksum() error {
	var first int
	var second int
	var err error

	checksum := a.generateChecksum()
	checksumStr := strconv.Itoa(checksum)

	if len(checksumStr) == 2 {
		first, err = strconv.Atoi(string(checksumStr[0]))
		if err != nil {
			return err
		}

		second, err = strconv.Atoi(string(checksumStr[1]))
		if err != nil {
			return err
		}
	} else {
		first = 0
		second, err = strconv.Atoi(checksumStr)
		if err != nil {
			return err
		}
	}

	a.checksum[0] = first
	a.checksum[1] = second

	return nil
}

func (a *asset) setEncoding() error {
	for i, checksumDigit := range a.checksum {
		enc, err := encodingForDigit(checksumDigit)
		if err != nil {
			return err
		}

		a.encoding[i+1] = enc
	}

	for i, idDigit := range a.id {
		enc, err := encodingForDigit(idDigit)
		if err != nil {
			return err
		}

		a.encoding[i+1+checksumLength] = enc
	}

	return nil
}

func (a *asset) generateChecksum() int {
	return (a.id[0] + (10 * a.id[1]) + (100 * a.id[2]) + (1000 * a.id[3])) % 97
}

func (a *asset) encodingPattern() string {
	keys := make([]int, 0)
	for k := range a.encoding {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	values := make([]string, 0, len(a.encoding))
	for _, k := range keys {
		values = append(values, a.encoding[k])
	}

	return strings.Join(values, "")
}

func (a *asset) buildImage() (*image.NRGBA, error) {
	encoding := a.encodingPattern()

	// Create a image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, pngWidth, pngHeight))

	for y := 0; y < pngHeight; y++ {
		// Bits 0–7 are reserved for other uses and should be set to zero.
		for x := 0; x < offset; x++ {
			img.Set(x, y, color.White)
		}

		// Bits 8-55 represent all six characters to be displayed
		for x := offset; x <= 55; x++ {
			if string(encoding[x-offset]) == "1" {
				img.Set(x, y, color.Black)
			} else {
				img.Set(x, y, color.White)
			}
		}

		// Bits 56–255 are reserved for other uses and should be set to zero.
		for x := 56; x < pngWidth; x++ {
			img.Set(x, y, color.White)
		}
	}

	return img, nil
}

func (a *asset) persistToFile(img *image.NRGBA) error {
	err := os.MkdirAll(outputDir(), os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/%s.png", outputDir(), a.idStr()))
	if err != nil {
		return err
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func outputDir() string {
	if flag.Lookup("test.v") == nil {
		return "outputs"
	}

	return "test_outputs"
}

func encodingForDigit(digit int) (string, error) {
	switch digit {
	case 0:
		return "01110111", nil
	case 1:
		return "01000010", nil
	case 2:
		return "10110101", nil
	case 3:
		return "11010110", nil
	case 4:
		return "11000011", nil
	case 5:
		return "11010101", nil
	case 6:
		return "11110101", nil
	case 7:
		return "01000110", nil
	case 8:
		return "11110111", nil
	case 9:
		return "11010111", nil
	default:
		return "", fmt.Errorf("No encoding for given digit: %d", digit)
	}
}

func sliceAtoi(sliceStr []string) ([]int, error) {
	var sliceInt = make([]int, 0)

	for _, str := range sliceStr {
		i, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		sliceInt = append(sliceInt, i)
	}

	return sliceInt, nil
}
