/*
 *@Description
 *@author          lirui
 *@create          2020-09-27 14:43
 */
package broadcaster

import (
	"errors"
	"sync"
)

type MessageType string

// Reference outside of package
const (
	Log              MessageType = "Log"
	NamespaceCreated MessageType = "NamespaceCreated"
	NamespaceUpdated MessageType = "NamespaceUpdated"
	NamespaceDeleted MessageType = "NamespaceDeleted"
)

// 封装信息
type SocketData struct {
	MessageType MessageType
	Payload     interface{}
}

// 广播定义
// Only a pointer to the struct should be used
type Broadcaster struct {
	expired bool
	// Explicit name to specify locking condition
	expiredLock sync.Mutex
	subscribers *sync.Map //map[*Subscriber]struct{}
	c           chan SocketData
}

// 订阅定义
// Wrapper return type for subscriptions
type Subscriber struct {
	subChan   chan SocketData
	unsubChan chan struct{}
}

// Read-Only access to the subscription channel
// Open or nil, never closed
func (s *Subscriber) SubChan() <-chan SocketData {
	return s.subChan
}

// Closed on unsubscribe or when broadcast parent channel closes
func (s *Subscriber) UnsubChan() <-chan struct{} {
	return s.unsubChan
}

var expiredError error = errors.New("Broadcaster expired")

// Creates broadcastertest from channel parameter and immediately starts broadcasting
// Without any subscribers, received data will be discarded
// Broadcaster should be the only channel reader
func NewBroadcaster(c chan SocketData) *Broadcaster {
	if c == nil {
		panic("Channel passed cannot be nil")
	}

	// 初始化一个广播器
	b := &Broadcaster{subscribers: new(sync.Map)}
	// 配置 SocketData
	b.c = c

	go func() {
		for {
			// msg是SocketData通道
			msg, channelOpen := <-b.c
			if channelOpen {
				// 通道打开
				println("SocketData通道打开了")
				// range函数
				// Range calls f sequentially for each key and value present in the map.
				// If f returns false, range stops the iteration.
				// 将map中的简直对依次作为参数传入到其 函数类型的参数中执行 如果返回false就会停止迭代
				b.subscribers.Range(func(key, value interface{}) bool {
					//
					subscriber := key.(*Subscriber)
					select {
					// 将数据给订阅者的接收数据通道
					case subscriber.subChan <- msg:
					case <-subscriber.unsubChan:
					}
					return true
				})
			} else {

				// 通道关闭
				println("SocketData通道关闭了")
				b.expiredLock.Lock()
				b.expired = true
				b.subscribers.Range(func(key, value interface{}) bool {
					subscriber := key.(*Subscriber)
					// 关闭一个接收着通道
					close(subscriber.unsubChan)
					return true
				})
				// Remove references
				b.subscribers = nil
				b.expiredLock.Unlock()
				return
			}
		}
	}()
	return b
}

func (b *Broadcaster) Expired() bool {
	b.expiredLock.Lock()
	defer b.expiredLock.Unlock()
	return b.expired
}

// Subscriber expected to constantly consume or unsubscribe
func (b *Broadcaster) Subscribe() (*Subscriber, error) {
	b.expiredLock.Lock()
	defer b.expiredLock.Unlock()

	if b.expired {
		return &Subscriber{}, expiredError
	}
	newSub := &Subscriber{
		subChan:   make(chan SocketData),
		unsubChan: make(chan struct{}),
	}
	// Generate unique key
	b.subscribers.Store(newSub, struct{}{})
	return newSub, nil
}

func (b *Broadcaster) Unsubscribe(sub *Subscriber) error {
	b.expiredLock.Lock()
	defer b.expiredLock.Unlock()

	if b.expired {
		return expiredError
	}
	if _, ok := b.subscribers.Load(sub); ok {
		b.subscribers.Delete(sub)
		close(sub.unsubChan)
		return nil
	}
	return errors.New("Subscription not found")
}

// Iterates over sync.Map and returns number of elements
// Response can be oversized if counted subscriptions are cancelled while counting
func (b *Broadcaster) PoolSize() (size int) {
	b.expiredLock.Lock()
	defer b.expiredLock.Unlock()

	if b.expired {
		return 0
	}
	b.subscribers.Range(func(key, value interface{}) bool {
		size++
		return true
	})
	return size
}
