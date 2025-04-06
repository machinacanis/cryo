package cryo

// EventHandler 是事件处理器的函数类型
type EventHandler[T Event] func(event T) T

// Middleware 用来定义事件处理过程中的中间件行为
type Middleware interface {
	IsGlobal() bool                             // 是否是全局中间件
	GetId() string                              // 获取中间件ID
	GetTag() []string                           // 获取中间件标签
	AddHandler(handlers ...EventHandler[Event]) // 添加事件处理器
	GetHandler() []EventHandler[Event]          // 获取事件处理器列表
	Do(event Event) bool                        // 执行中间件
	DoAsync(event Event)                        // 并发执行中间件，谨慎使用
}

// UniMiddleware 是一个中间件的实现，包含了事件处理器和中间件的基本信息
type UniMiddleware struct {
	receivingEventType []EventType           // 这个中间件接收的事件类型，如果是空，则接收所有事件（即全局中间件）
	id                 string                // 中间件ID
	tag                []string              // 中间件标签
	handlers           []EventHandler[Event] // 事件处理器列表
}

// IsGlobal 判断这个中间件是否是全局中间件
func (m *UniMiddleware) IsGlobal() bool {
	return len(m.receivingEventType) == 0
}

// GetId 获取中间件ID
func (m *UniMiddleware) GetId() string {
	return m.id
}

// GetTag 获取中间件标签
func (m *UniMiddleware) GetTag() []string {
	return m.tag
}

// AddHandler 添加事件处理器
func (m *UniMiddleware) AddHandler(handlers ...EventHandler[Event]) {
	m.handlers = append(m.handlers, handlers...)
}

// GetHandler 获取事件处理器列表
func (m *UniMiddleware) GetHandler() []EventHandler[Event] {
	return m.handlers
}

// Do 执行中间件
func (m *UniMiddleware) Do(event Event) bool {
	for _, handler := range m.handlers {
		event = handler(event)
		if event == nil {
			return false // 事件被中间件截断
		}
	}
	return true
}

// DoAsync 并发执行中间件，谨慎使用
func (m *UniMiddleware) DoAsync(event Event) {
	// 使用 goroutine 并发执行中间件
	go func() {
		for _, handler := range m.handlers {
			event = handler(event)
			if event == nil {
				return // 事件被中间件截断
			}
		}
	}
}
