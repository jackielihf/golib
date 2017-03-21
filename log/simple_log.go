
package log

import "fmt"

type SimpleLogger struct {

}

func (that *SimpleLogger) Debugf(format string, args ...interface{}) {
    fmt.Printf(format + "\n", args...)
}
func (that *SimpleLogger) Infof(format string, args ...interface{}) {
    fmt.Printf(format + "\n", args...)
}
func (that *SimpleLogger) Printf(format string, args ...interface{}) {
    fmt.Printf(format + "\n", args...)
}
func (that *SimpleLogger) Warnf(format string, args ...interface{}) {
    fmt.Printf(format + "\n", args...)
}
func (that *SimpleLogger) Warningf(format string, args ...interface{}) {
    fmt.Printf(format + "\n", args...)
}
func (that *SimpleLogger) Errorf(format string, args ...interface{}) {
    fmt.Printf(format + "\n", args...)
}
func (that *SimpleLogger) Fatalf(format string, args ...interface{}) {
    fmt.Printf(format + "\n", args...)
}
func (that *SimpleLogger) Panicf(format string, args ...interface{}) {
    fmt.Printf(format + "\n", args...)
}


func (that *SimpleLogger) Debug(args ...interface{}) {
    fmt.Println(args...)
}
func (that *SimpleLogger) Info(args ...interface{}) {
    fmt.Println(args...)
}
func (that *SimpleLogger) Print(args ...interface{}) {
    fmt.Println(args...)
}
func (that *SimpleLogger) Warn(args ...interface{}) {
    fmt.Println(args...)
}
func (that *SimpleLogger) Warning(args ...interface{}) {
    fmt.Println(args...)
}
func (that *SimpleLogger) Error(args ...interface{}) {
    fmt.Println(args...)
}
func (that *SimpleLogger) Fatal(args ...interface{}) {
    fmt.Println(args...)
}
func (that *SimpleLogger) Panic(args ...interface{}) {
    fmt.Println(args...)
}




