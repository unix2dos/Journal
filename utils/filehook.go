package utils

import (
	"os"
	"path/filepath"

	"fmt"

	"time"

	"runtime"

	"bufio"

	"os/signal"

	"syscall"

	"github.com/sirupsen/logrus"
)

type FileHook struct {
	levels   []logrus.Level
	path     string
	link     string
	currSize int64
	maxSize  int64
	file     *os.File
	buf      *bufio.Writer
	queue    chan *LogNode
	today    int
}

type LogNode struct {
	level   logrus.Level
	time    time.Time
	file    string
	line    int
	content string
}

func (f *FileHook) Levels() []logrus.Level {
	return f.levels
}

func (f *FileHook) Fire(entry *logrus.Entry) (err error) {

	node := &LogNode{}
	node.level = entry.Level
	node.time = entry.Time
	node.content = fmt.Sprintf("MODULE=%s|  %v \n", entry.Data["MODULE"], entry.Message)
	_, node.file, node.line, _ = runtime.Caller(5)
	node.file = filepath.Base(node.file)
	f.queue <- node

	return
}

func CreateFileHook() (f *FileHook, err error) {
	f = new(FileHook)

	levels := []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
	f.levels = levels
	f.path = "./logs"
	f.link = f.path + "/main.log"
	f.currSize = 0
	f.maxSize = 268435456

	if err := f.loadFile(); err != nil { //load软链接, 没有就新建当天日志
		if err = f.createFile(); err != nil {
			return nil, err
		}
	}

	f.queue = make(chan *LogNode, 1024)
	go f.run()
	go f.signalHandler()
	return
}

func (f *FileHook) loadFile() (err error) {
	filename, err := os.Readlink(f.link)
	if err != nil {
		return
	}

	newPath := f.path + "/" + filename
	info, err := os.Stat(newPath)
	if err != nil {
		return
	}

	file, err := os.OpenFile(newPath, os.O_WRONLY|os.O_SYNC|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	f.file = file
	f.currSize = info.Size()
	year, mon, day := info.ModTime().Date()
	f.today = year*1000 + int(mon)*100 + day
	return
}

func (f *FileHook) createFile() (err error) {
	err = os.MkdirAll(f.path, os.ModePerm)
	if err != nil {
		return
	}

	now := time.Now()
	newPath := f.path + "/" + now.Format("20060102_150405") + ".log"

	var currSize int64 = 0
	info, err := os.Stat(newPath)
	if err == nil {
		currSize = info.Size()
	}

	file, err := os.OpenFile(newPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_SYNC, 0666)
	if err != nil {
		return
	}

	f.file = file
	f.currSize = currSize
	f.createLink(newPath)
	year, mon, day := now.Date()
	f.today = year*1000 + int(mon)*100 + day
	return
}

func (f *FileHook) createLink(filename string) (err error) {
	//fmt.Println("createLink", filepath.Base(filename))
	os.Remove(f.link)
	return os.Symlink(filepath.Base(filename), f.link)
}

func (f *FileHook) isSwitchFile(time *time.Time) bool {
	year, mon, day := time.Date()
	today := year + int(mon) + day
	if f.currSize < f.maxSize && today <= f.today {
		return false
	}
	return true
}

func (f *FileHook) run() {
	ticker := time.NewTicker(time.Second)
	f.buf = bufio.NewWriterSize(f.file, 1048576)

	for {
		select {
		case node := <-f.queue:
			if node.level <= logrus.FatalLevel {
				f.buf.Flush()
				return
			}
			//是否该换文件了
			if f.isSwitchFile(&node.time) {
				f.buf.Flush()
				if f.file != nil {
					f.file.Close()
					f.file = nil
				}
				err := f.createFile()
				if err != nil {
					return
				}
				f.buf.Reset(f.file)
			}
			//写日志
			size, _ := f.buf.WriteString(fmt.Sprintf("[%2s][%02d:%02d:%02d][%s:%d]: ",
				logrus.Level(node.level).String(),
				node.time.Hour(),
				node.time.Minute(),
				node.time.Second(),
				node.file,
				node.line))
			f.currSize += int64(size)
			size, _ = f.buf.WriteString(node.content)
			f.currSize += int64(size)

		case <-ticker.C:
			if f.buf.Buffered() > 0 {
				f.buf.Flush()
			}
		}
	}
}

func (f *FileHook) signalHandler() {
	var sigChan = make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case sig := <-sigChan:
			fmt.Printf("received signal is:%v and exit the whole world\n", sig)
			f.buf.Flush()
			os.Exit(1)
		}
	}
}
