package models

type Value struct {
	ReadOnly bool
	Level    Level
	Value    string
}

type OValue struct {
	ReadOnly bool
	Value    string
}
