package cryo

import (
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/message"
	"github.com/machinacanis/cryo/log"
	"time"
)

// EventBind 绑定LagrangeGo的事件到cryobot的事件总线
func EventBind(c *LagrangeClient) {

	log.Infof("[Cryo] 正在将 %d 的消息事件绑定到事件总线", c.Client.Uin)
	// 断开连接
	c.Client.DisconnectedEvent.Subscribe(func(client *client.QQClient, event *client.DisconnectedEvent) {
		c.bus.Publish(&BotDisconnectedEvent{UniEvent{
			payload:     nil,
			EventType:   BotDisconnectedEventType,
			EventId:     newUUID(),
			EventTags:   []string{"cryo", "bot_disconnected"},
			Time:        uint32(time.Now().Unix()),
			botClient:   c,
			BotId:       c.Id,
			BotNickname: c.Nickname,
			BotUin:      c.Uin,
			BotUid:      c.Uid,
			Platform:    c.Platform,
		}})
	})

	// 私聊消息
	c.Client.PrivateMessageEvent.Subscribe(func(client *client.QQClient, event *message.PrivateMessage) {
		m := Message{}
		m.AddIMessageElement(event.Elements...)
		c.bus.Publish(&PrivateMessageEvent{
			UniMessageEvent: UniMessageEvent{
				UniEvent: UniEvent{
					payload:     nil,
					EventType:   PrivateMessageEventType,
					EventId:     newUUID(),
					EventTags:   []string{"message", "private_message"},
					Time:        event.Time,
					botClient:   c,
					BotId:       c.Id,
					BotNickname: c.Nickname,
					BotUin:      c.Uin,
					BotUid:      c.Uid,
					Platform:    c.Platform,
				},
				MessageId:       event.ID,
				SenderUin:       event.Sender.Uin,
				SenderUid:       event.Sender.UID,
				SenderNickname:  event.Sender.Nickname,
				SenderCardname:  event.Sender.CardName,
				IsSenderFriend:  event.Sender.IsFriend,
				MessageElements: m,
				GroupUin:        event.Sender.Uin,
				GroupName:       event.Sender.Nickname,
			},
			InternalId: event.InternalID,
			ClientSeq:  event.ClientSeq,
			TargetUin:  event.Target,
		})
	})

	// 群聊消息
	c.Client.GroupMessageEvent.Subscribe(func(client *client.QQClient, event *message.GroupMessage) {
		m := Message{}
		m.AddIMessageElement(event.Elements...)
		c.bus.Publish(&GroupMessageEvent{
			UniMessageEvent: UniMessageEvent{
				UniEvent: UniEvent{
					payload:     nil,
					EventType:   GroupMessageEventType,
					EventId:     newUUID(),
					EventTags:   []string{"message", "group_message"},
					Time:        event.Time,
					botClient:   c,
					BotId:       c.Id,
					BotNickname: c.Nickname,
					BotUin:      c.Uin,
					BotUid:      c.Uid,
					Platform:    c.Platform,
				},
				MessageId:       event.ID,
				SenderUin:       event.Sender.Uin,
				SenderUid:       event.Sender.UID,
				SenderNickname:  event.Sender.Nickname,
				SenderCardname:  event.Sender.CardName,
				IsSenderFriend:  event.Sender.IsFriend,
				MessageElements: m,
				GroupUin:        event.GroupUin,
				GroupName:       event.GroupName,
			},
			InternalId: event.InternalID,
		})
	})

	c.Client.TempMessageEvent.Subscribe(func(client *client.QQClient, event *message.TempMessage) {
		m := Message{}
		m.AddIMessageElement(event.Elements...)
		c.bus.Publish(&TempMessageEvent{
			UniMessageEvent: UniMessageEvent{
				UniEvent: UniEvent{
					payload:     nil,
					EventType:   TempMessageEventType,
					EventId:     newUUID(),
					EventTags:   []string{"message", "temp_message"},
					Time:        uint32(time.Now().Unix()),
					botClient:   c,
					BotId:       c.Id,
					BotNickname: c.Nickname,
					BotUin:      c.Uin,
					BotUid:      c.Uid,
					Platform:    c.Platform,
				},
				MessageId:       event.ID,
				SenderUin:       event.Sender.Uin,
				SenderUid:       event.Sender.UID,
				SenderNickname:  event.Sender.Nickname,
				SenderCardname:  event.Sender.CardName,
				IsSenderFriend:  event.Sender.IsFriend,
				MessageElements: m,
				GroupUin:        event.GroupUin,
				GroupName:       event.GroupName,
			},
		})
	})

	log.Infof("[Cryo] %d 的消息事件绑定完成", c.Client.Uin)
}
