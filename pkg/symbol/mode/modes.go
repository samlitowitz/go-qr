package mode

type Mode int

const (
	Alphanumeric Mode = iota
	Byte         Mode = iota
	Numeric      Mode = iota
)
