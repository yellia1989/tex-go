package log

import (
	"testing"
	_ "time"
)

//TestLogger test logger writes.
func TestLogger(t *testing.T) {
	SetLevel(DEBUG)
	lg := GetLogger("debug")
	// lg.SetConsole()
	// lg.SetFileRoller("./logs", 3, 1)
	// lg.SetDayRoller("./logs", 2)
	// lg.SetHourRoller("./logs", 2)
	bs := make([]byte, 1024)
	longmsg := string(bs)
	for i := 0; i < 10; i++ {
		lg.Debugf("debugxxxxxxxxxxxxxxxxxxxxxxxxxxx:%d",i+1)
		lg.Infof(":%s:%d", longmsg, i+1)
		//lg.Warn("warn")
		//lg.Error("ERROR")
		//time.Sleep(time.Second)
	}
	//time.Sleep(time.Millisecond * 100)

	FlushLogger()
}

//TestGetLogList test get log list
func TestGetLogList(t *testing.T) {
	w := NewDateWriter("./logs", "abc", HOUR, 0)
	w.cleanOldLogs()
}

//BenchmarkRogger benchmark rogger writes.
func BenchmarkRogger(b *testing.B) {
    b.ReportAllocs() 
	SetLevel(DEBUG)
	for i := 0; i < b.N; i++ {
		Debug("hello")
	}
}
