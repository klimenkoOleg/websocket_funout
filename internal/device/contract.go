//go:generate mockgen -source $GOFILE -destination mock_test.go -package ${GOPACKAGE}
package device

type Logger interface {
	Debug(args ...interface{})
	Warn(args ...interface{})
}

type Connector interface {
	WriteJSON(v interface{}) error
	Close() error
}
