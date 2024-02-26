package main

type Q struct {
	Type        int
	Id          int
	Qq          *Qq
	PickedCount int
}

type Qq struct {
	Marks float32
	Text  []byte // can be empty
	Qq    []*Qq  // sub-divisions
}
