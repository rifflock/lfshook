// Package LFShook allows users to write to the logfiles using logrus.
package lfshook

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"log"
	"os"
)

// Map for linking a log level to a log file
// Multiple levels may share a file, but multiple files may not be used for one level
type PathMap map[logrus.Level]string

// Hook to handle writing to local log files.
type lfsHook struct {
	paths  PathMap
	levels []logrus.Level
}

// Given a map with keys equal to log levels.
// We can generate our levels handled on the fly, and write to a specific file for each level.
// We can also write to the same file for all levels. They just need to be specified.
func NewHook(levelMap PathMap) *lfsHook {
	hook := &lfsHook{
		paths: levelMap,
	}
	for level, _ := range levelMap {
		hook.levels = append(hook.levels, level)
	}
	return hook
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
