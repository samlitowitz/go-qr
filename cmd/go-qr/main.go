package main

import (
	"flag"
	"log"
)

func main() {
	var encMode, ecLevel string

	flag.StringVar(&encMode, "mode", "numeric", "Mode for data encoding, one of [numeric, alphanumeric, byte]")
	flag.StringVar(&ecLevel, "eclevel", "l", "Error correction level, one of [l, m, q, h]")

	flag.Parse()

	// Choose EC Level
	switch ecLevel {
	case "l":
	case "m":
	case "q":
	case "h":

	default:
		log.Fatalf("Invalid error correction level `%s` provided, must be one of [l, m, q, h]", ecLevel)
	}

	// Determine Version

	switch encMode {
	case "numeric":
	case "alphanumeric":
	case "byte":

	default:
		log.Fatalf("Invalid mode `%s` provided, must be one of [numeric, alphanumeric, byte]", encMode)
	}

	// 1. Data Analysis
	// 1. Data Encoding
	// 1. Error Correction Coding
	// 1. Structure Final Message
	// 1. Module placement in Matrix
	// 1. Data Masking
	// 1. Format and Version Information
}
