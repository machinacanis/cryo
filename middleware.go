package cryo

// EventHandler 是事件处理器的函数类型
type EventHandler[T Event] func(event T) T

// Middleware 用来定义事件处理过程中的中间件行为
type Middleware interface {
	IsGlobal() bool                                   // 是否是全局中间件
	GetId() string                                    // 获取中间件ID
	AddTag(tag ...string) *UniMiddleware              // 添加中间件标签
	GetTag() []string                                 // 获取中间件标签
	HasTag(tag ...string) bool                        // 判断中间件是否包含有所有指定标签
	RemoveTag(tag ...string) *UniMiddleware           // 删除中间件标签
	AddType(eventType ...EventType) *UniMiddleware    // 添加中间件接收的事件类型
	GetType() []EventType                             // 获取中间件接收的事件类型
	HasType(eventType ...EventType) bool              // 判断中间件是否接收指定类型的事件
	RemoveType(eventType ...EventType) *UniMiddleware // 删除中间件接收的事件类型
	AddHandler(handlers ...EventHandler[Event])       // 添加事件处理器
	Do(event Event) bool                              // 执行中间件
	DoAsync(event Event)                              // 并发执行中间件，会破坏中间件的顺序性，谨慎使用
}

// UniMiddleware 是一个中间件的实现，包含了事件处理器和中间件的基本信息
type UniMiddleware struct {
	receivingType []EventType           // 这个中间件接收的事件类型，如果是空，则接收所有事件（即全局中间件）
	id            string                // 中间件ID
	tag           []string              // 中间件标签
	Handlers      []EventHandler[Event] // 事件处理器列表
}

// IsGlobal 判断这个中间件是否是全局中间件
func (m *UniMiddleware) IsGlobal() bool {
	return len(m.receivingType) == 0
}

// GetId 获取中间件ID
func (m *UniMiddleware) GetId() string {
	return m.id
}

// AddTag 添加中间件标签
func (m *UniMiddleware) AddTag(tag ...string) *UniMiddleware {
	m.tag = append(m.tag, tag...)
	return m
}

// GetTag 获取中间件标签
func (m *UniMiddleware) GetTag() []string {
	return m.tag
}

// HasTag 判断中间件是否有所有指定标签
func (m *UniMiddleware) HasTag(tag ...string) bool {
	for _, t := range tag {
		for _, mt := range m.tag {
			if t == mt {
				return true // 中间件有指定标签
			}
		}
	}
	return false // 中间件没有指定标签
}

// RemoveTag 删除中间件标签
func (m *UniMiddleware) RemoveTag(tag ...string) *UniMiddleware {
	for _, t := range tag {
		for i, mt := range m.tag {
			if t == mt {
				m.tag = append(m.tag[:i], m.tag[i+1:]...) // 删除标签
				break
			}
		}
	}
	return m
}

// AddType 添加中间件接收的事件类型
func (m *UniMiddleware) AddType(eventType ...EventType) *UniMiddleware {
	m.receivingType = append(m.receivingType, eventType...)
	return m
}

// GetType 获取中间件接收的事件类型
func (m *UniMiddleware) GetType() []EventType {
	return m.receivingType
}

// HasType 判断中间件是否接收指定类型的事件
func (m *UniMiddleware) HasType(eventType ...EventType) bool {
	for _, et := range eventType {
		for _, mt := range m.receivingType {
			if et == mt {
				return true // 中间件接收指定类型的事件
			}
		}
	}
	return false // 中间件不接收指定类型的事件
}

// RemoveType 删除中间件接收的事件类型
func (m *UniMiddleware) RemoveType(eventType ...EventType) *UniMiddleware {
	for _, et := range eventType {
		for i, mt := range m.receivingType {
			if et == mt {
				m.receivingType = append(m.receivingType[:i], m.receivingType[i+1:]...) // 删除事件类型
				break
			}
		}
	}
	return m
}

// AddHandler 添加事件处理器
func (m *UniMiddleware) AddHandler(handlers ...EventHandler[Event]) {
	m.Handlers = append(m.Handlers, handlers...)
}

// Do 执行中间件，返回值为false表示事件被中间件截断
func (m *UniMiddleware) Do(event Event) bool {
	for _, handler := range m.Handlers {
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
		for _, handler := range m.Handlers {
			handler(event)
		}
	}()
}

// NewUniMiddleware 创建一个新的中间件实例
func NewUniMiddleware(eventType ...EventType) *UniMiddleware {
	return &UniMiddleware{
		id:            newUUID(),
		receivingType: eventType,
		tag:           make([]string, 0),
		Handlers:      make([]EventHandler[Event], 0),
	}
}
