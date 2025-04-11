package cryo

import (
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/client/event"
	"github.com/LagrangeDev/LagrangeGo/message"
	"time"
)

// EventBind 绑定LagrangeGo的事件到cryobot的事件总线
func (c *LagrangeClient) eventBind() {

	c.logger.Infof("[Cryo] 正在将 %d 的消息事件绑定到事件总线", c.Client.Uin)
	// 断开连接
	c.Client.DisconnectedEvent.Subscribe(func(client *client.QQClient, event *client.DisconnectedEvent) {
		c.bus.Publish(&BotDisconnectedEvent{UniEvent{
			payload:        nil,
			EventType:      BotDisconnectedEventType,
			EventId:        newUUID(),
			EventTags:      []string{"cryo", "bot_disconnected"},
			Time:           uint32(time.Now().Unix()),
			botClient:      c,
			ClientId:       c.Id,
			ClientNickname: c.Nickname,
			ClientUin:      c.Uin,
			ClientUid:      c.Uid,
			Platform:       c.Platform,
		}})
	})

	// 私聊消息
	c.Client.PrivateMessageEvent.Subscribe(func(client *client.QQClient, event *message.PrivateMessage) {
		m := Message{}
		m.AddIMessageElement(event.Elements...)
		c.bus.Publish(&PrivateMessageEvent{
			UniMessageEvent: UniMessageEvent{
				UniEvent: UniEvent{
					payload:        nil,
					EventType:      PrivateMessageEventType,
					EventId:        newUUID(),
					EventTags:      []string{"message", "private_message"},
					Time:           event.Time,
					botClient:      c,
					ClientId:       c.Id,
					ClientNickname: c.Nickname,
					ClientUin:      c.Uin,
					ClientUid:      c.Uid,
					Platform:       c.Platform,
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
					payload:        nil,
					EventType:      GroupMessageEventType,
					EventId:        newUUID(),
					EventTags:      []string{"message", "group_message"},
					Time:           event.Time,
					botClient:      c,
					ClientId:       c.Id,
					ClientNickname: c.Nickname,
					ClientUin:      c.Uin,
					ClientUid:      c.Uid,
					Platform:       c.Platform,
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

	// 临时消息
	c.Client.TempMessageEvent.Subscribe(func(client *client.QQClient, event *message.TempMessage) {
		m := Message{}
		m.AddIMessageElement(event.Elements...)
		c.bus.Publish(&TempMessageEvent{
			UniMessageEvent: UniMessageEvent{
				UniEvent: UniEvent{
					payload:        nil,
					EventType:      TempMessageEventType,
					EventId:        newUUID(),
					EventTags:      []string{"message", "temp_message"},
					Time:           uint32(time.Now().Unix()),
					botClient:      c,
					ClientId:       c.Id,
					ClientNickname: c.Nickname,
					ClientUin:      c.Uin,
					ClientUid:      c.Uid,
					Platform:       c.Platform,
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

	// 好友请求
	c.Client.NewFriendRequestEvent.Subscribe(func(client *client.QQClient, event *event.NewFriendRequest) {
		c.bus.Publish(&NewFriendRequestEvent{
			UniEvent: UniEvent{
				payload:        nil,
				EventType:      NewFriendRequestEventType,
				EventId:        newUUID(),
				EventTags:      []string{"notice", "new_friend_request"},
				Time:           uint32(time.Now().Unix()),
				botClient:      c,
				ClientId:       c.Id,
				ClientNickname: c.Nickname,
				ClientUin:      c.Uin,
				ClientUid:      c.Uid,
				Platform:       c.Platform,
			},
			Uin:      event.SourceUin,
			Uid:      event.SourceUID,
			Nickname: event.SourceNick,
			Message:  event.Msg,
			From:     event.Source,
		})
	})

	// 新好友
	c.Client.NewFriendEvent.Subscribe(func(client *client.QQClient, event *event.NewFriend) {
		c.bus.Publish(&NewFriendEvent{
			UniEvent: UniEvent{
				payload:        nil,
				EventType:      NewFriendEventType,
				EventId:        newUUID(),
				EventTags:      []string{"notice", "new_friend"},
				Time:           uint32(time.Now().Unix()),
				botClient:      c,
				ClientId:       c.Id,
				ClientNickname: c.Nickname,
				ClientUin:      c.Uin,
				ClientUid:      c.Uid,
				Platform:       c.Platform,
			},
			Uin:      event.FromUin,
			Uid:      event.FromUID,
			Nickname: event.FromNick,
			Message:  event.Msg,
		})
	})

	// 好友撤回
	c.Client.FriendRecallEvent.Subscribe(func(client *client.QQClient, event *event.FriendRecall) {
		c.bus.Publish(&FriendRecallEvent{
			UniEvent: UniEvent{
				payload:        nil,
				EventType:      FriendRecallEventType,
				EventId:        newUUID(),
				EventTags:      []string{"notice", "friend_recall"},
				Time:           uint32(time.Now().Unix()),
				botClient:      c,
				ClientId:       c.Id,
				ClientNickname: c.Nickname,
				ClientUin:      c.Uin,
				ClientUid:      c.Uid,
				Platform:       c.Platform,
			},
			Uin:     event.FromUin,
			Uid:     event.FromUID,
			Seqence: event.Sequence,
			Random:  event.Random,
		})
	})

	// 改名
	c.Client.RenameEvent.Subscribe(func(client *client.QQClient, event *event.Rename) {
		c.bus.Publish(&FriendRenameEvent{
			UniEvent: UniEvent{
				payload:        nil,
				EventType:      FriendRenameEventType,
				EventId:        newUUID(),
				EventTags:      []string{"notice", "friend_rename"},
				Time:           uint32(time.Now().Unix()),
				botClient:      c,
				ClientId:       c.Id,
				ClientNickname: c.Nickname,
				ClientUin:      c.Uin,
				ClientUid:      c.Uid,
				Platform:       c.Platform,
			},
			IsSelf:   event.SubType == 0,
			Uin:      event.Uin,
			Uid:      event.UID,
			Nickname: event.Nickname,
		})
	})

	// 好友戳一戳 暂时不支持

	// 群成员权限变动
	c.Client.GroupMemberPermissionChangedEvent.Subscribe(func(client *client.QQClient, event *event.GroupMemberPermissionChanged) {
		c.bus.Publish(&GroupMemberPermissionUpdatedEvent{
			UniEvent: UniEvent{
				payload:        nil,
				EventType:      GroupMemberPermissionUpdatedEventType,
				EventId:        newUUID(),
				EventTags:      []string{"notice", "group_member_permission_updated"},
				Time:           uint32(time.Now().Unix()),
				botClient:      c,
				ClientId:       c.Id,
				ClientNickname: c.Nickname,
				ClientUin:      c.Uin,
				ClientUid:      c.Uid,
				Platform:       c.Platform,
			},
			GroupUin: event.GroupUin,
			Uin:      event.UserUin,
			Uid:      event.UserUID,
			IsAdmin:  event.IsAdmin,
		})
	})

	// 群改名
	c.Client.GroupNameUpdatedEvent.Subscribe(func(client *client.QQClient, event *event.GroupNameUpdated) {
		c.bus.Publish(&GroupNameUpdatedEvent{
			UniEvent: UniEvent{
				payload:        nil,
				EventType:      GroupNameUpdatedEventType,
				EventId:        newUUID(),
				EventTags:      []string{"notice", "group_name_updated"},
				Time:           uint32(time.Now().Unix()),
				botClient:      c,
				ClientId:       c.Id,
				ClientNickname: c.Nickname,
				ClientUin:      c.Uin,
				ClientUid:      c.Uid,
				Platform:       c.Platform,
			},
			GroupUin: event.GroupUin,
			NewName:  event.NewName,
		})
	})

	// 群禁言
	c.Client.GroupMuteEvent.Subscribe(func(client *client.QQClient, event *event.GroupMute) {
		c.bus.Publish(&GroupMuteEvent{
			UniEvent: UniEvent{
				payload:        nil,
				EventType:      GroupMuteEventType,
				EventId:        newUUID(),
				EventTags:      []string{"notice", "group_mute"},
				Time:           uint32(time.Now().Unix()),
				botClient:      c,
				ClientId:       c.Id,
				ClientNickname: c.Nickname,
				ClientUin:      c.Uin,
				ClientUid:      c.Uid,
				Platform:       c.Platform,
			},
			GroupUin:    event.GroupUin,
			OperatorUin: event.OperatorUin,
			OperatorUid: event.OperatorUID,
			TargetUin:   event.UserUin,
			TargetUid:   event.UserUID,
			Duration:    event.Duration,
			isMuteAll:   event.MuteAll(),
		})
	})

	// 群撤回
	c.Client.GroupRecallEvent.Subscribe(func(client *client.QQClient, event *event.GroupRecall) {
		c.bus.Publish(&GroupRecallEvent{
			UniEvent: UniEvent{
				payload:        nil,
				EventType:      GroupRecallEventType,
				EventId:        newUUID(),
				EventTags:      []string{"notice", "group_recall"},
				Time:           uint32(time.Now().Unix()),
				botClient:      c,
				ClientId:       c.Id,
				ClientNickname: c.Nickname,
				ClientUin:      c.Uin,
				ClientUid:      c.Uid,
				Platform:       c.Platform,
			},
			GroupUin:    event.GroupUin,
			OperatorUin: event.OperatorUin,
			OperatorUid: event.OperatorUID,
			SenderUin:   event.UserUin,
			SenderUid:   event.UserUID,
			Random:      event.Random,
			Seqence:     event.Sequence,
		})
	})

	// 群成员入群请求
	c.Client.GroupMemberJoinRequestEvent.Subscribe(func(client *client.QQClient, event *event.GroupMemberJoinRequest) {
		c.bus.Publish(&GroupMemberJoinRequestEvent{
			UniEvent: UniEvent{
				payload:        nil,
				EventType:      GroupMemberJoinRequestEventType,
				EventId:        newUUID(),
				EventTags:      []string{"notice", "group_member_join_request"},
				Time:           uint32(time.Now().Unix()),
				botClient:      c,
				ClientId:       c.Id,
				ClientNickname: c.Nickname,
				ClientUin:      c.Uin,
				ClientUid:      c.Uid,
				Platform:       c.Platform,
			},
			GroupUin:       event.GroupUin,
			SenderUin:      event.UserUin,
			SenderUid:      event.UserUID,
			SenderNickname: event.TargetNick,
			InviterUin:     event.UserUin,
			InviterUid:     event.UserUID,
			Answer:         event.Answer,
			RequestSeqence: event.RequestSeq,
		})
	})

	// 群成员增加
	c.Client.GroupMemberJoinEvent.Subscribe(func(client *client.QQClient, event *event.GroupMemberIncrease) {
		c.bus.Publish(&GroupMemberIncreaseEvent{
			UniEvent: UniEvent{
				payload:        nil,
				EventType:      GroupMemberIncreaseEventType,
				EventId:        newUUID(),
				EventTags:      []string{"notice", "group_member_increase"},
				Time:           uint32(time.Now().Unix()),
				botClient:      c,
				ClientId:       c.Id,
				ClientNickname: c.Nickname,
				ClientUin:      c.Uin,
				ClientUid:      c.Uid,
				Platform:       c.Platform,
			},
			GroupUin:   event.GroupUin,
			Uin:        event.UserUin,
			Uid:        event.UserUID,
			InviterUin: event.InvitorUin,
			InviterUid: event.InvitorUID,
			IsSelf:     event.UserUin == c.Uin,
		})
	})

	// 群成员减少
	c.Client.GroupMemberLeaveEvent.Subscribe(func(client *client.QQClient, event *event.GroupMemberDecrease) {
		c.bus.Publish(&GroupMemberDecreaseEvent{
			UniEvent: UniEvent{
				payload:        nil,
				EventType:      GroupMemberDecreaseEventType,
				EventId:        newUUID(),
				EventTags:      []string{"notice", "group_member_decrease"},
				Time:           uint32(time.Now().Unix()),
				botClient:      c,
				ClientId:       c.Id,
				ClientNickname: c.Nickname,
				ClientUin:      c.Uin,
				ClientUid:      c.Uid,
				Platform:       c.Platform,
			},
			GroupUin: event.GroupUin,
			Uin:      event.UserUin,
			Uid:      event.UserUID,
			IsSelf:   event.UserUin == c.Uin,
			IsKicked: event.IsKicked(),
		})
	})

	// 群精华消息
	c.Client.GroupDigestEvent.Subscribe(func(client *client.QQClient, event *event.GroupDigestEvent) {
		c.bus.Publish(&GroupDigestEvent{
			UniEvent: UniEvent{
				payload:        nil,
				EventType:      GroupDigestEventType,
				EventId:        newUUID(),
				EventTags:      []string{"notice", "group_digest"},
				Time:           uint32(time.Now().Unix()),
				botClient:      c,
				ClientId:       c.Id,
				ClientNickname: c.Nickname,
				ClientUin:      c.Uin,
				ClientUid:      c.Uid,
				Platform:       c.Platform,
			},
			GroupUin:         event.GroupUin,
			MessageId:        event.MessageID,
			InternalId:       event.InternalMessageID,
			SenderUin:        event.UserUin,
			SenderUid:        event.UserUID,
			SenderNickname:   event.SenderNick,
			OperatorUin:      event.OperatorUin,
			OperatorNickname: event.OperatorNick,
			IsRemove:         event.OperationType == 2,
		})
	})

	// 群表态事件
	c.Client.GroupReactionEvent.Subscribe(func(client *client.QQClient, event *event.GroupReactionEvent) {
		c.bus.Publish(&GroupReactionEvent{
			UniEvent: UniEvent{
				payload:        nil,
				EventType:      GroupReactionEventType,
				EventId:        newUUID(),
				EventTags:      []string{"notice", "group_reaction"},
				Time:           uint32(time.Now().Unix()),
				botClient:      c,
				ClientId:       c.Id,
				ClientNickname: c.Nickname,
				ClientUin:      c.Uin,
				ClientUid:      c.Uid,
				Platform:       c.Platform,
			},
			GroupUin:  event.GroupUin,
			Uin:       event.UserUin,
			Uid:       event.UserUID,
			TargetSeq: event.TargetSeq,
			IsAdd:     event.IsAdd,
			IsEmoji:   event.IsEmoji,
			Code:      event.Code,
			Count:     event.Count,
		})
	})

	// 群成员头衔变更
	c.Client.MemberSpecialTitleUpdatedEvent.Subscribe(func(client *client.QQClient, event *event.MemberSpecialTitleUpdated) {
		c.bus.Publish(&GroupMemberSpecialTitleUpdated{
			UniEvent: UniEvent{
				payload:        nil,
				EventType:      GroupMemberSpecialTitleUpdatedEventType,
				EventId:        newUUID(),
				EventTags:      []string{"notice", "group_member_special_title_updated"},
				Time:           uint32(time.Now().Unix()),
				botClient:      c,
				ClientId:       c.Id,
				ClientNickname: c.Nickname,
				ClientUin:      c.Uin,
				ClientUid:      c.Uid,
				Platform:       c.Platform,
			},
			GroupUin: event.GroupUin,
			Uin:      event.UserUin,
			Uid:      event.UserUID,
			NewTitle: event.NewTitle,
		})
	})

	// 群邀请
	c.Client.GroupInvitedEvent.Subscribe(func(client *client.QQClient, event *event.GroupInvite) {
		c.bus.Publish(&GroupInviteEvent{
			UniEvent: UniEvent{
				payload:        nil,
				EventType:      GroupInviteEventType,
				EventId:        newUUID(),
				EventTags:      []string{"notice", "group_invite"},
				Time:           uint32(time.Now().Unix()),
				botClient:      c,
				ClientId:       c.Id,
				ClientNickname: c.Nickname,
				ClientUin:      c.Uin,
				ClientUid:      c.Uid,
				Platform:       c.Platform,
			},
			GroupUin:        event.GroupUin,
			GroupName:       event.GroupName,
			InviterUin:      event.InvitorUin,
			InviterUid:      event.InvitorUID,
			InviterNickname: event.InvitorNick,
			RequestSeqence:  event.RequestSeq,
		})
	})

	c.logger.Successf("[Cryo] %d 的消息事件绑定完成", c.Client.Uin)
}
