package main

import (
	"flag"
	"fmt"
	"os"
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
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
	}
	flag.Usage()

}
