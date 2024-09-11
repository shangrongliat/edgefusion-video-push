package service

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/robfig/cron"
)

// Queue 结构体定义了一个线程安全的队列，并增加了一个用于通知的channel
type Queue struct {
	mu         sync.Mutex
	StatusChan chan string // 状态通知通道
	status     int32       // 0 - Active, 1 - Idle
}

const (
	Active = iota
	Idle
)

// NewQueue 创建一个新的队列实例，并初始化itemChan
func NewQueue() *Queue {
	q := Queue{
		StatusChan: make(chan string), // 状态通知通道
	}
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
	atomic.StoreInt32(&q.status, Active)
	go q.listenForData()
	c.Start()
	log.Println("队列初始化成功....")
	return &q
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
		}
	}
}

// Put 向队列中添加一个元素，并通过channel发送信号
func (q *Queue) Put() {
	if q.status == Idle {
		atomic.StoreInt32(&q.status, Active)
	}
}

func (q *Queue) Status() int32 {
	return q.status
}

func (q *Queue) Close() {
	<-q.StatusChan // 等待状态通道关闭
}
