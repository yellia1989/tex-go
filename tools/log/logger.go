package log

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"time"
    "sync"
    "sync/atomic"
    "github.com/yellia1989/tex-go/tools/util"
)

//DEBUG loglevel
const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type config struct {
	logLevel LogLevel
    framework_logLevel LogLevel

	currUnixTime int64
	currDateTime string
	currDateHour string
	currDateDay  string
}

var (
	logQueue  = make(chan *logValue, 10000)
	loggerMap = make(map[string]*Logger)
	writeDone = make(chan struct{})
    defaultLogger *Logger

    mu sync.Mutex
    cfg atomic.Value
)

//Logger is the struct with name and wirter.
type Logger struct {
	name   string
	writer LogWriter
}

//LogLevel is uint8 type
type LogLevel uint8

type logValue struct {
	level  LogLevel
	value  []byte
	fileNo string
	writer LogWriter
}

func init() {
	now := time.Now()
    defaultLogger = GetLogger("AB0D2927-C7EF-4E17-AB72-938D027B0D08")

    cfg.Store(config{
	    currUnixTime: now.Unix(),
	    currDateTime: now.Format("2006-01-02 15:04:05"),
	    currDateHour: now.Format("2006010215"),
	    currDateDay: now.Format("20060102"),
        logLevel: DEBUG,
        framework_logLevel: INFO,
    })

	go func() {
        defer func() {
		    if err := recover(); err != nil {
                // avoid timer panic
		    }
        }()
		tm := time.NewTimer(time.Second)
		for {
			now := time.Now()
			d := time.Second - time.Duration(now.Nanosecond())
			tm.Reset(d)
			<-tm.C
			now = time.Now()
            
            mu.Lock()
            old := cfg.Load().(config)
			old.currUnixTime = now.Unix()
			old.currDateTime = now.Format("2006-01-02 15:04:05")
			old.currDateHour = now.Format("2006010215")
			old.currDateDay = now.Format("20060102")
            cfg.Store(old)
            mu.Unlock()
		}
	}()
    
	go func() {
        for v := range logQueue {
            if v.writer == nil {
                // 保证所有日志都写成功
                writeDone <- struct{}{}
                return
            }
			v.writer.Write(v.value)
		}
    }()
}

//String return turns the LogLevel to string.
func (lv *LogLevel) String() string {
	switch *lv {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
    default:
        panic("UNKNOWN")
	}
}

// GetLogger return an logger instance
func GetLogger(name string) *Logger {
	if lg, ok := loggerMap[name]; ok {
		return lg
	}
	lg := &Logger{
		name:   name,
		writer: &ConsoleWriter{},
	}
	loggerMap[name] = lg
	return lg
}

//SetLevel sets the log level
func SetLevel(level LogLevel) {
    mu.Lock()
    old := cfg.Load().(config)
    old.logLevel = level
    cfg.Store(old)
    mu.Unlock()
}

func GetLevel() LogLevel {
    old := cfg.Load().(config)
    return old.logLevel
}

func SetFrameworkLevel(level LogLevel) {
    mu.Lock()
    old := cfg.Load().(config)
    old.framework_logLevel = level
    cfg.Store(old)
    mu.Unlock()
}

func GetFrameworkLevel() LogLevel {
    old := cfg.Load().(config)
    return old.framework_logLevel
}

//StringToLevel turns string to LogLevel
func StringToLevel(level string) LogLevel {
	switch level {
	case "DEBUG","debug":
		return DEBUG
	case "INFO","info":
		return INFO
	case "WARN","warn":
		return WARN
	case "ERROR","error":
		return ERROR
    default:
        panic("UNKNOWN")
	}
}

//SetLogName sets the log name
func (l *Logger) SetLogName(name string) {
	l.name = name
}

//SetFileRoller sets the file rolled by size in MB, with max num of files.
func (l *Logger) SetFileRoller(logpath string, num int, sizeMB int) error {
	if err := os.MkdirAll(logpath, 0755); err != nil {
		panic(err)
	}
	w := NewRollFileWriter(logpath, l.name, num, sizeMB)
	l.writer = w
	return nil
}

//IsConsoleWriter returns whether is consoleWriter or not.
func (l *Logger) IsConsoleWriter() bool {
	if reflect.TypeOf(l.writer) == reflect.TypeOf(&ConsoleWriter{}) {
		return true
	}
	return false
}

//SetWriter sets the writer to the logger.
func (l *Logger) SetWriter(w LogWriter) {
	l.writer = w
}

func (l *Logger) GetWriter() LogWriter {
    return l.writer
}

//SetDayRoller sets the logger to rotate by day, with max num files.
func (l *Logger) SetDayRoller(logpath string, num int) error {
	if err := os.MkdirAll(logpath, 0755); err != nil {
		return err
	}
	w := NewDateWriter(logpath, l.name, DAY, num)
	l.writer = w
	return nil
}

//SetHourRoller sets the logger to rotate by hour, with max num files.
func (l *Logger) SetHourRoller(logpath string, num int) error {
	if err := os.MkdirAll(logpath, 0755); err != nil {
		return err
	}
	w := NewDateWriter(logpath, l.name, HOUR, num)
	l.writer = w
	return nil
}

//SetConsole sets the logger with console writer.
func (l *Logger) SetConsole() {
	l.writer = &ConsoleWriter{}
}

//Debug logs interface in debug loglevel.
func (l *Logger) Debug(v ...interface{}) {
	l.writef(false, DEBUG, "", v)
}

//Info logs interface in Info loglevel.
func (l *Logger) Info(v ...interface{}) {
	l.writef(false, INFO, "", v)
}

//Warn logs interface in warning loglevel
func (l *Logger) Warn(v ...interface{}) {
	l.writef(false, WARN, "", v)
}

//Error logs interface in Error loglevel
func (l *Logger) Error(v ...interface{}) {
	l.writef(false, ERROR, "", v)
}

//Debugf logs interface in debug loglevel with formating string
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.writef(false, DEBUG, format, v)
}

//Infof logs interface in Infof loglevel with formating string
func (l *Logger) Infof(format string, v ...interface{}) {
	l.writef(false, INFO, format, v)
}

//Warnf logs interface in warning loglevel with formating string
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.writef(false, WARN, format, v)
}

//Errorf logs interface in Error loglevel with formating string
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.writef(false, ERROR, format, v)
}

func (l *Logger) writef(framework bool, level LogLevel, format string, v []interface{}) {
    cfg := cfg.Load().(config)
	if framework && level < cfg.framework_logLevel {
		return
	}
	if !framework && level < cfg.logLevel {
		return
	}

	buf := bytes.NewBuffer(nil)
	if l.writer.NeedPrefix() {
		fmt.Fprintf(buf, "%s|", cfg.currDateTime)
		if cfg.logLevel == DEBUG {
			_, file, line, ok := runtime.Caller(2)
			if !ok {
				file = "???"
				line = 0
			} else {
				file = filepath.Base(file)
			}
			fmt.Fprintf(buf, "%s:%d|g%d|", file, line, util.CurGoroutineID())
		}

        if framework {
            buf.WriteString("[framework]|")
        }

		buf.WriteString(level.String())
		buf.WriteByte('|')
	}

	if format == "" {
		fmt.Fprint(buf, v...)
	} else {
		fmt.Fprintf(buf, format, v...)
	}
	if l.writer.NeedPrefix() {
		buf.WriteByte('\n')
	}
	logQueue <- &logValue{value: buf.Bytes(), writer: l.writer}
}

//FlushLogger flushs all log to disk.
func FlushLogger() {
    log := &logValue{}
    logQueue <- log
    <-writeDone
}

//默认创建的logger
func GetDefaultLogger() *Logger {
    return defaultLogger
}
func Debugf(format string, v ...interface{}) {
	defaultLogger.writef(false, DEBUG, format, v)
}
func Infof(format string, v ...interface{}) {
	defaultLogger.writef(false, INFO, format, v)
}
func Warnf(format string, v ...interface{}) {
	defaultLogger.writef(false, WARN, format, v)
}
func Errorf(format string, v ...interface{}) {
	defaultLogger.writef(false, ERROR, format, v)
}
func Debug(v ...interface{}) {
	defaultLogger.writef(false, DEBUG, "", v)
}
func Info(v ...interface{}) {
	defaultLogger.writef(false, INFO, "", v)
}
func Warn(v ...interface{}) {
	defaultLogger.writef(false, WARN, "", v)
}
func Error(v ...interface{}) {
	defaultLogger.writef(false, ERROR, "", v)
}
func FDebugf(format string, v ...interface{}) {
	defaultLogger.writef(true, DEBUG, format, v)
}
func FInfof(format string, v ...interface{}) {
	defaultLogger.writef(true, INFO, format, v)
}
func FWarnf(format string, v ...interface{}) {
	defaultLogger.writef(true, WARN, format, v)
}
func FErrorf(format string, v ...interface{}) {
	defaultLogger.writef(true, ERROR, format, v)
}
func FDebug(v ...interface{}) {
	defaultLogger.writef(true, DEBUG, "", v)
}
func FInfo(v ...interface{}) {
	defaultLogger.writef(true, INFO, "", v)
}
func FWarn(v ...interface{}) {
	defaultLogger.writef(true, WARN, "", v)
}
func FError(v ...interface{}) {
	defaultLogger.writef(true, ERROR, "", v)
}
