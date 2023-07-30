package utils

type Logger interface {
	Logf(format string, args ...interface{})
}

type SMB interface {
	Get(filename string) ([]byte, error)
}
