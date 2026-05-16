package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func findAllOccurrences(data, pattern []byte) []int {
	patternMatchOffsets := []int{}
	startIndex := 0
	for {
		// search for matching bytes from startindex and forward
		patternMatchOffset := bytes.Index(data[startIndex:], pattern)
		if patternMatchOffset == -1 {
			break // if no more matches are found then break
		}
		nextStartIndex := patternMatchOffset + startIndex
		patternMatchOffsets = append(patternMatchOffsets, nextStartIndex) // append to the list of matching offsets
		startIndex = nextStartIndex + 1
	}
	return patternMatchOffsets
}

func exportFileType(data, startBytes, endBytes []byte, fileType string) {
	headerMatches := findAllOccurrences(data, startBytes)
	trailerMatches := findAllOccurrences(data, endBytes)
	if len(headerMatches) <= 0 || len(trailerMatches) <= 0 {
		return
	}
	fmt.Println("[?] Parsing passed bytes (", len(data), ")")
	fmt.Printf("[?] Found %d %s headers and %d %s trailers\n", len(headerMatches), fileType, len(trailerMatches), fileType)

	filesExtracted := 0
	for i := 0; i < len(headerMatches) && i < len(trailerMatches); i++ {

		start := headerMatches[i]
		end := trailerMatches[i] + len(endBytes)

		if end > len(data) {
			fmt.Printf("[!] Trailer ending at offset %d exceeds data length %d skipping... %d\n", end, len(data), i)
			continue
		}
		if start > end {
			fmt.Printf("[!] Header at offset %d is after trailer at offset %d skipping...\n", start, end)
			continue
		}

		fmt.Printf("[*] Extracting file %d: starts at %d, ends at %d, size %d bytes\n", i, start, end, end-start)
		fileByteSequence := data[start:end]
		fileName := "export_png_" + strconv.Itoa(i) + "." + fileType

		err := os.WriteFile(fileName, fileByteSequence, 0644)
		if err != nil {
			fmt.Printf("[-] Failed to write %s: %v\n", fileName, err)
		} else {
			fmt.Printf("[+] Exported file: %s\n", fileName)
			filesExtracted++
		}
	}

	fmt.Printf("[?] Extraction completed. Total files extracted: %d\n", filesExtracted)

}

func exportAllImplementedFileTypes(data []byte) {
	exportFileType(data,
		[]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
		[]byte{0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82},
		"png")

	exportFileType(data,
		[]byte{0xFF, 0xD8, 0xFF},
		[]byte{0xFF, 0xD9},
		"jpg")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wpengineunpack <input_file>")
		return
	}
	inputFile := os.Args[1]
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("[-] Failed to read file:", err)
		return
	}
	exportAllImplementedFileTypes(data)
}

