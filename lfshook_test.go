package lfshook

import (
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

const expectedMsg = "This is the expected test message."

// Tests that writing to a tempfile log works.
// Matches the 'msg' of the output and deletes the tempfile.
func TestLogEntryWritten(t *testing.T) {
	log := logrus.New()
	// The colors were messing with the regexp so I turned them off.
	log.Formatter = &logrus.TextFormatter{DisableColors: true}
	tmpfile, err := ioutil.TempFile("", "test_lfshook")
	if err != nil {
		t.Errorf("Unable to generate logfile due to err: %s", err)
	}
	fname := tmpfile.Name()
	defer func() {
		tmpfile.Close()
		os.Remove(fname)
	}()
	hook := NewHook(PathMap{
		logrus.InfoLevel: fname,
	})
	log.Hooks.Add(hook)

	log.Info(expectedMsg)

	if contents, err := ioutil.ReadAll(tmpfile); err != nil {
		t.Errorf("Error while reading from tmpfile: %s", err)
	} else if matched, err := regexp.Match("msg=\""+expectedMsg+"\"", contents); err != nil || !matched {
		t.Errorf("Message read (%s) doesnt match message written (%s) for file: %s", contents, expectedMsg, fname)
	}
}
