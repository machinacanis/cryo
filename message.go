package cryo

import (
	lgrmessage "github.com/LagrangeDev/LagrangeGo/message"
	"github.com/machinacanis/cryo/log"
	"io"
)

// Message 消息元素切片
type Message []MessageElement

// ToIMessageElements 将消息元素切片转换为LagrangeGo的消息元素切片
func (m *Message) ToIMessageElements() []lgrmessage.IMessageElement {
	var elements []lgrmessage.IMessageElement
	for _, element := range *m {
		elements = append(elements, element.GetIMessageElement())
	}
	return elements
}

// ToString 将消息元素切片转换为字符串
func (m *Message) ToString() string {
	var result string
	for _, element := range *m {
		result += element.ToString()
	}
	return result
}

// HasType 是否包含指定类型的消息元素
func (m *Message) HasType(ets ...ElementType) bool {
	for _, element := range *m {
		for _, et := range ets {
			if element.GetType() == et {
				return true
			}
		}
	}
	return false
}

// Add 添加消息元素
func (m *Message) Add(e ...MessageElement) *Message {
	*m = append(*m, e...)
	return m
}

// AddIMessageElement 添加LagrangeGo的消息元素
func (m *Message) AddIMessageElement(e ...lgrmessage.IMessageElement) *Message {
	for _, element := range e {
		switch element.(type) {
		case *lgrmessage.TextElement:
			*m = append(*m, &Text{*element.(*lgrmessage.TextElement)})
		case *lgrmessage.AtElement:
			*m = append(*m, &At{*element.(*lgrmessage.AtElement)})
		case *lgrmessage.FaceElement:
			*m = append(*m, &Face{*element.(*lgrmessage.FaceElement)})
		case *lgrmessage.ReplyElement:
			*m = append(*m, &Reply{*element.(*lgrmessage.ReplyElement)})
		case *lgrmessage.VoiceElement:
			*m = append(*m, &Voice{*element.(*lgrmessage.VoiceElement)})
		case *lgrmessage.ImageElement:
			*m = append(*m, &Image{*element.(*lgrmessage.ImageElement)})
		case *lgrmessage.FileElement:
			*m = append(*m, &File{*element.(*lgrmessage.FileElement)})
		case *lgrmessage.ShortVideoElement:
			*m = append(*m, &ShortVideo{*element.(*lgrmessage.ShortVideoElement)})
		case *lgrmessage.LightAppElement:
			*m = append(*m, &LightApp{*element.(*lgrmessage.LightAppElement)})
		case *lgrmessage.XMLElement:
			*m = append(*m, &XML{*element.(*lgrmessage.XMLElement)})
		case *lgrmessage.ForwardMessage:
			*m = append(*m, &ForwardMessage{*element.(*lgrmessage.ForwardMessage)})
		case *lgrmessage.MarketFaceElement:
			*m = append(*m, &MarketFace{*element.(*lgrmessage.MarketFaceElement)})
		default:
			continue
		}
	}
	return m
}

// AddText 添加文本消息元素
func (m *Message) AddText(text ...string) *Message {
	for _, t := range text {
		*m = append(*m, &Text{*lgrmessage.NewText(t)})
	}
	return m
}

// AddAt 添加提及消息元素
func (m *Message) AddAt(uin uint32, display ...string) *Message {
	*m = append(*m, &At{*lgrmessage.NewAt(uin, display...)})
	return m
}

// AddFace 添加表情消息元素
func (m *Message) AddFace(id uint32) *Message {
	*m = append(*m, &Face{*lgrmessage.NewFace(id)})
	return m
}

// AddReply 添加回复消息元素，需要传入一个MessageEvent作为目标
func (m *Message) AddReply(msg MessageEvent) *Message {
	replySeq, senderUin, time, elements := msg.GetReplyDetail()
	*m = append(*m, &Reply{lgrmessage.ReplyElement{
		ReplySeq:  replySeq,
		SenderUin: senderUin,
		Time:      time,
		Elements:  elements,
	}})
	return m
}

// AddReplyDirect 添加回复消息元素
func (m *Message) AddReplyDirect(replySeq uint32, senderUin uint32, time uint32, elements []lgrmessage.IMessageElement) *Message {
	*m = append(*m, &Reply{lgrmessage.ReplyElement{
		ReplySeq:  replySeq,
		SenderUin: senderUin,
		Time:      time,
		Elements:  elements,
	}})
	return m
}

// AddImage 添加图片消息元素
func (m *Message) AddImage(imgData []byte, summary ...string) *Message {
	*m = append(*m, &Image{*lgrmessage.NewImage(imgData, summary...)})
	return m
}

// AddImageStream 添加图片消息元素，使用io.ReadSeeker
func (m *Message) AddImageStream(r io.ReadSeeker, summary ...string) *Message {
	*m = append(*m, &Image{*lgrmessage.NewStreamImage(r, summary...)})
	return m
}

// AddImageFile 添加图片消息元素，使用文件路径
func (m *Message) AddImageFile(filePath string, summary ...string) *Message {
	imgElement, err := lgrmessage.NewFileImage(filePath, summary...)
	if err != nil {
		log.Errorf("打开位于 %s 的图片时失败: %v", filePath, err)
	}
	*m = append(*m, &Image{*imgElement})
	return m
}
