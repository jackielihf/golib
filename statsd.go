
// statsd client
package golib
import "fmt"
import "sync"
import "os"
import statsd "gopkg.in/alexcesaro/statsd.v2"


var onceOfStatsd sync.Once
var client *statsd.Client

func Statsd() *statsd.Client{
    onceOfStatsd.Do(func(){
        addr := os.Getenv("statsd_addr")
        var err error
        if addr != "" {
            client, err = statsd.New(statsd.Address(addr))
            if err != nil {
                defer client.Close()
                fmt.Println("failed to create statsd client: " + err.Error())
            }
        }else{
            client, _ = statsd.New()  // default: 8125
        }
    })
    return client
}

