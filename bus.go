package cryo

import (
	"sync"
)

// MiddlewareOrdering 是中间件的执行顺序类型别名
type MiddlewareOrdering int

const (
	PreMiddlewareType   MiddlewareOrdering = iota // 预处理中间件
	PostMiddlewareType                            // 后处理中间件
	SyncMiddlewareType                            // 同步处理中间件
	AsyncMiddlewareType                           // 异步处理中间件
)

// EventBus 是事件总线，负责管理中间件和事件的分发
type EventBus struct {
	middlewareMutex sync.RWMutex // 保护中间件的读写锁
	preMiddleware   []Middleware // 预处理中间件列表
	postMiddleware  []Middleware // 后处理中间件列表
	syncMiddleware  []Middleware // 中间件列表
	asyncMiddleware []Middleware // 并发中间件列表
}

// NewEventBus 创建一个新的事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		preMiddleware:   make([]Middleware, 0),
		syncMiddleware:  make([]Middleware, 0),
		postMiddleware:  make([]Middleware, 0),
		asyncMiddleware: make([]Middleware, 0),
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

func (bus *EventBus) applySyncMiddleware(event Event) Event {
	if len(bus.syncMiddleware) == 0 {
		return event // 如果没有处理中间件，则直接返回事件
	}
	eventType := event.GetEventType() // 获取事件类型

	bus.middlewareMutex.RLock() // 持有读锁
	// 创建一个中间件切片的副本，减小锁的粒度
	var middlewareCopy []Middleware
	if len(bus.syncMiddleware) > 0 {
		middlewareCopy = make([]Middleware, len(bus.syncMiddleware))
		copy(middlewareCopy, bus.syncMiddleware)
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
		eventCopy := event.Clone()      // 事件的副本

		if middleware.IsGlobal() || middleware.HasType(eventType) {
			wg.Add(1)
			// 为每个中间件创建一个 goroutine
			go func(m Middleware) {
				defer wg.Done()
				m.Do(eventCopy)
			}(middleware) // 传递中间件实例作为参数
		}
	}

	// 等待所有中间件处理完成
	// wg.Wait()

	// 处理中间件是用来替代事件订阅的，它不会对事件进行修改，只是用来把事件分发到对应的处理器
	// 自然也不需要等待所有中间件处理完成，直接返回原来的事件即可
	return event // 返回事件
}

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

func (bus *EventBus) applyAsyncMiddleware(event Event) Event {
	if len(bus.asyncMiddleware) == 0 {
		return event // 如果没有并发处理中间件，则直接返回事件
	}
	eventType := event.GetEventType() // 获取事件类型

	bus.middlewareMutex.RLock() // 持有读锁
	// 创建一个中间件切片的副本，减小锁的粒度
	var middlewareCopy []Middleware
	if len(bus.asyncMiddleware) > 0 {
		middlewareCopy = make([]Middleware, len(bus.asyncMiddleware))
		copy(middlewareCopy, bus.asyncMiddleware)
	}
	bus.middlewareMutex.RUnlock() // 释放读锁

	// 对于并发中间件，它的逻辑是完全无序的，只是传入一个事件然后让它们全部并发执行
	var wg sync.WaitGroup

	for i := range middlewareCopy {
		middleware := middlewareCopy[i] // 在循环内部创建局部变量，避免闭包陷阱
		eventCopy := event.Clone()

		if middleware.IsGlobal() || middleware.HasType(eventType) {
			wg.Add(1)
			// 为每个中间件创建一个 goroutine
			go func(m Middleware) {
				defer wg.Done()
				m.DoAsync(eventCopy)
			}(middleware) // 传递中间件实例作为参数
		}
	}
	return event // 返回源事件
}

// AddPreMiddleware 添加预处理中间件
func (bus *EventBus) AddPreMiddleware(middleware ...Middleware) {
	bus.middlewareMutex.Lock() // 持有写锁
	defer bus.middlewareMutex.Unlock()
	for _, m := range middleware {
		bus.preMiddleware = append(bus.preMiddleware, m)
	}
}

// AddSyncMiddleware 添加同步中间件
func (bus *EventBus) AddSyncMiddleware(middleware ...Middleware) {
	bus.middlewareMutex.Lock() // 持有写锁
	defer bus.middlewareMutex.Unlock()
	for _, m := range middleware {
		bus.syncMiddleware = append(bus.syncMiddleware, m)
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

// AddAsyncMiddleware 添加异步中间件
func (bus *EventBus) AddAsyncMiddleware(middleware ...Middleware) {
	bus.middlewareMutex.Lock() // 持有写锁
	defer bus.middlewareMutex.Unlock()
	for _, m := range middleware {
		bus.asyncMiddleware = append(bus.asyncMiddleware, m)
	}
}

// RemoveMiddlewareById 删除指定Id的中间件
func (bus *EventBus) RemoveMiddlewareById(id ...string) bool {
	if len(id) == 0 {
		return false
	}

	bus.middlewareMutex.Lock() // 持有写锁
	defer bus.middlewareMutex.Unlock()

	// 创建一个map来存储要删除的ID，提高查找效率
	idMap := make(map[string]struct{}, len(id))
	for _, i := range id {
		if i != "" {
			idMap[i] = struct{}{}
		}
	}

	if len(idMap) == 0 {
		return false
	}

	// 创建一个新的切片来存储不包含指定ID的中间件
	newPreMiddleware := make([]Middleware, 0, len(bus.preMiddleware))
	for _, m := range bus.preMiddleware {
		if _, exists := idMap[m.GetId()]; !exists {
			newPreMiddleware = append(newPreMiddleware, m)
		}
	}
	bus.preMiddleware = newPreMiddleware

	newSyncMiddleware := make([]Middleware, 0, len(bus.syncMiddleware))
	for _, m := range bus.syncMiddleware {
		if _, exists := idMap[m.GetId()]; !exists {
			newSyncMiddleware = append(newSyncMiddleware, m)
		}
	}
	bus.syncMiddleware = newSyncMiddleware

	newPostMiddleware := make([]Middleware, 0, len(bus.postMiddleware))
	for _, m := range bus.postMiddleware {
		if _, exists := idMap[m.GetId()]; !exists {
			newPostMiddleware = append(newPostMiddleware, m)
		}
	}
	bus.postMiddleware = newPostMiddleware

	newAsyncMiddleware := make([]Middleware, 0, len(bus.asyncMiddleware))
	for _, m := range bus.asyncMiddleware {
		if _, exists := idMap[m.GetId()]; !exists {
			newAsyncMiddleware = append(newAsyncMiddleware, m)
		}
	}
	bus.asyncMiddleware = newAsyncMiddleware

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

	newSyncMiddleware := make([]Middleware, 0, len(bus.syncMiddleware))
	for _, m := range bus.syncMiddleware {
		if !m.HasTag(tag...) {
			newSyncMiddleware = append(newSyncMiddleware, m)
		}
	}
	bus.syncMiddleware = newSyncMiddleware

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

// ClearSyncMiddleware 清空同步中间件
func (bus *EventBus) ClearSyncMiddleware() {
	bus.middlewareMutex.Lock()
	defer bus.middlewareMutex.Unlock()
	bus.syncMiddleware = make([]Middleware, 0)
}

// ClearPostMiddleware 清空后处理中间件
func (bus *EventBus) ClearPostMiddleware() {
	bus.middlewareMutex.Lock()
	defer bus.middlewareMutex.Unlock()
	bus.postMiddleware = make([]Middleware, 0)
}

// ClearAsyncMiddleware 清空异步中间件
func (bus *EventBus) ClearAsyncMiddleware() {
	bus.middlewareMutex.Lock()
	defer bus.middlewareMutex.Unlock()
	bus.asyncMiddleware = make([]Middleware, 0)
}

// ClearAllMiddleware 清空所有中间件
func (bus *EventBus) ClearAllMiddleware() {
	bus.middlewareMutex.Lock()
	defer bus.middlewareMutex.Unlock()
	bus.preMiddleware = make([]Middleware, 0)
	bus.syncMiddleware = make([]Middleware, 0)
	bus.postMiddleware = make([]Middleware, 0)
}

// Publish 发布事件并按顺序执行中间件
func (bus *EventBus) Publish(event Event) {
	// 先执行预处理中间件
	event = bus.applyPreMiddleware(event)
	if event == nil {
		return // 事件被中间件截断
	}

	// 然后执行处理中间件
	bus.applyAsyncMiddleware(event)
	bus.applySyncMiddleware(event)

	// 最后执行后处理中间件
	event = bus.applyPostMiddleware(event)
	if event == nil {
		return // 事件被中间件截断
	}
}
