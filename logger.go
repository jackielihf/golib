
/**
 Logger client for application.
**/
package golib

import "os"
import "sync"
import "encoding/json"
import "bytes"
import "strconv"
import "github.com/Sirupsen/logrus"
import "github.com/rifflock/lfshook"

var onceOfLogger sync.Once
var logger *logrus.Logger
type Fields logrus.Fields


// custom formatter
type MyFormatter struct {
    appName string
}

func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    var b *bytes.Buffer = new(bytes.Buffer)
    // format
    timestampFormat := logrus.DefaultTimestampFormat
    b.WriteString(entry.Time.Format(timestampFormat))
    b.WriteString(" [")
    b.WriteString(entry.Level.String())
    b.WriteString("] ")
    b.WriteString(f.appName)
    b.WriteString(" $ ")
    b.WriteString(entry.Message)

    //buffer
    if entry.Buffer != nil {
        b.WriteString(entry.Buffer.String())
    }
    //entry.Data
    if len(entry.Data) > 0 {
        serialized, err := json.Marshal(entry.Data)
        if err == nil {
            b.Write(serialized)
        }    
    }
    b.WriteByte('\n')
  return b.Bytes(), nil
}

// logger
func Init(){

    appName := os.Getenv("app_name")
    logFormat := os.Getenv("log_format")
    enableFile := os.Getenv("log_enable_file")
    logDir := os.Getenv("log_dir")
    // enableSyslog := os.Getenv("log_enable_syslog")

    if logFormat == "json" {
        logger.Formatter = new(logrus.JSONFormatter)
    }else {
        logger.Formatter = &MyFormatter{appName}
    }
    // file
    if v, err := strconv.ParseBool(enableFile); v && err == nil {
        if logDir == "" {
            logDir = "."
        }
        logger.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
            logrus.InfoLevel: logDir + "/info.log",
        }))
    }
    // syslog
    
}
func Logger() *logrus.Logger{
    // init once
    onceOfLogger.Do(func(){
        logger = logrus.New()
        Init()
    
    })
    return logger
}



