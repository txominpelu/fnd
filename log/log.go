package log

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

type StandardLogger struct {
	logger *logrus.Logger
}

// NewLogger initializes the standard logger
func NewLogger(logFile string) *StandardLogger {
	if logFile != "" {
		file, err := os.Create(logFile)
		if err != nil {
			panic(fmt.Sprintf("Error: %s. when creating file: %s", err, logFile))
		}
		return createLogger(file)
	} else {
		return createLogger(os.Stderr)
	}
}

func createLogger(file *os.File) *StandardLogger {
	baseLogger := logrus.New()
	baseLogger.Out = file
	baseLogger.Formatter = &logrus.TextFormatter{}
	var standardLogger = &StandardLogger{
		logger: baseLogger,
	}
	return standardLogger
}

func (s *StandardLogger) CheckError(err error, msg string) {
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		c := s.logger.
			WithError(err).
			WithField("file", file).
			WithField("line", line).
			WithField("func", details.Name())
		c.Errorln(msg)
		panic(fmt.Sprintf("Error: %s %s\n", err, msg))
	}
}

func (s *StandardLogger) WarnIfErr(err error, msg string) {
	if err != nil {
		s.logger.Warn(fmt.Sprintf("Warn: %s %s\n", err, msg))
	}
}
