package main

// Q for q store
type Q struct {
	Type        int
	Id          int
	Qq          *Question
	PickedCount int
}

type Question struct {
	Unit         string
	Section      string
	Marks        float32
	Text         string     // can be empty
	SubQuestions []Question // sub-divisions
}
