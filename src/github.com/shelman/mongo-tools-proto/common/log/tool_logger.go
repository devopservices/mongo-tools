package log

import (
	"fmt"
	"github.com/shelman/mongo-tools-proto/common/options"
	"github.com/shelman/mongo-tools-proto/common/util"
	"io"
	"os"
	"sync"
	"time"
)

//====== Tool Logger Definition ======

type ToolLogger struct {
	mutex     *sync.Mutex
	writer    io.Writer
	format    string
	verbosity int
}

func (tl *ToolLogger) SetVerbosity(verbosity *options.Verbosity) {
	if verbosity.Quiet {
		tl.verbosity = -1
	} else {
		tl.verbosity = len(verbosity.Verbose)
	}
}

func (tl *ToolLogger) SetWriter(writer io.Writer) {
	tl.writer = writer
}

func (tl *ToolLogger) SetDateFormat(dateFormat string) {
	tl.format = dateFormat
}

func (tl *ToolLogger) Logf(minVerb int, format string, a ...interface{}) {
	if minVerb < 0 {
		panic("cannot set a minimum log verbosity that is less than 0")
	}

	if minVerb <= tl.verbosity {
		tl.mutex.Lock()
		defer tl.mutex.Unlock()
		tl.log(fmt.Sprintf(format, a...))
	}
}

func (tl *ToolLogger) Log(minVerb int, msg string) {
	if minVerb < 0 {
		panic("cannot set a minimum log verbosity that is less than 0")
	}

	if minVerb <= tl.verbosity {
		tl.mutex.Lock()
		defer tl.mutex.Unlock()
		tl.log(msg)
	}
}

func (tl *ToolLogger) log(msg string) {
	fmt.Fprintf(tl.writer, "%v\t%v\n", time.Now().Format(tl.format), msg)
}

func NewToolLogger(verbosity *options.Verbosity) *ToolLogger {
	tl := &ToolLogger{
		mutex:  &sync.Mutex{},
		writer: os.Stderr, // default to stderr
		format: util.ToolTimeFormat,
	}
	tl.SetVerbosity(verbosity)
	return tl
}

//====== Global Logging ======

var globalToolLogger *ToolLogger

func assertGlobalToolLoggerInitialized() {
	if globalToolLogger == nil {
		panic("must initialize the global ToolLogger before use")
	}
}

func InitToolLogger(verbosity *options.Verbosity) {
	if globalToolLogger != nil {
		panic("global ToolLogger already initialized")
	}
	globalToolLogger = NewToolLogger(verbosity)
}

func Logf(minVerb int, format string, a ...interface{}) {
	assertGlobalToolLoggerInitialized()
	globalToolLogger.Logf(minVerb, format, a...)
}

func Log(minVerb int, msg string) {
	assertGlobalToolLoggerInitialized()
	globalToolLogger.Log(minVerb, msg)
}

func SetVerbosity(verbosity *options.Verbosity) {
	assertGlobalToolLoggerInitialized()
	globalToolLogger.SetVerbosity(verbosity)
}

func SetWriter(writer io.Writer) {
	assertGlobalToolLoggerInitialized()
	globalToolLogger.SetWriter(writer)
}

func SetDateFormat(dateFormat string) {
	assertGlobalToolLoggerInitialized()
	globalToolLogger.SetDateFormat(dateFormat)
}