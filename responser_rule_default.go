package cryo

// ToMeRule 内置的 提及我 规则，接收到群聊消息时会检查是否有At到当前用户，否则退出消息事件的处理
//
// 如果指定了 removeAt 参数为 false，则不会移除 At 元素
func ToMeRule(removeAt ...bool) Rule[Event] {
	// 如果没有指定是否移除At，默认移除
	if len(removeAt) == 0 {
		removeAt = append(removeAt, true)
	}
	if removeAt[0] {
		return RuleFor(
			func(e *GroupMessageEvent) bool {
				msg := *e.GetMessage()

				// 快速检查：如果消息为空，直接返回
				if len(msg) == 0 {
					return false
				}

				// 查找需要移除的 At 元素的索引
				atIndex := -1
				for i, element := range msg {
					if element.GetType() == AtType {
						if at, ok := element.(*At); ok && at.TargetUin == e.GetClient().Uin {
							atIndex = i
							break
						}
					}
				}

				// 如果没找到 At，直接返回
				if atIndex == -1 {
					return false
				}

				// 处理找到 At 的情况
				// 检查 At 后面的文本是否需要处理空格
				if atIndex+1 < len(msg) && msg[atIndex+1].GetType() == TextType {
					if text, ok := msg[atIndex+1].(*Text); ok && len(text.Content) > 0 && text.Content[0] == ' ' {
						text.Content = text.Content[1:]
					}
				}

				// 构建新消息，移除 At 元素
				newMsg := make(Message, 0, len(msg)-1)
				newMsg = append(newMsg, msg[:atIndex]...)
				newMsg = append(newMsg, msg[atIndex+1:]...)

				// 更新消息并返回
				*e.GetMessage() = newMsg
				return true
			})
	} else {
		return RuleFor(
			func(e *GroupMessageEvent) bool {
				msg := *e.GetMessage()

				// 快速检查：如果消息为空，直接返回
				if len(msg) == 0 {
					return false
				}

				// 遍历所有消息元素，检查是否有At当前用户的元素
				for _, element := range msg {
					if element.GetType() == AtType {
						if at, ok := element.(*At); ok && at.TargetUin == e.GetClient().Uin {
							// 找到At当前用户的元素，直接返回true
							return true
						}
					}
				}

				// 没找到At当前用户的元素
				return false
			})
	}
}

// AtRule 内置的 At 规则，接收到群聊消息时会检查是否有At到指定用户，否则退出消息事件的处理
func AtRule(target ...uint32) Rule[Event] {
	return RuleFor(
		func(e *GroupMessageEvent) bool {
			// 如果没有指定目标用户，直接返回false
			if len(target) == 0 {
				return false
			}

			msg := *e.GetMessage()
			// 快速检查：如果消息为空，直接返回
			if len(msg) == 0 {
				return false
			}

			// 将targets转换为map以提高查找效率
			targetMap := make(map[uint32]struct{}, len(target))
			for _, t := range target {
				targetMap[t] = struct{}{}
			}

			// 遍历所有消息元素，检查是否有At指定用户的元素
			for _, element := range msg {
				if element.GetType() == AtType {
					if at, ok := element.(*At); ok {
						if _, exists := targetMap[at.TargetUin]; exists {
							// 找到At指定用户的元素，直接返回true
							return true
						}
					}
				}
			}

			// 没找到At指定用户的元素
			return false
		})
}

// StartWithRule 内置的前缀匹配规则，检查消息是否恰好只包含一个文本元素，且以指定的前缀开头
//
// 如果消息包含多个元素或非文本元素，直接返回false
func StartWithRule(prefix ...string) Rule[Event] {
	return RuleFor(
		func(e *GroupMessageEvent) bool {
			// 如果没有指定前缀，直接返回false
			if len(prefix) == 0 {
				return false
			}

			msg := *e.GetMessage()
			// 检查消息是否恰好只包含一个元素
			if len(msg) != 1 {
				return false
			}

			// 检查唯一的元素是否为文本类型
			element := msg[0]
			if element.GetType() != TextType {
				return false
			}

			// 类型断言为文本元素
			text, ok := element.(*Text)
			if !ok {
				return false
			}

			content := text.Content

			// 检查前缀
			for _, p := range prefix {
				if len(content) >= len(p) && content[:len(p)] == p {
					return true
				}
			}

			return false
		})
}

// EndWithRule 内置的后缀匹配规则，检查消息是否恰好只包含一个文本元素，且以指定的后缀结尾
//
// 如果消息包含多个元素或非文本元素，直接返回false
func EndWithRule(suffix ...string) Rule[Event] {
	return RuleFor(
		func(e *GroupMessageEvent) bool {
			// 如果没有指定后缀，直接返回false
			if len(suffix) == 0 {
				return false
			}

			msg := *e.GetMessage()
			// 检查消息是否恰好只包含一个元素
			if len(msg) != 1 {
				return false
			}

			// 检查唯一的元素是否为文本类型
			element := msg[0]
			if element.GetType() != TextType {
				return false
			}

			// 类型断言为文本元素
			text, ok := element.(*Text)
			if !ok {
				return false
			}

			content := text.Content

			// 检查后缀
			for _, s := range suffix {
				if len(content) >= len(s) && content[len(content)-len(s):] == s {
					return true
				}
			}

			return false
		})
}

// FullMatchRule 内置的完全匹配规则，检查消息是否恰好只包含一个文本元素，且内容与指定的内容完全相同
//
// 如果消息包含多个元素或非文本元素，直接返回false
func FullMatchRule(content ...string) Rule[Event] {
	return RuleFor(
		func(e *GroupMessageEvent) bool {
			// 如果没有指定内容，直接返回false
			if len(content) == 0 {
				return false
			}

			msg := *e.GetMessage()
			// 检查消息是否恰好只包含一个元素
			if len(msg) != 1 {
				return false
			}

			// 检查唯一的元素是否为文本类型
			element := msg[0]
			if element.GetType() != TextType {
				return false
			}

			// 类型断言为文本元素
			text, ok := element.(*Text)
			if !ok {
				return false
			}

			contentStr := text.Content

			// 检查内容
			for _, c := range content {
				if contentStr == c {
					return true
				}
			}

			return false
		})
}

// KeyWordRule 内置的关键词匹配规则，检查消息是否包含文本元素，如果包含则遍历文本元素在其中查找关键词
func KeyWordRule(keyword ...string) Rule[Event] {
	return RuleFor(
		func(e *GroupMessageEvent) bool {
			// 如果没有指定关键词，直接返回false
			if len(keyword) == 0 {
				return false
			}

			msg := *e.GetMessage()
			// 快速检查：如果消息为空，直接返回
			if len(msg) == 0 {
				return false
			}

			// 将关键词转换为map以提高查找效率
			keywordMap := make(map[string]struct{}, len(keyword))
			for _, k := range keyword {
				keywordMap[k] = struct{}{}
			}

			// 遍历所有消息元素，寻找文本元素
			for _, element := range msg {
				if element.GetType() == TextType {
					if text, ok := element.(*Text); ok {
						content := text.Content

						// 检查每个关键词是否在文本内容中
						for _, k := range keyword {
							if len(k) > 0 && ContainsKeyword(content, k) {
								return true
							}
						}
					}
				}
			}

			return false
		})
}
