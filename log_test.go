package log

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	Printf("this is printf, %d", 1)
	Printf("this is printf, %d", 2)
	Infof("this is info, %d", 3)
	Debugf("this is debug, %d", 4)
	Errorf("this is error, %d", 5)
	Warnf("this is warn, %d", 6)
}

func TestLogger(t *testing.T) {
	log := NewLogger(OpenPrint(), OpenWriteFile(), Filepath("./log/"), FileName(time.Now().Format("060102150405")+"_"+localIP()+".log"))
	log.Printf("this is printf, %d", 1)
	log.Infof("this is info, %d", 2)
	log.Debugf("this is debug, %d", 3)
	log.Errorf("this is error, %d", 4)
	log.Warnf("this is warn, %d", 5)
}
