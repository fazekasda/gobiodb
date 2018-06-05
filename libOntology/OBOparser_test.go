package libOntology

import (
	"testing"
)

func TestparseTagValuePair_extractError(t *testing.T) {
	tvp, err := parseTagValuePair("incorrect tag value pais")
	if tvp != nil || err.Error() != "Could not extract tag" {
		t.Error("tag value pair syntax checking not worked")
	}
}
