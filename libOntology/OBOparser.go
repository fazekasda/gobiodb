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
	var nextline string
	for oboScanner.Scan() {
		line := oboScanner.Text()
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") {
			nextline = line
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
	stanzaLinesBuf := make([]string, 0)
	stanzaLinesBuf = append(stanzaLinesBuf, nextline)
	stanzaLinesBufChan := make(chan []string)
	stanzaParsedChan := make(chan *Stanza)
	doneCollectingStanzas := make(chan bool)
	// parser goroutine
	go func() {
		for lines := range stanzaLinesBufChan {
			stanza, err := parseStanza(lines)
			if err != nil {
				log.Printf("Could not parse stanza %q: %v", stanzaLinesBuf, err)
				continue
			}
			stanzaParsedChan <- stanza
		}
		close(stanzaParsedChan)
	}()
	// collect parsed stanzaz
	go func() {
		for s := range stanzaParsedChan {
			document.Stanzas = append(document.Stanzas, s)
		}
		doneCollectingStanzas <- true
	}()
	// scan lines
	for oboScanner.Scan() {
		line := oboScanner.Text()
		if line == "" {
			continue
		} else if strings.HasPrefix(line, "[") {
			if len(stanzaLinesBuf) > 0 {
				stanzaLinesBufChan <- stanzaLinesBuf
				stanzaLinesBuf = make([]string, 0)
			}
			stanzaLinesBuf = append(stanzaLinesBuf, line)
		} else {
			stanzaLinesBuf = append(stanzaLinesBuf, line)
		}
	}
	close(stanzaLinesBufChan)
	<-doneCollectingStanzas
	if err := oboScanner.Err(); err != nil {
		return nil, fmt.Errorf("Could not scan: %v", err)
	}

	return document, nil
}

func parseStanza(lines []string) (*Stanza, error) {

	if len(lines) < 2 {
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
		if tag.Tag == "id" {
			s.ID = tag.Value
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

func extractTrailingModifiers(line string) (string, map[string]string, error) {

	//line, brace := extractAmongRunes(line, '{', '}')

	return line, make(map[string]string), nil
}

func extractAmongRunes(line string, openRune, closeRune rune) (string, string) {
	if strings.IndexRune(line, openRune) == -1 || strings.IndexRune(line, closeRune) == -1 {
		return line, ""
	}

	openPos := -1
	closePos := -1
	for i, c := range line {
		switch {
		case i == 0 && c == openRune:
			openPos = i
		case c == openRune && line[i-1] != '\\' && closePos == -1 && openPos == -1:
			openPos = i
		case c == openRune && line[i-1] != '\\' && closePos != -1 && openPos != -1:
			openPos = i
			closePos = -1
		case c == closeRune && line[i-1] != '\\' && closePos == -1 && openPos != -1:
			closePos = i
		}
	}
	if openPos == -1 || closePos == -1 {
		return line, ""
	}

	lineFront := ""
	lineBack := ""
	if openPos > 0 {
		lineFront = line[0:openPos]
	}
	if closePos+1 < len(line) {
		lineFront = line[closePos+1 : len(line)]
	}
	return lineFront + lineBack, line[openPos+1 : closePos]
}
