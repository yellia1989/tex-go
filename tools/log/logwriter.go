package log

import (
	"os"
	"path/filepath"
	"strconv"
	"time"
)

//DAY for rotate log by day
const (
	DAY DateType = iota
	HOUR
)

//LogWriter is interface for different writer.
type LogWriter interface {
	Write(v []byte) (n int, err error)
	NeedPrefix() bool
}

//ConsoleWriter writes the logs to the console.
type ConsoleWriter struct {
}

//RollFileWriter struct for rotate logs by file size.
type RollFileWriter struct {
	logpath  string
	name     string
	num      int
	size     int64
	currSize int64
	currFile *os.File
	openTime int64
}

//DateWriter rotate logs by date.
type DateWriter struct {
	logpath   string
	name      string
	dateType  DateType
	num       int
	currDate  string
	currFile  *os.File
	openTime  int64
	hasPrefix bool
}

//HourWriter for rotate logs by hour
type HourWriter struct {
}

//DateType is uint8
type DateType uint8

func reOpenFile(path string, currFile **os.File, openTime *int64) error {
    cfg := cfg.Load().(config)
	*openTime = cfg.currUnixTime
	if *currFile != nil {
		(*currFile).Close()
	}
	of, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
        return err
    }

    *currFile = of
    return nil
}

func (w *ConsoleWriter) Write(v []byte) (n int, err error){
	return os.Stdout.Write(v)
}

//NeedPrefix shows whether needs the prefix for the console writer.
func (w *ConsoleWriter) NeedPrefix() bool {
	return true
}

//Write for writing []byte to the writter.
func (w *RollFileWriter) Write(v []byte) (n int, err error) {
    cfg := cfg.Load().(config)
	if w.currFile == nil || w.openTime+10 < cfg.currUnixTime {
		fullPath := filepath.Join(w.logpath, w.name+".log")
		err = reOpenFile(fullPath, &w.currFile, &w.openTime)
        if err != nil {
            return
        }
	}
	n, err = w.currFile.Write(v)
    if err != nil {
        return
    }
	w.currSize += int64(n)
	if w.currSize >= w.size {
		w.currSize = 0
		for i := w.num - 1; i >= 1; i-- {
			var n1, n2 string
			if i > 1 {
				n1 = strconv.Itoa(i - 1)
			}
			n2 = strconv.Itoa(i)
			p1 := filepath.Join(w.logpath, w.name+n1+".log")
			p2 := filepath.Join(w.logpath, w.name+n2+".log")
			if _, err := os.Stat(p1); !os.IsNotExist(err) {
				os.Rename(p1, p2)
			}
		}
		fullPath := filepath.Join(w.logpath, w.name+".log")
		err = reOpenFile(fullPath, &w.currFile, &w.openTime)
        if err != nil {
            return
        }
	}

    return
}

//NewRollFileWriter returns a RollFileWriter, rotate logs in sizeMB , and num files are keeped.
func NewRollFileWriter(logpath, name string, num, sizeMB int) *RollFileWriter {
	w := &RollFileWriter{
		logpath: logpath,
		name:    name,
		num:     num,
		size:    int64(sizeMB) * 1024 * 1024,
	}
	fullPath := filepath.Join(logpath, name+".log")
	st, _ := os.Stat(fullPath)
	if st != nil {
		w.currSize = st.Size()
	}
	return w
}

//NeedPrefix shows need prefix or not.
func (w *RollFileWriter) NeedPrefix() bool {
	return true
}

//Write method implement for the DateWriter
func (w *DateWriter) Write(v []byte) (n int, err error) {
    cfg := cfg.Load().(config)
	if w.currFile == nil || w.openTime+10 < cfg.currUnixTime {
		fullPath := filepath.Join(w.logpath, w.name+"_"+w.currDate+".log")
		err = reOpenFile(fullPath, &w.currFile, &w.openTime)
        if err != nil {
            return
        }
	}

	n, err = w.currFile.Write(v)
    if err != nil {
        return
    }
	currDate := w.getCurrDate()
	if w.currDate != currDate {
		w.currDate = currDate
		w.cleanOldLogs()
		fullPath := filepath.Join(w.logpath, w.name+"_"+w.currDate+".log")
		err = reOpenFile(fullPath, &w.currFile, &w.openTime)
        if err != nil {
            return
        }
	}

    return
}

//NeedPrefix shows whether needs prefix info for DateWriter or not.
func (w *DateWriter) NeedPrefix() bool {
	return w.hasPrefix
}

func (w *DateWriter) SetPrefix(enable bool) {
	w.hasPrefix = enable
}

//NewDateWriter returns a writer which keeps logs in hours or day format.
func NewDateWriter(logpath, name string, dateType DateType, num int) *DateWriter {
	w := &DateWriter{
		logpath:   logpath,
		name:      name,
		num:       num,
		dateType:  dateType,
		hasPrefix: true,
	}
	w.currDate = w.getCurrDate()
	return w
}

func (w *DateWriter) cleanOldLogs() {
	format := "20060102"
	duration := -time.Hour * 24
	if w.dateType == HOUR {
		format = "2006010215"
		duration = -time.Hour
	}

	t := time.Now()
	t = t.Add(duration * time.Duration(w.num))
	for i := 0; i < 30; i++ {
		t = t.Add(duration)
		k := t.Format(format)
		fullPath := filepath.Join(w.logpath, w.name+"_"+k+".log")
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			os.Remove(fullPath)
		}
	}
	return
}

func (w *DateWriter) getCurrDate() string {
    cfg := cfg.Load().(config)
	if w.dateType == HOUR {
		return cfg.currDateHour
	}
	return cfg.currDateDay // DAY
}
