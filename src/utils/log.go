package utils

import (
	"fmt"
	stdlog "log"
	"os"
	"time"
)

const (
	DEBUG = iota
	INFO
	NOTICE
	WARN
	ERROR
	CRIT
	PANIC
)

var Level = map[int]string{
	DEBUG:  "DEBUG",
	WARN:   "WARN",
	INFO:   "INFO",
	NOTICE: "NOTICE",
	ERROR:  "ERROR",
	CRIT:   "CRIT",
	PANIC:  "PANIC",
}

type Logger struct {
	log     *stdlog.Logger
	level   int
	logfile string
}

func NewLogger(logfile string, flag int, level int) (*Logger, error) {
	l := new(Logger)

	if logfile != "std" {
		out, err := os.Create(logfile)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		l.log = stdlog.New(out, "", flag)
	} else {
		l.log = stdlog.New(os.Stdout, "", flag)
	}

	l.level = level

	return l, nil
}

func (self *Logger) Printf(level int, format string, v ...interface{}) {
	self.log.Printf(Level[level]+" "+format, v...)
}

func (self *Logger) Debug(format string, v ...interface{}) {
	if DEBUG >= self.level {
		self.log.Printf(Level[DEBUG]+" "+format, v...)
	}
}

func (self *Logger) Info(format string, v ...interface{}) {
	if INFO >= self.level {
		self.log.Printf(Level[INFO]+" "+format, v...)
	}
}

func (self *Logger) Warn(format string, v ...interface{}) {
	if WARN >= self.level {
		self.log.Printf(Level[WARN]+" "+format, v...)
	}
}

func (self *Logger) Notice(format string, v ...interface{}) {
	if NOTICE >= self.level {
		self.log.Printf(Level[NOTICE]+" "+format, v...)
	}
}

func (self *Logger) Crit(format string, v ...interface{}) {
	if CRIT >= self.level {
		self.log.Printf(Level[CRIT]+" "+format, v...)
	}
}

func (self *Logger) Error(format string, v ...interface{}) {
	if ERROR >= self.level {
		self.log.Printf(Level[ERROR]+" "+format, v...)
	}
}

func (self *Logger) Panic(format string, v ...interface{}) {
	if PANIC >= self.level {
		self.log.Printf(Level[PANIC]+" "+format, v...)
	}
}

func (self *Logger) SetOutfile(logfile string) error {
	if logfile != "std" {
		out, err := os.Create(logfile)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		self.logfile = logfile
		self.log.SetOutput(out)
	} else {
		self.logfile = "std"
		self.log.SetOutput(os.Stdout)
	}

	return nil
}

func (self *Logger) CheckSize(size int64) (bool, error) {
	info, err := os.Stat(self.logfile)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}

	return info.Size() < size, nil
}

func MakeLogfile() string {
	return time.Now().Format("20060102-150405") + ".log"
}

/*
func main() {
	dir := "./log"
	file := MakeLogfile()

	log, err := NewLogger(dir+"/"+file, 2, WARN)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	log.Debug("%s%d%s", "test", 123, Level[DEBUG])
	log.Info("%s%d%s", "test", 123, Level[INFO])
	log.Warn("%s%d%s", "test", 123, Level[WARN])
	log.Error("%s%d%s", "test", 123, Level[ERROR])
}
*/
