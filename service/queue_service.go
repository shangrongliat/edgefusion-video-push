package service

import (
	"container/list"
	"github.com/robfig/cron"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// Queue 结构体定义了一个线程安全的队列，并增加了一个用于通知的channel
type Queue struct {
	mu         sync.Mutex
	queue      *list.List
	timer      *time.Timer
	DataChan   chan interface{}
	StatusChan chan string // 状态通知通道
	status     int32       // 0 - Active, 1 - Idle
}

const (
	Active = iota
	Idle
)

// NewQueue 创建一个新的队列实例，并初始化itemChan
func NewQueue() *Queue {
	q := &Queue{
		DataChan:   make(chan interface{}), // 数据通道的缓冲区大小
		StatusChan: make(chan string),      // 状态通知通道
		queue:      list.New(),
	}
	atomic.StoreInt32(&q.status, Active)
	c := cron.New()
	if err := c.AddFunc("@every 20s", func() {
		q.mu.Lock()
		defer q.mu.Unlock()
		if atomic.LoadInt32(&q.status) == Active {
			atomic.StoreInt32(&q.status, Idle)
			q.StatusChan <- "Idle" // 通知状态变为 Idle
		}
	}); err != nil {
		log.Println("定时队列状态监测启动失败....", err)
		return nil
	}
	go q.listenForData()
	c.Start()
	return q
}

func (q *Queue) listenForData() {
	for {
		select {
		case status, ok := <-q.StatusChan:
			if !ok {
				// 数据通道关闭，退出
				return
			}
			if status == "Active" {
				atomic.StoreInt32(&q.status, Active)
			} else if status == "Idle" {
				atomic.StoreInt32(&q.status, Idle)
			}
			log.Println("定时队列活跃状态监测2:", q.status)
		}
	}
}

// Put 向队列中添加一个元素，并通过channel发送信号
func (q *Queue) Put(data interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue.PushBack(data)
	if q.status == Idle {
		q.StatusChan <- "Active" // 通知状态变为 Active
	}
	q.DataChan <- struct{}{}
}

// Pull 从队列中取出一个元素并移除
func (q *Queue) Pull() (interface{}, bool) {
	if element := q.queue.Front(); element != nil {
		q.queue.Remove(element)
		return element.Value, true
	}
	return nil, false
}

func (q *Queue) Status() int32 {
	return q.status
}

func (q *Queue) Close() {
	close(q.DataChan) // 关闭数据通道
	<-q.StatusChan    // 等待状态通道关闭
}
