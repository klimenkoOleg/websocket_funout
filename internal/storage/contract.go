//go:generate mockgen -source $GOFILE -destination mock_test.go -package ${GOPACKAGE}
package storage

type Logger interface {
	Debug(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
}
