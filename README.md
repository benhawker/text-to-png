### ID to PNG

### Running this program

Run the program:
```
go run main.go
```

Run (with another input file):
```
go run main.go -inputFile=inputs/test_input.txt
```

Run the tests:
```
go test -v
```

The output from running the program against `inputs/test_input.txt` is provided in the committed `outputs/` folder.

A binary `./test-to-png` built with `go1.14.2 darwin/amd64` is also committed. No external dependencies outside the Go standard library are required.

### Potential Improvements

- e2e testing of file creation

- Write to tmp file(s) when running tests enabling the removal of `outputDir()` and it's call to `flag.Lookup("test.v")`.

- ...or if we refactor to allow us to pass in something that meets the https://golang.org/pkg/io/#Writer interface we can more easily test output without writing to a file.

- Consider storing data on the asset as `[]byte`. Program can use https://golang.org/pkg/encoding/binary/ to convert to `int` where required (checksum calculation) and string conversion with `string([]byte{65, 66, 67})`

- This would also allow us to improve the fn `setChecksum()` and the check against it's length which limits the flexibility.

 
### Task

A customer has requested that we use our robot platform to survey their premises once a day and report the location of various key assets. The customer has an existing installation of remotely-updateable, six-character, seven-segment displays, and one of these is attached to each asset.

During normal working hours, these displays show customer-speciﬁc information. We have agreed with the customer that we will will use these displays outside of working hours to determine the location of their assets. The customer will remotely update each display so that it shows a unique ID code that identiﬁes the asset, then our robot will autonomously drive through the entire premises and use a camera-based OCR (Optical Character Recognition) system to detect each display, read the code, and work out the real-world position of the asset.

The content of each display can be controlled by uploading a specially-crafted PNG image ﬁle that contains 256 bits of information encoded as black and white pixels. 48 pixels in this image directly control the state of the segments on the displays, and the rest are reserved for other purposes.

We will be provided with a list of asset IDs, each of which is a numeric value between 0 and 9999, supporting up to 10,000 unique displays. For each of these asset IDs, we will supply the customer with a PNG ﬁle that will be uploaded to the display and will program the display to show the asset ID.

We only require four characters to represent each unique ID, and we have decided to use two extra characters to implement error detection – we will include a checksum that we can verify when the display is read, allowing us to discard most codes which have an error in them.

The task to be completed in this exercise is to consume the list of asset IDs and output the corresponding PNG ﬁles that will be supplied to the customer.

### Functional specification

Take the following steps to generate the required output:

**1.** Given a text ﬁle containing a list of unique, four-digit numerical asset IDs in a text ﬁle format, calculate a two-digit checksum for each ID and add it as a preﬁx to the ID to generate a six-digit code.

Each individual asset ID is represented as a 4-digit number between `0` and `9999`. A simple checksum can be implemented by using modular arithmetic with a base of 97 to calculate two check digits. The checksum c of a number a with 4 digits can be calculated as:

```
c = (a1 + (10 * a2) + (100 * a3) + (1000 * a4)) mod 97
```

For example, the checksum for the asset ID 1337 can be calculated as follows:
```
(1 + (10 * 3) + (100 * 3) + (1000 * 7)) mod 97 = 7331 mod 97 = 56
```

This can be preﬁxed to the original asset ID to give the six-digit checksummed code `561337`.

**2.**  For each asset ID and checksum, calculate the pattern of bits required to activate the required segements on the display, so that all six characters will be shown.

Each of the six characters on the display consists of seven segments which can be turned on or oﬀ. The ten decimal digits can be displayed as 

2 Each character is represented by a diﬀerent 8-bit pattern in which each bit controls one segment. The bits are arranged as seen in `lcd-bits.png`

NB: bit 4 is always set to zero. In this system, the character `5` would be encoded as the bit string `11010101`.

I believe this should result in the following encodings:

```
x 0123 4 567
------------
0 0111 0 111
1 0100 0 010
2 1011 0 101
3 1101 0 110
4 1100 0 011
5 1101 0 101
6 1111 0 101
7 0100 0 110
8 1111 0 111
9 1101 0 111

```

**3.**  For each pattern of bits, generate a PNG ﬁle that encodes this pattern, and output it to disk with a ﬁle name of the original asset ID.

After constructing a pattern of bits for each character, they must be assembled into a PNG ﬁle which encodes this information. Your code must generate a valid PNG ﬁle which is 256 pixels wide, 1 pixel in height, and uses a 1-bit colour depth.

Each pixel in the PNG ﬁle will represent the value of a single bit in a 256-bit pattern – if the pixel is white, the bit is 0, with a 1 bit represented by a black pixel.

Only 48 bits (6 * 8) are required to represent all six characters to be displayed, using the bit pattern you constructed in the last step. You should encode this bit pattern into the PNG ﬁle by setting the appropriate pixel values, starting with an oﬀset of 8 bits. Bits 0–7 and bits 56–255 are reserved for other uses and should be set to zero.

For example, encoding the checksummed asset ID `561337` as described above would result in the bit pattern:

```
110101011111010101000010110101101101011001000110
```