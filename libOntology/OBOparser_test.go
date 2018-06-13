package libOntology

import (
	"testing"
)

func TestParseTagValuePair(t *testing.T) {
	tvp, err := parseTagValuePair("incorrect tag value pais")
	if tvp != nil || err.Error() != "Could not extract tag" {
		t.Error("tag value pair syntax checking not worked")
	}
	tvpText := "key: value ! comment"
	tvp, err = parseTagValuePair(tvpText)
	if err != nil {
		t.Errorf("error during parse %q: %v", tvpText, err)
	}
	if tvp.Tag != "key" {
		t.Errorf("could not parse Tag from %q, found: %q", tvpText, tvp.Tag)
	}
	if tvp.Value != "value" {
		t.Errorf("could not parse Value from %q, found: %q", tvpText, tvp.Value)
	}
	if tvp.Comment != "comment" {
		t.Errorf("could not parse Comment from %q, found: %q", tvpText, tvp.Comment)
	}
}

func TestParseStanza(t *testing.T) {
	//correct stanza
	sanzaTex1 := []string{
		"[Term]",
		"key1: value1 ! comment1",
		"key2: value2 ! comment2",
		"incorrect tag value pais",
	}
	stanza, err := parseStanza(sanzaTex1)
	if err != nil {
		t.Errorf("error during parse %q: %v", sanzaTex1, err)
	}
	if stanza.Type != "Term" {
		t.Errorf("could not parse stanza Type from %q, found: %q", sanzaTex1, stanza.Type)
	}
	if len(stanza.Tags) != 2 {
		t.Error("could not parse all tag value pair in stanza")
	}

	//stanza name error
	sanzaTex2 := []string{
		"Term]",
		"key1: value1 ! comment1",
		"key2: value2 ! comment2",
	}
	stanza, err = parseStanza(sanzaTex2)
	if stanza != nil || err.Error() != "Stanza name not correct" {
		t.Errorf("error checking stanza name %q: %v", sanzaTex2, err)
	}

	//empty stanza
	sanzaTex3 := []string{
		"[Term]",
	}
	stanza, err = parseStanza(sanzaTex3)
	if stanza != nil || err.Error() != "A Stanza must contains at least 2 lines" {
		t.Errorf("error checking stanza length %q: %v", sanzaTex3, err)
	}
}

func TestExtractAmongRunes(t *testing.T) {
	rem, ext := extractAmongRunes("there aren't any to extract", '{', '}')
	if rem != "there aren't any to extract" || ext != "" {
		t.Errorf("Error with find open and close runes")
	}
	rem, ext = extractAmongRunes("tex t\\{\\}o {ex{tract} {secound}", '{', '}')
	if rem != "tex t\\{\\}o {ex{tract} " || ext != "secound" {
		t.Errorf("Error with extract text: %q - %q", rem, ext)
	}
}
