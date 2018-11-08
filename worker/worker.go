
package worker

import (
    "fmt"
    "errors"
    "sync"
    "sync/atomic"
    "runtime/debug"
    "time"
)


var (
    ErrJobTimeout = errors.New("err: job handling timeout")
    ErrWorkerClosed = errors.New("err: workder already closed")
    ErrBufferFull = errors.New("err: buffer full, request dropped")
    ErrResultChanClosed = errors.New("err: resultChan closed")
    ErrPoolClosed = errors.New("err: pool closed")
    ErrJobInterrupt = errors.New("err: job interrupt")
)

type Handler interface {

    Handle(interface{}) interface{} // 处理函数
}

type Result struct {
    Data interface{}
    Err error
}

type WorkItem struct {
    payload interface{}
    resultChan chan Result
    interrupt chan bool
}

func NewWorkItem(payload interface{}) *WorkItem {
    item := &WorkItem{
        payload: payload,
        resultChan: make(chan Result),
        interrupt: make(chan bool),
    }
    return item
}

func(wi *WorkItem) release(){
    // fmt.Println("release item")
    wi.payload = nil
    close(wi.interrupt)
}

type Worker struct {
    id string
    pool *Pool
    handler Handler

}

func NewWorker(id string, pool *Pool, handler Handler) *Worker{

    wk := &Worker{
        id: id,
        pool: pool,
        handler: handler,
    }
    return wk
}


func(w *Worker) run(item *WorkItem){

    if item == nil {
        return
    }

    go func(){

        defer func(){
            v := recover()
            if v != nil {
                fmt.Printf("worker err recover: %v. print stack:\n", v)
                debug.PrintStack()
                item.resultChan <- Result{Data: nil, Err: errors.New(fmt.Sprintf("worker err recover: %v\n", v))}
            }
            // 将空闲worker加入pool
            w.pool.workerChan <- w
        }() 

        output := w.handler.Handle(item.payload)
        // fmt.Printf("worker %s, input: %v, output: %v\n", w.id, item.payload, output)

        select {
        case <-item.interrupt:
        case item.resultChan <- Result{Data: output, Err: nil}:    
        }
    }()    
}

func(w *Worker) interrupt(){
    
}



type Pool struct {

    workers []*Worker 
    workerChan chan *Worker
    queue chan *WorkItem
    quit chan bool

    poolSize int
    queueSize int64 // queue的实时大小
    bufferSize int64 // queue buffer大小
    mtx sync.Mutex
    timeout time.Duration
}

func NewPool(poolSize int, builder func() Handler) *Pool {

    if poolSize < 1 {
        poolSize = 1
    }

    bufferSize := int64(poolSize * 1000)

    pool := &Pool{
        workers: make([]*Worker, poolSize),
        workerChan: make(chan *Worker, poolSize),
        queue: make(chan *WorkItem, bufferSize),
        quit: make(chan bool),
        poolSize: poolSize,
        queueSize: 0,
        bufferSize: bufferSize,
    }

    for i := 0; i < pool.poolSize; i ++ {
        id := fmt.Sprintf("%d", i)
        pool.workers[i] = NewWorker(id, pool, builder())
    }

    fmt.Printf("Pool created: poolSize=%d\n", pool.poolSize)
    return pool
}

// 在Start前调用
func(p *Pool) WithBufferSize(bs int64) *Pool {
    p.mtx.Lock()
    defer p.mtx.Unlock()

    if bs < 1 {
        bs = 1
    }
    p.bufferSize = bs
    p.queue = make(chan *WorkItem, bs)
    return p
}


func(p *Pool) WithTimeout(d time.Duration) *Pool {
    p.mtx.Lock()
    defer p.mtx.Unlock()

    if d < time.Millisecond {
        d = time.Millisecond
    }
    p.timeout = d
    return p
}

func(p *Pool) QueueSize() int64 {
    return atomic.LoadInt64(&p.queueSize)
}

func(p *Pool) IdleWorker() int {
    p.mtx.Lock()
    defer p.mtx.Unlock()
    return len(p.workerChan)
}

func(p *Pool) Start() *Pool {
    p.mtx.Lock()
    defer p.mtx.Unlock()

    for i := 0; i < p.poolSize; i++ {
        wk := p.workers[i]
        p.workerChan <- wk
        fmt.Printf("start worker[%s]\n", wk.id)
    }
    p.dispatch()
    return p
}

func(p *Pool) Close() {
    p.mtx.Lock()
    defer p.mtx.Unlock()
    
    close(p.quit)
    p.poolSize = 0
    p.workers = nil
}

func(p *Pool) dispatch() {

    go func(){
        for {
            select {
            case <-p.quit:
                return
            case item := <- p.queue:
                atomic.AddInt64(&p.queueSize, -1)
                select {
                case <-p.quit:
                    return
                case worker := <-p.workerChan: // 从worker池取一个空闲worker
                    worker.run(item)
                }
            }
        }
        fmt.Println("pool dispatch stopped")
    }()
}


func(p *Pool) enqueueTimed(item *WorkItem, timer *time.Timer) error {

    if timer != nil {
        select {
        case <-p.quit:
            return ErrPoolClosed
        case <-timer.C:
            return ErrJobTimeout
        case p.queue <- item:
            atomic.AddInt64(&p.queueSize, 1)
        default: // 若buffer满, 则丢弃item
            return ErrBufferFull
        }        
    }else{
        select {
        case <-p.quit:
            return ErrPoolClosed
        case p.queue <- item:
            atomic.AddInt64(&p.queueSize, 1)
        default: // 若buffer满, 则丢弃item
            return ErrBufferFull
        }
    }
    return nil
}

func(p *Pool) Process(payload interface{}) (interface{}, error) {


    select {
    case <-p.quit:
        return nil, ErrPoolClosed
    default:

    }

    var timer *time.Timer = nil
    if p.timeout > time.Millisecond { // 启用超时
        timer = time.NewTimer(p.timeout)
    }

    item := NewWorkItem(payload)
    defer func(){
        if item != nil {
            item.release()
        }
        if timer != nil {
            timer.Stop()
        }
    }()

    // enqueue
    if err := p.enqueueTimed(item, timer); err != nil {
        return nil, err
    }

    var result Result
    var open bool

    // 等待结果
    if timer != nil {
        select {
        case result, open = <- item.resultChan:
        case <-timer.C:
            return nil, ErrJobTimeout            
        }    
    }else{
        result, open = <- item.resultChan    
    }
    
    if !open {
        return nil, ErrResultChanClosed
    }
    if &result == nil {
        return nil, nil
    }
    return result.Data, result.Err
}

// no block
// 只保证发送成功，不等待结果返回
func(p *Pool) ProcessNB(payload interface{}) (error) {

    select {
    case <-p.quit:
        return ErrPoolClosed
    default:
    }

    item := NewWorkItem(payload)
    defer func(){
        if item != nil {
            item.release()
        }
    }()

    // enqueue
    if err := p.enqueueTimed(item, nil); err != nil {
        return err
    }

    return nil
}


