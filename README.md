# Local Filesystem Hook for Logrus

Sometimes developers like to write directly to a file on the filesystem. This is a hook for [logrus](https://github.com/Sirupsen/logrus) designed to allow users to do just that.  The log levels are dynamic at instanciation of the hook, so it is capable of logging at some or all levels.

## Example
```go
import (
	log "github.com/Sirupsen/logrus"
	"github.com/rifflock/lfshook"
)

var Log *log.Logger

func NewLogger( config map[string]interface{} ) *log.Logger {
	if Log != nil {
		return Log
	}
	
	Log = log.New()
	Log.Formatter = new(log.JSONFormatter)
	Log.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
		log.InfoLevel : "/var/log/info.log",
		log.ErrorLevel : "/var/log/error.log",
	}))
	return Log
}
```

### Note:
Whichever user is running the go application must have read/write permissions to the log files selected, or if the files do not yet exists, then to the directory in which the files will be created.