package logger

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	logLevel = DEBUG

	logStrPool = sync.Pool{
		New: func() interface{} {
			return new(strings.Builder)
		},
	}
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
}

func SetLogLevel(level string) {
	switch level {
	case "debug":
		logLevel = DEBUG
	case "info":
		logLevel = INFO
	case "warn":
		logLevel = WARN
	case "error":
		logLevel = ERROR
	case "fatal":
		logLevel = FATAL
	default:
		log.Fatalf("invalid log level %s", level)
	}
}

func printf(level string, tag string, args ...interface{}) {
	log.Println(formatf(level, tag, args...))
}

func fatalf(tag string, args ...interface{}) {
	log.Fatalln(formatf("F", tag, args...))
}

func formatf(level string, tag string, args ...interface{}) string {
	b := logStrPool.Get().(*strings.Builder)
	b.Reset()
	defer logStrPool.Put(b)

	b.WriteString(level)
	b.WriteRune('\t')
	b.WriteString(tag)
	b.WriteRune('\t')
	for idx, v := range args {
		switch v := v.(type) {
		default:
			b.WriteString(fmt.Sprintf("%v", v))
		case error:
			b.WriteString(fmt.Sprintf("%s", v.Error()))
		}
		if idx < len(args)-1 {
			b.WriteRune(' ')
		}
	}

	return b.String()
}

func Debug(tag string, args ...interface{}) {
	if logLevel <= DEBUG {
		printf("D", tag, args...)
	}
}

func Debugf(tag string, template string, args ...interface{}) {
	if logLevel <= DEBUG {
		printf("D", tag, fmt.Sprintf(template, args...))
	}
}

func Info(tag string, args ...interface{}) {
	if logLevel <= INFO {
		printf("I", tag, args...)
	}
}

func Infof(tag string, template string, args ...interface{}) {
	if logLevel <= INFO {
		printf("I", tag, fmt.Sprintf(template, args...))
	}
}

func Warn(tag string, args ...interface{}) {
	if logLevel <= WARN {
		printf("W", tag, args...)
	}
}

func Warnf(tag string, template string, args ...interface{}) {
	if logLevel <= WARN {
		printf("W", tag, fmt.Sprintf(template, args...))
	}
}

func Error(tag string, args ...interface{}) {
	if logLevel <= ERROR {
		printf("E", tag, args...)
	}
}

func Errorf(tag string, template string, args ...interface{}) {
	if logLevel <= ERROR {
		printf("E", tag, fmt.Sprintf(template, args...))
	}
}

func Fatal(tag string, args ...interface{}) {
	if logLevel <= FATAL {
		fatalf(tag, args...)
	}
}

func Fatalf(tag string, template string, args ...interface{}) {
	if logLevel <= FATAL {
		fatalf(tag, fmt.Sprintf(template, args...))
	}
}
