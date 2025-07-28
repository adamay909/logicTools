package gentzen

import (
	"io"
	"log"
	"strings"
)

var checkLog strings.Builder
var dLog strings.Builder

var logger *log.Logger
var debug *log.Logger

// SetDebug turns on verbose logging for debugging
func SetDebug(v bool) {
	oDEBUG = v
}

// SetDebuglog sets debugLog to w
func SetDebuglog(w io.Writer) {
	SetDebug(true)
	debug = log.New(w, "", 0)
}

// Debug adds debug message
func Debug(a ...any) {
	if !oDEBUG {
		return
	}
	debug.Print(a...)
}

// ShowDebugLog displays debug log.
func ShowDebugLog() string {

	return dLog.String()

}

// ShowLog displays log. Currently, logging is only done by proof checker.
func ShowLog() string {

	return checkLog.String()

}

// ClearLog clears log.
func ClearLog() {
	checkLog.Reset()
	return
}

// WriteLog writes s to log.
func WriteLog(s string, p string) {
	logger.SetPrefix(p)
	logger.Print(s)
	return
}
