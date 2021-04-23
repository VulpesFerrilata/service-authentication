package mapper

type modelState int

const (
	New modelState = iota
	Modified
	Unchanged
)
