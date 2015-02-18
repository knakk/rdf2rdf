package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"unicode/utf8"

	"github.com/knakk/rdf"
)

var usage = `rdf2rdf
-------
Convert between different RDF serialization formats.

By default the converter is streaming both input and output, emitting
converted triples/quads as soon as they are available. This ensures you can
convert huge files with minimum memory footprint. However, if you have
small datasets you can choose to load all data into memory before conversion.
This makes it possible to sort the data, and generate more compact Turtle
serializations, maximizing predicate and object lists. Do this by setting the
flag stream=false.

Conversion from a quad-format to a triple-format will disregard the triple's
context (graph). Conversion from a triple-format to a quad-format is not
supported.

Input and ouput formats are determined by file extensions, according to
the following table:

  Format    | File extension
  ----------|-------------------
  N-Triples | .nt
  N-Quads   | .nq
  RDF/XML   | .rdf .rdfxml .xml
  Turtle    | .ttl

Usage:
	rdf2rdf -in=input.xml -out=output.ttl

Options:
	-h --help      Show this message.
	-in            Input file. 
	-out           Output file.
	-stream=true   Streaming mode.

`

func main() {
	log.SetFlags(0)
	log.SetPrefix("ERROR: ")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
	}
	input := flag.String("in", "", "Input file")
	output := flag.String("out", "", "Output file")
	flag.Parse()

	if *input == "" || *output == "" {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	inFile, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}

	outFile, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}

	inExt := fileExtension(*input)
	outExt := fileExtension(*output)

	if inExt == outExt {
		log.Fatal("No conversion necessary. Input and output formats are identical.")
	}

	var inFormat, outFormat rdf.Format

	switch inExt {
	case "nt":
		inFormat = rdf.FormatNT
	case "nq":
		inFormat = rdf.FormatNQ
	case "ttl":
		inFormat = rdf.FormatTTL
	case "":
		log.Fatal("Unknown file format. No file extension on input file.")
	default:
		log.Fatalf("Unsopported file exension on input file: %s", *inFile)
	}

	switch outExt {
	case "nt":
		outFormat = rdf.FormatNT
	case "nq":
		// No other quad-formats supported ATM
		log.Fatal("Serializing to N-Quads currently not supported.")
	case "ttl":
		outFormat = rdf.FormatTTL
	case "":
		log.Fatal("Unknown file format. No file extension on output file.")
	default:
		log.Fatalf("Unsopported file exension on output file: %s", *outFile)
	}

	tripleToTriple(inFile, outFile, inFormat, outFormat)
}

func tripleToTriple(inFile, outFile *os.File, inFormat, outFormat rdf.Format) {
	dec := rdf.NewTripleDecoder(inFile, inFormat)
	// TODO set base to file name?
	enc := rdf.NewTripleEncoder(outFile, outFormat)

	for t, err := dec.Decode(); err != io.EOF; t, err = dec.Decode() {
		if err != nil {
			log.Fatal(err)
		}
		err = enc.Encode(t)
		if err != nil {
			log.Fatal(err)
		}
	}
	err := enc.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func fileExtension(s string) string {
	i := len(s)
	for i > 0 {
		r, w := utf8.DecodeLastRuneInString(s[0:i])
		if r == '.' {
			return s[i:len(s)]
			break
		}
		i -= w
	}
	return "not found"
}
