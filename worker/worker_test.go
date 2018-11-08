package worker

import (
    "fmt"
    "testing"
    "time"
    "os"
    "os/signal"
    "sync"
)

func waitForQuit() {

    signals := make(chan os.Signal, 1)
    signal.Notify(signals, os.Interrupt, os.Kill)
    <-signals
}

type MyHandler struct {

}

func(h *MyHandler) Handle(input interface{}) interface{} {

    if v, _ := input.(string); v == "slow" {
        time.Sleep(10 * time.Millisecond)
    }else if v == "veryslow" {
        time.Sleep(5 * time.Second)
    }else if v == "exception" {
        var b int
        b = 0
        b = 1 / b
    }

    result := fmt.Sprintf("R_%v", input)
    return result
}

func expect(input interface{}) interface{} {
    result := fmt.Sprintf("R_%v", input)
    return result   
}

func NewMyHandler() Handler {
    return &MyHandler{}
}

func Test_Process(t *testing.T) {

    pool := NewPool(2, NewMyHandler).Start()
    var wg sync.WaitGroup
    max := 10
    wg.Add(max)

    for i := 0; i < max; i++ {
        go func(a int){
            result, _ := pool.Process(a)
            if result != expect(a) {
                t.Error("wrong result of process")
            }
            wg.Done()
        }(i)
    }
    
    wg.Wait()
}

func Test_Close(t *testing.T) {

    pool := NewPool(2, NewMyHandler).Start()
    input := "1"

    result, _ := pool.Process(input)
    if result != expect(input) {
        t.Error("wrong result of process")
    }
            
    pool.Close()

    _, err := pool.Process(input)
    if err != ErrPoolClosed {
        fmt.Println(err)
        t.Error("wrong: pool not closed")
    }
}

func Test_timeout(t *testing.T) {

    pool := NewPool(2, NewMyHandler).WithTimeout(1 * time.Second).Start()

    _, err := pool.Process("veryslow")
    fmt.Println(err)
    if err != ErrJobTimeout {
        t.Error("wrong: should be timeouted")
    }
}

func Test_exception(t *testing.T) {

    pool := NewPool(1, NewMyHandler).WithBufferSize(1).Start()

    _, err := pool.Process("exception")
    fmt.Println(err)
    if err == nil {
        t.Error("wrong: should catch exception")
    }

    // buffer exceed
    var bErr error = nil
    var wg sync.WaitGroup
    max := 4
    wg.Add(max)
    for i := 0; i < max; i++ {
        go func(){
            if _, tmpErr := pool.Process("slow"); tmpErr != nil {
                bErr = tmpErr
            }
            wg.Done()
        }()    
    }
    wg.Wait()
    if bErr != ErrBufferFull {
        t.Error("wrong: should drop item when buffer is full")
    }
}

func Test_parallel(t *testing.T) {

    size := 10
    pool := NewPool(size, NewMyHandler).WithTimeout(1 * time.Second).Start()

    if pool.IdleWorker() != size {
        t.Error("wrong: init size of worker")
    }

    num := 6
    for i := 0; i < num; i++ {
        go func(){
            pool.Process("veryslow")        
        }()    
    }
    time.Sleep(1*time.Second)
    
    fmt.Println(pool.IdleWorker())
    if pool.IdleWorker() > size - num {
        t.Error("wrong: worker not parallel")
    }
}

func Benchmark_1(b *testing.B) {

    pool := NewPool(2, NewMyHandler).Start()

    b.RunParallel(func(pb *testing.PB){
        for pb.Next() {
            // pool.Process("helloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworld")    
            _, err := pool.Process("slow")
            if err != nil {
                fmt.Println(err)
            }
        }
    })
    
    

}

