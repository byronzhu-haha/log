package log

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	defaultFlushSec = 30
	writeByteSize   = 2048
)

type Logger interface {
	Printf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Close()
}

var (
	filepath = "./log/"
	filename = time.Now().Format("060102150405") + "_" + localIP() + ".log"
	opt      = Options{
		isPrint:     true,
		isWroteFile: true,
		flushSec:    defaultFlushSec,
		filePath:    filepath,
		fileName:    filename,
	}
	vlog Logger = newLogger(opt)
)

func localIP() string {
	localAddrs, err := net.InterfaceAddrs()
	if err != nil {
		println(err.Error())
	}
	var ip = "localhost"
	for _, address := range localAddrs {
		ipn, ok := address.(*net.IPNet)
		if !ok {
			continue
		}
		if ipn.IP.IsLoopback() {
			continue
		}
		if ipn.IP.To4() != nil {
			ip = ipn.IP.String()
			break
		}
	}
	return ip
}

func Printf(format string, v ...interface{}) {
	vlog.Printf(format, v...)
}

func Infof(format string, v ...interface{}) {
	vlog.Infof(format, v...)
}

func Debugf(format string, v ...interface{}) {
	vlog.Debugf(format, v...)
}

func Warnf(format string, v ...interface{}) {
	vlog.Warnf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	vlog.Errorf(format, v...)
}

type logger struct {
	opt  Options
	file *os.File
	wr   *bufio.Writer
	buf  *bytes.Buffer
	mu   sync.Mutex
}

func NewLogger(opts ...Option) Logger {
	option := Options{}
	for _, opt := range opts {
		option = opt(option)
	}

	return newLogger(option)
}

func newLogger(opt Options) Logger {
	l := &logger{
		opt: opt,
	}
	if !l.opt.isWroteFile {
		return l
	}
	if l.opt.filePath == "" {
		l.opt.filePath = filepath
	}
	if l.opt.fileName == "" {
		l.opt.fileName = filename
	}
	if l.opt.flushSec == 0 {
		l.opt.flushSec = defaultFlushSec
	}
	if !strings.HasSuffix(l.opt.filePath, string(os.PathSeparator)) {
		l.opt.filePath += string(os.PathSeparator)
	}
	l.mkdir()
	f, err := os.OpenFile(l.opt.filePath+l.opt.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	l.buf = &bytes.Buffer{}
	l.file = f
	l.wr = bufio.NewWriter(f)
	go l.backgroundWrite()
	return l
}

func (l *logger) Printf(format string, v ...interface{}) {
	l.print("", format, v...)
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.print("[INFO]", format, v...)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.print("[DEBUG]", format, v...)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.print("[WARN]", format, v...)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.print("[ERROR]", format, v...)
}

func (l *logger) Close() {
	if l.opt.isWroteFile {
		l.write()
		_ = l.wr.Flush()
		_ = l.file.Close()
	}
}

func (l *logger) print(head, format string, v ...interface{}) {
	if !l.opt.isPrint {
		return
	}

	format = fmt.Sprintf(format, v...)
	if head != "" {
		format = fmt.Sprintf("%s %s: %s", head, l.now(), format)
	}
	println(format)

	if !l.opt.isWroteFile {
		return
	}
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	l.mu.Lock()
	_, err := l.buf.WriteString(format)
	if err != nil {
		println(err)
		l.mu.Unlock()
		return
	}
	if l.buf.Len() >= writeByteSize {
		l.write()
	}
	l.mu.Unlock()
}

func (l *logger) now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (l *logger) backgroundWrite() {
	t := time.NewTicker(time.Duration(l.opt.flushSec) * time.Second)
	for range t.C {
		l.mu.Lock()
		l.write()
		l.mu.Unlock()
	}
}

func (l *logger) write() {
	if l.buf.Len() <= 0 {
		return
	}
	_, err := l.wr.WriteString(l.buf.String())
	if err != nil {
		println(err)
	}
}

func (l *logger) mkdir() {
	f, err := os.Stat(l.opt.filePath)
	if err != nil || f.IsDir() == false {
		if err := os.Mkdir(l.opt.filePath, os.ModePerm); err != nil {
			panic("日志目录创建失败, " + err.Error())
		}
	}
}
