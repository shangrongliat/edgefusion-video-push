package service

import (
	"container/list"
	"log"
	"sync"
)

// Queue 结构体定义了一个线程安全的队列，并增加了一个用于通知的channel
type Queue struct {
	mu       sync.Mutex
	queue    *list.List
	ItemChan chan interface{}
}

// NewQueue 创建一个新的队列实例，并初始化itemChan
func NewQueue() *Queue {
	return &Queue{
		queue:    list.New(),
		ItemChan: make(chan interface{}),
	}
}

// Put 向队列中添加一个元素，并通过channel发送信号
func (q *Queue) Put(data interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue.PushBack(data)
	select {
	case q.ItemChan <- struct{}{}: // 通知有新元素加入
	default: // 避免channel满时阻塞
	}
}

// Pull 从队列中取出一个元素并移除
func (q *Queue) Pull() (interface{}, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if element := q.queue.Front(); element != nil {
		q.queue.Remove(element)
		log.Fatalln("消费并移除的应用名称：", element)
		return element.Value, true
	}
	return nil, false
}
