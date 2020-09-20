/**
 * yy captain
 * time: 2020/09/20
 **/

package mylog

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type logger struct {
	fileName string     /** 文件名 **/
	file     *os.File   /** 文件句柄 **/
	mu       sync.Mutex /** 互斥量 **/
}

/**
 * 回滚日志
 **/
func (l *logger) rotate() {
	if l.file != nil {
		l.file.Close()
		l.file = nil
	}

	now := time.Now().Format("20060102150405")
	os.Rename(l.fileName, l.fileName+"-"+now)
}

/**
 * 打开一个新的日志文件
 **/
func (l *logger) open() (err error) {
	l.file, err = os.OpenFile(l.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open LogFile %v failed: %v", l.fileName, err)
	}
	return nil
}

/**
 * 向日志写入内容
 **/
func (l *logger) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil {
		if err = l.open(); err != nil {
			return 0, err
		}
	}

	sizeLog, err := l.file.Write(p)
	if err != nil {
		return -1, err
	}
	return sizeLog, nil
}

/**
 * 定时器，定期对日志文件进行回滚
 **/
func rotateTimer(l *logger, maxAge int) {
	for {
		now := time.Now().Unix()
		next := (now/int64(maxAge) + 1) * int64(maxAge)

		timer := time.NewTimer(time.Second * time.Duration(next-now))
		<-timer.C

		l.mu.Lock()
		l.rotate()
		l.mu.Unlock()
	}
}

/**
 * 创建一个日志句柄
 **/
func newlogger(fileName string, maxAge int) *logger {
	logger := &logger{
		fileName: fileName,
	}

	if maxAge > 0 {
		go rotateTimer(logger, maxAge)
	}

	return logger
}
