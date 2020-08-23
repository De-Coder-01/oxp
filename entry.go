package oxp

type Sense struct {
	Shcut    string
	Def      string
	Examples []string
}

type Entry struct {
	Headword string
	POS      string
	Senses   []Sense
}
