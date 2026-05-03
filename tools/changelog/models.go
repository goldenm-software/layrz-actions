package main

type Commit struct {
	ShortHash   string
	FullHash    string
	Subject     string
	Body        string
	Author      string
	CoAuthors   []string
	Type        string
	Scope       string
	IsBreaking  bool
	Description string
}

type section struct {
	key    string
	header string
	lines  []string
}
