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
