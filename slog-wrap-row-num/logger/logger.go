package logger

type SampleLogger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}

type NoOp struct{}

func NewNoOp() SampleLogger {
	return &NoOp{}
}

func (NoOp) Info(_ string, _ ...interface{}) {}

func (NoOp) Error(_ string, _ ...interface{}) {}

var (
	loggerImpl SampleLogger
)

func init() {
	loggerImpl = NewNoOp()
}

// SetLogger sets logger implementation.
// Note: not thread safe!
func SetLogger(l SampleLogger) {
	loggerImpl = l
}

func Info(format string, v ...interface{}) {
	loggerImpl.Info(format, v...)
}

func Error(format string, v ...interface{}) {
	loggerImpl.Error(format, v...)
}
