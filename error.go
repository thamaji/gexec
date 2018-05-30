package gexec

type Error struct {
	err    error
	Stderr string
}

func (err *Error) Error() string {
	return err.err.Error()
}
