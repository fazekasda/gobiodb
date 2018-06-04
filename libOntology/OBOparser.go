package libOntology

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func ParseOBOfromfile(path string) (*OBOdocument, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return ParseOBO(f)
}

func ParseOBO(r io.Reader) (*OBOdocument, error) {

	oboScanner := bufio.NewScanner(r)
	document := NewOBOdocument()

	// parse header
	for oboScanner.Scan() {
		line := oboScanner.Text()
		if line == "" {
			break
		} else {
			tag, err := parseTagValuePair(line)
			if err != nil {
				log.Printf("Could not parse tag %q: %v", line, err)
			} else {
				document.Header = append(document.Header, tag)
			}
		}
	}
	if err := oboScanner.Err(); err != nil {
		return nil, fmt.Errorf("Could not scan: %v", err)
	}

	// parse stanzas
	// for oboScanner.Scan() {
	// 	line := oboScanner.Text()
	// }
	// if err := oboScanner.Err(); err != nil {
	// 	return nil, fmt.Errorf("Could not scan: %v", err)
	// }

	return nil, nil
}

func parseTagValuePair(line string) (*TagValuePair, error) {
	tvp := NewTagValuePair()

	// extract tag
	tag := strings.SplitN(line, ":", 2)
	if len(tag) != 2 {
		return nil, errors.New("Could not ertaxt tag")
	}
	tvp.Tag = strings.Trim(tag[0], " ")

	// cut comment
	val := strings.Trim(tag[1], " ")
	if strings.Contains(val, "!") {
		comm := strings.SplitN(val, "!", 2)
		tvp.Comment = strings.Trim(comm[1], " ")
		val = strings.Trim(comm[0], " ")
	}

	// extract TrailingModifiers
	// TODU: implement

	// set value
	tvp.Value = val
	return tvp, nil
}
