package gexec

import (
	"bytes"
	"io"
	"os/exec"
)

type Options struct {
	Dir   string
	Stdin io.Reader
}

func Exec(name string, args ...string) (*Result, error) {
	return ExecWithOpt(nil, name, args...)
}

func ExecWithOpt(o *Options, name string, args ...string) (*Result, error) {
	if o == nil {
		o = &Options{}
	}

	cmd := exec.Command(name, args...)

	stdout := bytes.NewBuffer(make([]byte, 0, 1024))
	stderr := bytes.NewBuffer(make([]byte, 0, 1024))

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if o.Dir != "" {
		cmd.Dir = o.Dir
	}

	if o.Stdin != nil {
		cmd.Stdin = o.Stdin
	}

	if err := cmd.Run(); err != nil {
		return nil, &Error{err: err, Stderr: stderr.String()}
	}

	result := &Result{
		stdout: stdout.Bytes(),
	}

	return result, nil
}
