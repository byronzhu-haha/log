package log

type Options struct {
	isPrint     bool
	isWroteFile bool
	flushSec    int32
	filePath    string
	fileName    string
}

type Option func(o Options) Options

func OpenPrint() Option {
	return func(o Options) Options {
		o.isPrint = true
		return o
	}
}

func OpenWriteFile() Option {
	return func(o Options) Options {
		o.isWroteFile = true
		return o
	}
}

func Filepath(path string) Option {
	return func(o Options) Options {
		if path == "" {
			o.filePath = filepath
			return o
		}
		o.filePath = path
		return o
	}
}

func FileName(name string) Option {
	return func(o Options) Options {
		if name == "" {
			o.fileName = filename
			return o
		}
		o.fileName = name
		return o
	}
}

func FlushSec(sec int32) Option {
	return func(o Options) Options {
		if sec == 0 {
			o.flushSec = defaultFlushSec
			return o
		}
		o.flushSec = sec

		return o
	}
}
