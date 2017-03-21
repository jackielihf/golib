
package log

type ICommonLogger interface {
    
    Debugf(format string, args ...interface{})
    Infof(format string, args ...interface{})
    Printf(format string, args ...interface{})
    Warnf(format string, args ...interface{})
    Warningf(format string, args ...interface{})
    Errorf(format string, args ...interface{})
    Fatalf(format string, args ...interface{})
    Panicf(format string, args ...interface{})

    Debug(args ...interface{})
    Info(args ...interface{})
    Print(args ...interface{})
    Warn(args ...interface{})
    Warning(args ...interface{})
    Error(args ...interface{})
    Fatal(args ...interface{})
    Panic(args ...interface{})

}

// logger instance
var CommonLogger ICommonLogger

func init() {
    CommonLogger = new(SimpleLogger)
}

func Debugf(format string, args ...interface{}) {
    CommonLogger.Debugf(format, args...)
}
func Infof(format string, args ...interface{}) {
    CommonLogger.Infof(format, args...)
}
func Printf(format string, args ...interface{}) {
    CommonLogger.Printf(format, args...)
}
func Warnf(format string, args ...interface{}) {
    CommonLogger.Warnf(format, args...)
}
func Warningf(format string, args ...interface{}) {
    CommonLogger.Warningf(format, args...)
}
func Errorf(format string, args ...interface{}) {
    CommonLogger.Errorf(format, args...)
}
func Fatalf(format string, args ...interface{}) {
    CommonLogger.Fatalf(format, args...)
}
func Panicf(format string, args ...interface{}) {
    CommonLogger.Panicf(format, args...)
}


func Debug(args ...interface{}) {
    CommonLogger.Debug(args...)
}
func Info(args ...interface{}) {
    CommonLogger.Info(args...)    
}
func Print(args ...interface{}) {
    CommonLogger.Print(args...)    
}
func Warn(args ...interface{}) {
    CommonLogger.Warn(args...)    
}
func Warning(args ...interface{}) {
    CommonLogger.Warning(args...)    
}
func Error(args ...interface{}) {
    CommonLogger.Error(args...)    
}
func Fatal(args ...interface{}) {
    CommonLogger.Fatal(args...)    
}
func Panic(args ...interface{}) {
    CommonLogger.Panic(args...)    
}


