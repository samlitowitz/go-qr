//go:generate go run golang.org/x/tools/cmd/stringer -type=Level
package symbol

type Level int

const (
	L Level = iota
	M Level = iota
	Q Level = iota
	H Level = iota
)
