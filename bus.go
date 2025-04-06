package cryo

import (
	"sync"
)

// EventBus 是事件总线，负责管理中间件和事件的分发
type EventBus struct {
	middlewareMutex      sync.RWMutex // 保护中间件的读写锁
	preMiddleware        []Middleware // 预处理中间件列表
	processingMiddleware []Middleware // 中间件列表
	postMiddleware       []Middleware // 后处理中间件列表
}

// NewEventBus 创建一个新的事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		preMiddleware:        make([]Middleware, 0),
		processingMiddleware: make([]Middleware, 0),
		postMiddleware:       make([]Middleware, 0),
	}
}

func (bus *EventBus) applyPreMiddleware(event Event) Event {
	if len(bus.preMiddleware) == 0 {
		return event // 如果没有预处理中间件，则直接返回事件
	}
	eventType := event.GetEventType() // 获取事件类型

	bus.middlewareMutex.RLock() // 持有读锁
	// 创建一个中间件切片的副本，减小锁的粒度
	var middlewareCopy []Middleware
	if len(bus.preMiddleware) > 0 {
		middlewareCopy = make([]Middleware, len(bus.preMiddleware))
		copy(middlewareCopy, bus.preMiddleware)
	}
	bus.middlewareMutex.RUnlock() // 释放读锁

	// 预处理中间件是按顺序执行的
	for _, middleware := range middlewareCopy {
		if middleware.IsGlobal() || middleware.HasType(eventType) {
			if !middleware.Do(event) {
				return nil // 事件被中间件截断
			}
		}
	}
	return event // 返回事件
}

// applyProcessingMiddleware 应用处理中间件
func (bus *EventBus) applyProcessingMiddleware(event Event) Event {
	if len(bus.processingMiddleware) == 0 {
		return event // 如果没有处理中间件，则直接返回事件
	}
	eventType := event.GetEventType() // 获取事件类型
	eventCopy := event                // 事件的副本

	bus.middlewareMutex.RLock() // 持有读锁
	// 创建一个中间件切片的副本，减小锁的粒度
	var middlewareCopy []Middleware
	if len(bus.processingMiddleware) > 0 {
		middlewareCopy = make([]Middleware, len(bus.processingMiddleware))
		copy(middlewareCopy, bus.processingMiddleware)
	}
	bus.middlewareMutex.RUnlock() // 释放读锁

	// for _, middleware := range middlewareCopy {
	// 	// 为每个中间件创建一个 goroutine
	// 	go func() {
	// 		if middleware.IsGlobal() || middleware.HasType(eventType) {
	// 			middleware.Do(event)
	// 		}
	// 	}()
	// }

	// 让Claude Sonnet 3.7优化了一下循环里面的闭包问题
	// 使用WaitGroup等待所有goroutine完成
	var wg sync.WaitGroup

	// 每个处理中间件之间是并发执行的，但中间件内部仍然是顺序执行的
	for i := range middlewareCopy {
		middleware := middlewareCopy[i] // 在循环内部创建局部变量，避免闭包陷阱

		if middleware.IsGlobal() || middleware.HasType(eventType) {
			wg.Add(1)
			// 为每个中间件创建一个 goroutine
			go func(m Middleware) {
				defer wg.Done()
				m.Do(event)
			}(middleware) // 传递中间件实例作为参数
		}
	}

	// 等待所有中间件处理完成
	// wg.Wait()

	// 处理中间件是用来替代事件订阅的，它不会对事件进行修改，只是用来把事件分发到对应的处理器
	// 自然也不需要等待所有中间件处理完成，直接返回原来的事件即可
	return eventCopy // 返回事件副本
}

// applyPostMiddleware 应用后处理中间件
func (bus *EventBus) applyPostMiddleware(event Event) Event {
	if len(bus.postMiddleware) == 0 {
		return event // 如果没有后处理中间件，则直接返回事件
	}
	eventType := event.GetEventType() // 获取事件类型

	bus.middlewareMutex.RLock() // 持有读锁
	// 创建一个中间件切片的副本，减小锁的粒度
	var middlewareCopy []Middleware
	if len(bus.postMiddleware) > 0 {
		middlewareCopy = make([]Middleware, len(bus.postMiddleware))
		copy(middlewareCopy, bus.postMiddleware)
	}
	bus.middlewareMutex.RUnlock() // 释放读锁

	// 后处理中间件是按顺序执行的
	for _, middleware := range middlewareCopy {
		if middleware.IsGlobal() || middleware.HasType(eventType) {
			if !middleware.Do(event) {
				return nil // 事件被中间件截断
			}
		}
	}
	return event // 返回事件
}

// AddPreMiddleware 添加预处理中间件
func (bus *EventBus) AddPreMiddleware(middleware ...Middleware) {
	bus.middlewareMutex.Lock() // 持有写锁
	defer bus.middlewareMutex.Unlock()
	for _, m := range middleware {
		bus.preMiddleware = append(bus.preMiddleware, m)
	}
}

// AddProcessingMiddleware 添加处理中间件
func (bus *EventBus) AddProcessingMiddleware(middleware ...Middleware) {
	bus.middlewareMutex.Lock() // 持有写锁
	defer bus.middlewareMutex.Unlock()
	for _, m := range middleware {
		bus.processingMiddleware = append(bus.processingMiddleware, m)
	}
}

// AddPostMiddleware 添加后处理中间件
func (bus *EventBus) AddPostMiddleware(middleware ...Middleware) {
	bus.middlewareMutex.Lock() // 持有写锁
	defer bus.middlewareMutex.Unlock()
	for _, m := range middleware {
		bus.postMiddleware = append(bus.postMiddleware, m)
	}
}

// RemoveMiddlewareById 删除指定Id的中间件
func (bus *EventBus) RemoveMiddlewareById(id string) bool {
	if id == "" {
		return false
	}

	bus.middlewareMutex.Lock() // 持有写锁
	defer bus.middlewareMutex.Unlock()

	// 创建一个新的切片来存储不包含指定ID的中间件
	newPreMiddleware := make([]Middleware, 0, len(bus.preMiddleware))
	for _, m := range bus.preMiddleware {
		if m.GetId() != id {
			newPreMiddleware = append(newPreMiddleware, m)
		}
	}
	bus.preMiddleware = newPreMiddleware

	newProcessingMiddleware := make([]Middleware, 0, len(bus.processingMiddleware))
	for _, m := range bus.processingMiddleware {
		if m.GetId() != id {
			newProcessingMiddleware = append(newProcessingMiddleware, m)
		}
	}
	bus.processingMiddleware = newProcessingMiddleware

	newPostMiddleware := make([]Middleware, 0, len(bus.postMiddleware))
	for _, m := range bus.postMiddleware {
		if m.GetId() != id {
			newPostMiddleware = append(newPostMiddleware, m)
		}
	}
	bus.postMiddleware = newPostMiddleware

	return true
}

// RemoveMiddlewareByTag 删除包含所有指定标签的中间件
func (bus *EventBus) RemoveMiddlewareByTag(tag ...string) bool {
	if len(tag) == 0 {
		return false
	}

	bus.middlewareMutex.Lock() // 持有写锁
	defer bus.middlewareMutex.Unlock()

	// 创建一个新的切片来存储不包含指定标签的中间件
	newPreMiddleware := make([]Middleware, 0, len(bus.preMiddleware))
	for _, m := range bus.preMiddleware {
		if !m.HasTag(tag...) {
			newPreMiddleware = append(newPreMiddleware, m)
		}
	}
	bus.preMiddleware = newPreMiddleware

	newProcessingMiddleware := make([]Middleware, 0, len(bus.processingMiddleware))
	for _, m := range bus.processingMiddleware {
		if !m.HasTag(tag...) {
			newProcessingMiddleware = append(newProcessingMiddleware, m)
		}
	}
	bus.processingMiddleware = newProcessingMiddleware

	newPostMiddleware := make([]Middleware, 0, len(bus.postMiddleware))
	for _, m := range bus.postMiddleware {
		if !m.HasTag(tag...) {
			newPostMiddleware = append(newPostMiddleware, m)
		}
	}
	bus.postMiddleware = newPostMiddleware

	return true
}

// ClearPreMiddleware 清空预处理中间件
func (bus *EventBus) ClearPreMiddleware() {
	bus.middlewareMutex.Lock()
	defer bus.middlewareMutex.Unlock()
	bus.preMiddleware = make([]Middleware, 0)
}

// ClearProcessingMiddleware 清空处理中间件
func (bus *EventBus) ClearProcessingMiddleware() {
	bus.middlewareMutex.Lock()
	defer bus.middlewareMutex.Unlock()
	bus.processingMiddleware = make([]Middleware, 0)
}

// ClearPostMiddleware 清空后处理中间件
func (bus *EventBus) ClearPostMiddleware() {
	bus.middlewareMutex.Lock()
	defer bus.middlewareMutex.Unlock()
	bus.postMiddleware = make([]Middleware, 0)
}

// ClearAllMiddleware 清空所有中间件
func (bus *EventBus) ClearAllMiddleware() {
	bus.middlewareMutex.Lock()
	defer bus.middlewareMutex.Unlock()
	bus.preMiddleware = make([]Middleware, 0)
	bus.processingMiddleware = make([]Middleware, 0)
	bus.postMiddleware = make([]Middleware, 0)
}

// Publish 发布事件，执行中间件
func (bus *EventBus) Publish(event Event) {
	// 先执行预处理中间件
	event = bus.applyPreMiddleware(event)
	if event == nil {
		return // 事件被中间件截断
	}

	// 然后执行处理中间件
	event = bus.applyProcessingMiddleware(event)
	if event == nil {
		return // 事件被中间件截断
	}

	// 最后执行后处理中间件
	event = bus.applyPostMiddleware(event)
	if event == nil {
		return // 事件被中间件截断
	}
}
