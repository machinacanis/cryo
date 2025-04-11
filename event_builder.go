package cryo

import "time"

func SendBotConnectedEvent(c *LagrangeClient) {
	// 发送bot连接事件
	c.bus.Publish(&BotConnectedEvent{
		UniEvent: UniEvent{
			payload:        nil,
			EventType:      BotConnectedEventType,
			EventId:        newUUID(),
			EventTags:      []string{"cryo", "bot_connected"},
			Time:           uint32(time.Now().Unix()),
			botClient:      c,
			ClientId:       c.Id,
			ClientNickname: c.Nickname,
			ClientUin:      c.Uin,
			ClientUid:      c.Uid,
			Platform:       c.Platform,
		},
		Version: c.Version,
	})
}

// SendScheduledTaskRegisteredEvent 发送定时任务注册事件
func SendScheduledTaskRegisteredEvent(b *Bot, task *ScheduledTask) {
	event := &ScheduledTaskRegisteredEvent{
		UniEvent: UniEvent{
			EventType:      ScheduledTaskRegisteredEventType,
			EventId:        newUUID(),
			EventTags:      []string{"cryo", "scheduled_task_registered"},
			Time:           uint32(time.Now().Unix()),
			botClient:      nil,
			ClientId:       "",
			ClientNickname: "",
			ClientUin:      0,
			ClientUid:      "",
			Platform:       "",
		},
		task: task,
	}
	b.bus.Publish(event) // 发布事件
}

// SendScheduledTaskSuccessEvent 发送定时任务成功事件
func SendScheduledTaskSuccessEvent(b *Bot, task *ScheduledTask) {
	event := &ScheduledTaskSuccessEvent{
		UniEvent: UniEvent{
			EventType:      ScheduledTaskSuccessEventType,
			EventId:        newUUID(),
			EventTags:      []string{"cryo", "scheduled_task_success"},
			Time:           uint32(time.Now().Unix()),
			botClient:      nil,
			ClientId:       "",
			ClientNickname: "",
			ClientUin:      0,
			ClientUid:      "",
			Platform:       "",
		},
		task: task,
	}
	b.bus.Publish(event) // 发布事件
}

// SendScheduledTaskFailedEvent 发送定时任务失败事件
func SendScheduledTaskFailedEvent(b *Bot, task *ScheduledTask) {
	event := &ScheduledTaskFailedEvent{
		UniEvent: UniEvent{
			EventType:      ScheduledTaskFailedEventType,
			EventId:        newUUID(),
			EventTags:      []string{"cryo", "scheduled_task_failed"},
			Time:           uint32(time.Now().Unix()),
			botClient:      nil,
			ClientId:       "",
			ClientNickname: "",
			ClientUin:      0,
			ClientUid:      "",
			Platform:       "",
		},
		task: task,
	}
	b.bus.Publish(event) // 发布事件
}

// SendScheduledTaskStoppedEvent 发送定时任务停止事件
func SendScheduledTaskStoppedEvent(b *Bot, task *ScheduledTask) {
	event := &ScheduledTaskStoppedEvent{
		UniEvent: UniEvent{
			EventType:      ScheduledTaskStoppedEventType,
			EventId:        newUUID(),
			EventTags:      []string{"cryo", "scheduled_task_stopped"},
			Time:           uint32(time.Now().Unix()),
			botClient:      nil,
			ClientId:       "",
			ClientNickname: "",
			ClientUin:      0,
			ClientUid:      "",
			Platform:       "",
		},
		task: task,
	}
	b.bus.Publish(event) // 发布事件
}
