package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

type logWrapper struct {
	*log.Logger
}

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func getWriter(level int, cutoff int) io.Writer {
	if level >= cutoff {
		return os.Stdout
	} else {
		return ioutil.Discard
	}
}

func Init(level int) {
	Trace = log.New(getWriter(level, 4),
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(getWriter(level, 3),
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(getWriter(level, 2),
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(getWriter(level, 1),
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
