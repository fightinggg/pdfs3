package pdfs_err

import "os"

type PdfsError struct {
	msg string
}

func (p PdfsError) Error() string {
	return p.msg
}

func BuildPdfsError(msg string) error {
	return PdfsError{
		msg: msg,
	}
}

func PdfsNotFoundError() error {
	return os.ErrNotExist
}

func PdfsOutOfDiskError() error {
	return PdfsError{
		msg: "space enouth",
	}
}
