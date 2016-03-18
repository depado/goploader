package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/Depado/goploader/server/conf"

	"github.com/gin-gonic/gin"
)

// Exported colors to use outside the package
var (
	Red      = "\x1b[31m"
	Green    = "\x1b[32m"
	Cyan     = "\x1b[36m"
	Yellow   = "\x1b[33m"
	Reset    = string([]byte{27, 91, 48, 109})
	colorMap = map[string]string{
		"ERROR":      Red,
		"INFO":       Green,
		"DEBUG":      Yellow,
		"server":     Cyan,
		"monitoring": Yellow,
	}
)

func getColor(str string) string {
	if val, ok := colorMap[str]; ok {
		return val
	}
	return Reset
}

func generic(status, from, ip, message string, extra ...interface{}) {
	args := []interface{}{
		getColor(status), status, Reset,
		getColor(from), from, Reset,
		time.Now().Format("2006/01/02 - 15:04:05"),
		ip,
		message,
	}
	strfmt := "%s%5s %s|%s %10v %s| %v | %15v | %s\n"
	if len(extra) > 0 {
		var x string
		for i, xv := range extra {
			if i == 0 {
				x = fmt.Sprintf("%v", xv)
			} else {
				x = fmt.Sprintf("%s %v", x, xv)
			}
		}
		args = append(args, x)
		strfmt = "%s%5s %s|%s %10v %s| %v | %15v | %s : %s\n"
	}
	fmt.Fprintf(os.Stdout, strfmt, args...)
}

// Err logs a simple error without filling the IP field
func Err(from, message string, extra ...interface{}) {
	generic("ERROR", from, "", message, extra...)
}

// ErrC logs a simple error and fills the IP field with the gin.Context
func ErrC(c *gin.Context, from, message string, extra ...interface{}) {
	generic("ERROR", from, c.ClientIP(), message, extra...)
}

// Info logs a simple info message without filling the IP field
func Info(from, message string, extra ...interface{}) {
	generic("INFO", from, "", message, extra...)
}

// InfoC logs a simple info message and fills the IP field with the gin.Context
func InfoC(c *gin.Context, from, message string, extra ...interface{}) {
	generic("INFO", from, c.ClientIP(), message, extra...)
}

// Debug logs a simple debug message without filling the IP field
func Debug(from, message string, extra ...interface{}) {
	if conf.C.Debug {
		generic("DEBUG", from, "", message, extra...)
	}
}

// DebugC logs a simple debug message and fills the IP field with the gin.Context
func DebugC(c *gin.Context, from, message string, extra ...interface{}) {
	if conf.C.Debug {
		generic("DEBUG", from, c.ClientIP(), message, extra...)
	}
}
