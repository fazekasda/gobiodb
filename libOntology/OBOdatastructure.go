package libOntology

// OBOdocument store an entire OBO file
type OBOdocument struct {
	OBOversion string
	Header     []*TagValuePair
	Stanzas    []*Stanza
}

// Stanza store an OBO entry
type Stanza struct {
	Type string // Term, Typedef or Instance
	ID   string
	Tags []*TagValuePair
}

// TagValuePair handle a line of an OBO file
type TagValuePair struct {
	Tag               string
	Value             string
	TrailingModifiers map[string]string
	Comment           string
}

// NewOBOdocument create empty OBOdocument struct
func NewOBOdocument() *OBOdocument {
	obod := new(OBOdocument)
	obod.Header = make([]*TagValuePair, 0)
	obod.Stanzas = make([]*Stanza, 0)
	return obod
}

// NewStanza create empty Stanza struct
func NewStanza() *Stanza {
	s := new(Stanza)
	s.Tags = make([]*TagValuePair, 0)
	return s
}

// NewTagValuePair create empty TagValuePair struct
func NewTagValuePair() *TagValuePair {
	tvp := new(TagValuePair)
	tvp.TrailingModifiers = make(map[string]string)
	return tvp
}
