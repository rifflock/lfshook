// Package LFShook allows users to write to the logfiles using logrus.
package lfshook

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"log"
	"os"
	"sync"
)

// We are logging to file, strip colors to make the output more readable
var txtFormatter = &logrus.TextFormatter{DisableColors: true}

// Map for linking a log level to a log file
// Multiple levels may share a file, but multiple files may not be used for one level
type PathMap map[logrus.Level]string

// Hook to handle writing to local log files.
type lfsHook struct {
	paths  PathMap
	levels []logrus.Level
	lock   *sync.Mutex
	formatter logrus.Formatter
}

// Given a map with keys equal to log levels.
// We can generate our levels handled on the fly, and write to a specific file for each level.
// We can also write to the same file for all levels. They just need to be specified.
func NewHook(levelMap PathMap) *lfsHook {
	hook := &lfsHook{
		paths: levelMap,
		lock:  new(sync.Mutex),
		formatter: txtFormatter,
	}
	for level, _ := range levelMap {
		hook.levels = append(hook.levels, level)
	}
	return hook
}

func (hook *lfsHook) SetFormatter(formatter logrus.Formatter) {
	hook.formatter = formatter

	switch hook.formatter.(type) {
		case *logrus.TextFormatter:
			textFormatter := hook.formatter.(*logrus.TextFormatter)
			textFormatter.DisableColors = true
	}
}

// Open the file, write to the file, close the file.
// Whichever user is running the function needs write permissions to the file or directory if the file does not yet exist.
func (hook *lfsHook) Fire(entry *logrus.Entry) error {
	var (
		fd   *os.File
		path string
		msg  string
		err  error
		ok   bool
	)

	hook.lock.Lock()
	defer hook.lock.Unlock()

	if path, ok = hook.paths[entry.Level]; !ok {
		err = fmt.Errorf("no file provided for loglevel: %d", entry.Level)
		log.Println(err.Error())
		return err
	}
	fd, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("failed to open logfile:", path, err)
		return err
	}
	defer fd.Close()

	// only modify Formatter if we are using a TextFormatter so we can strip colors
	switch entry.Logger.Formatter.(type) {
	case *logrus.TextFormatter:
		// swap to colorless TextFormatter
		formatter := entry.Logger.Formatter
		entry.Logger.Formatter = txtFormatter
		defer func() {
			// assign back original formatter
			entry.Logger.Formatter = formatter
		}()
	}

	msg, err = entry.String()

	if err != nil {
		log.Println("failed to generate string for entry:", err)
		return err
	}
	fd.WriteString(msg)
	return nil
}

func (hook *lfsHook) Levels() []logrus.Level {
	return hook.levels
}
