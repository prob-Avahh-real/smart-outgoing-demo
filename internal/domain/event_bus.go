package domain

// EventBus 事件总线接口
type EventBus interface {
	Publish(event DomainEvent) error
	Subscribe(eventType string, handler EventHandler) error
	Unsubscribe(eventType string, handler EventHandler) error
}

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(event DomainEvent) error
}

// InMemoryEventBus 内存事件总线实现
type InMemoryEventBus struct {
	handlers map[string][]EventHandler
}

// NewInMemoryEventBus 创建内存事件总线
func NewInMemoryEventBus() EventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Publish 发布事件
func (bus *InMemoryEventBus) Publish(event DomainEvent) error {
	if event == nil {
		return nil
	}

	handlers, exists := bus.handlers[event.EventType()]
	if !exists {
		return nil // 没有处理器，直接返回
	}

	// 异步处理所有事件处理器
	for _, handler := range handlers {
		go func(h EventHandler) {
			if err := h.Handle(event); err != nil {
				// 记录错误但不中断其他处理器
				// 在实际实现中应该有更好的错误处理
			}
		}(handler)
	}

	return nil
}

// Subscribe 订阅事件
func (bus *InMemoryEventBus) Subscribe(eventType string, handler EventHandler) error {
	if eventType == "" || handler == nil {
		return nil
	}

	bus.handlers[eventType] = append(bus.handlers[eventType], handler)
	return nil
}

// Unsubscribe 取消订阅事件
func (bus *InMemoryEventBus) Unsubscribe(eventType string, handler EventHandler) error {
	handlers, exists := bus.handlers[eventType]
	if !exists {
		return nil
	}

	// 移除指定的处理器
	for i, h := range handlers {
		// 简化比较，实际实现中可能需要更复杂的比较逻辑
		if h == handler {
			bus.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}

	return nil
}
