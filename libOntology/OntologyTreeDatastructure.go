package libOntology

type dbref struct {
	DB string
	ID string
}

type TermStanza struct {
	ID           string
	name         string
	is_anonymous string
	alt_id       []string
	def          string
	comment      string
}

type TypedefStanza struct {
}

type InstanceStanza struct {
}
