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

// ParseOBOfromfile process entire OBO file from file
func ParseOBOfromfile(path string) (*OBOdocument, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return ParseOBO(f)
}

// ParseOBO process entire OBO file from io.Reader
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
				continue
			}
			document.Header = append(document.Header, tag)
			if tag.Tag == "format-version" {
				document.OBOversion = tag.Value
			}
		}
	}
	if err := oboScanner.Err(); err != nil {
		return nil, fmt.Errorf("Could not scan: %v", err)
	}
	if document.OBOversion == "" {
		log.Print("\"format-version\" tag did not found in header!")
	}

	// parse stanzas
	// for oboScanner.Scan() {
	// 	line := oboScanner.Text()
	// }
	// if err := oboScanner.Err(); err != nil {
	// 	return nil, fmt.Errorf("Could not scan: %v", err)
	// }

	return document, nil
}

func parseStanza(lines []string) (*Stanza, error) {

	if len(lines) > 2 {
		return nil, errors.New("A Stanza must contains at least 2 lines")
	}

	s := NewStanza()

	tagname := strings.Trim(lines[0], " ")
	if !strings.HasPrefix(tagname, "[") || !strings.HasSuffix(tagname, "]") {
		return nil, errors.New("Stanza name not correct")
	}
	s.Type = tagname[1 : len(tagname)-1]

	for _, line := range lines[1:] {
		tag, err := parseTagValuePair(line)
		if err != nil {
			log.Printf("Could not parse tag %q: %v", line, err)
			continue
		}
		s.Tags = append(s.Tags, tag)
	}

	return s, nil
}

func parseTagValuePair(line string) (*TagValuePair, error) {
	tvp := NewTagValuePair()

	// extract tag
	tag := strings.SplitN(line, ":", 2)
	if len(tag) != 2 {
		return nil, errors.New("Could not extract tag")
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
	// TUDO: implement

	// set value
	tvp.Value = val
	return tvp, nil
}
