package mode

import "fmt"

type EncodingError struct {
	Pos int
	Err error
}

func (err *EncodingError) Error() string {
	return fmt.Sprintf(
		"Error at position %d: %s",
		err.Pos,
		err.Err,
	)
}

func (err *EncodingError) Unwrap() error {
	return err.Err
}

type OutOfBoundsError struct {
	Given  string
	Bounds string
}

func (err *OutOfBoundsError) Error() string {
	return fmt.Sprintf(
		"Out of Bounds: `%s` Given, expected in %s",
		err.Given,
		err.Bounds,
	)
}
