package cryo

import (
	"fmt"
	lgrmessage "github.com/LagrangeDev/LagrangeGo/message"
	"time"
)

// ElementType 消息元素类型
type ElementType int

const (
	TextType       ElementType = iota // 文本元素类型
	ImageType                         // 图片元素类型
	FaceType                          // 表情元素类型
	AtType                            // At元素类型
	ReplyType                         // 回复元素类型
	ServiceType                       // 服务元素类型
	ForwardType                       // 转发元素类型
	FileType                          // 文件元素类型
	VoiceType                         // 语音元素类型
	VideoType                         // 视频元素类型
	LightAppType                      // 轻应用元素类型
	RedBag                            // 红包（无实际意义）
	MarketFaceType                    // 魔法表情元素类型
)

type MessageElement interface {
	GetType() ElementType                           // 获取元素类型
	GetLgrElementType() lgrmessage.ElementType      // 获取LagrangeGo元素类型
	GetIMessageElement() lgrmessage.IMessageElement // 获取LagrangeGo的消息元素
	ToString() string                               // 获取元素的字符串表示
}

// Text 文本元素
type Text struct {
	lgrmessage.TextElement
}

// At 提及元素
type At struct {
	lgrmessage.AtElement
}

// Face 表情元素
type Face struct {
	lgrmessage.FaceElement
}

// Reply 回复元素
type Reply struct {
	lgrmessage.ReplyElement
}

// Voice 语音元素
type Voice struct {
	lgrmessage.VoiceElement
}

// Image 图片元素
type Image struct {
	lgrmessage.ImageElement
}

// File 文件元素
type File struct {
	lgrmessage.FileElement
}

// ShortVideo 短视频元素
type ShortVideo struct {
	lgrmessage.ShortVideoElement
}

// LightApp 轻应用元素
type LightApp struct {
	lgrmessage.LightAppElement
}

// XML 服务元素
type XML struct {
	lgrmessage.XMLElement
}

// ForwardMessage 转发消息元素
type ForwardMessage struct {
	lgrmessage.ForwardMessage
}

// MarketFace 魔法表情元素
type MarketFace struct {
	lgrmessage.MarketFaceElement
}

func (e *Text) GetType() ElementType           { return TextType }
func (e *At) GetType() ElementType             { return AtType }
func (e *Face) GetType() ElementType           { return FaceType }
func (e *Reply) GetType() ElementType          { return ReplyType }
func (e *Voice) GetType() ElementType          { return VoiceType }
func (e *Image) GetType() ElementType          { return ImageType }
func (e *File) GetType() ElementType           { return FileType }
func (e *ShortVideo) GetType() ElementType     { return VideoType }
func (e *LightApp) GetType() ElementType       { return LightAppType }
func (e *XML) GetType() ElementType            { return ServiceType }
func (e *ForwardMessage) GetType() ElementType { return ForwardType }
func (e *MarketFace) GetType() ElementType     { return MarketFaceType }

func (e *Text) GetLgrElementType() lgrmessage.ElementType {
	return lgrmessage.Text
}
func (e *At) GetLgrElementType() lgrmessage.ElementType {
	return lgrmessage.At
}
func (e *Face) GetLgrElementType() lgrmessage.ElementType {
	return lgrmessage.Face
}
func (e *Reply) GetLgrElementType() lgrmessage.ElementType {
	return lgrmessage.Reply
}
func (e *Voice) GetLgrElementType() lgrmessage.ElementType {
	return lgrmessage.Voice
}
func (e *Image) GetLgrElementType() lgrmessage.ElementType {
	return lgrmessage.Image
}
func (e *File) GetLgrElementType() lgrmessage.ElementType {
	return lgrmessage.File
}
func (e *ShortVideo) GetLgrElementType() lgrmessage.ElementType {
	return lgrmessage.Video
}
func (e *LightApp) GetLgrElementType() lgrmessage.ElementType {
	return lgrmessage.LightApp
}
func (e *XML) GetLgrElementType() lgrmessage.ElementType {
	return lgrmessage.Service
}
func (e *ForwardMessage) GetLgrElementType() lgrmessage.ElementType {
	return lgrmessage.Forward
}
func (e *MarketFace) GetLgrElementType() lgrmessage.ElementType {
	return lgrmessage.MarketFace
}

func (e *Text) GetIMessageElement() lgrmessage.IMessageElement {
	return &e.TextElement
}

func (e *At) GetIMessageElement() lgrmessage.IMessageElement {
	return &e.AtElement
}

func (e *Face) GetIMessageElement() lgrmessage.IMessageElement {
	return &e.FaceElement
}

func (e *Reply) GetIMessageElement() lgrmessage.IMessageElement {
	return &e.ReplyElement
}

func (e *Voice) GetIMessageElement() lgrmessage.IMessageElement {
	return &e.VoiceElement
}

func (e *Image) GetIMessageElement() lgrmessage.IMessageElement {
	return &e.ImageElement
}

func (e *File) GetIMessageElement() lgrmessage.IMessageElement {
	return &e.FileElement
}

func (e *ShortVideo) GetIMessageElement() lgrmessage.IMessageElement {
	return &e.ShortVideoElement
}

func (e *LightApp) GetIMessageElement() lgrmessage.IMessageElement {
	return &e.LightAppElement
}

func (e *XML) GetIMessageElement() lgrmessage.IMessageElement {
	return &e.XMLElement
}

func (e *ForwardMessage) GetIMessageElement() lgrmessage.IMessageElement {
	return &e.ForwardMessage
}

func (e *MarketFace) GetIMessageElement() lgrmessage.IMessageElement {
	return &e.MarketFaceElement
}

func (e *Text) ToString() string {
	return e.Content
}

func (e *At) ToString() string {
	return fmt.Sprintf("[@%d %s]", e.TargetUin, e.Display)
}

func (e *Face) ToString() string {
	return fmt.Sprintf("[表情 %d]", e.FaceID)
}

func (e *Reply) ToString() string {
	return fmt.Sprintf(
		"[回复 %d 于 %s 发送的消息 %s]",
		e.SenderUin, time.Unix(int64(e.Time), 0).Format("2006-01-02 15:04:05"),
		IMessageElementsToString(e.Elements),
	) // DONE：有空记得改了这里的回复消息显示
}

func (e *Voice) ToString() string {
	return fmt.Sprintf("[语音 %ds]", e.Duration)
}

func (e *Image) ToString() string {
	return fmt.Sprintf("[图片 %s]", e.URL)
}

func (e *File) ToString() string {
	return fmt.Sprintf("[文件 %s]", e.FileURL)
}

func (e *ShortVideo) ToString() string {
	return fmt.Sprintf("[视频 %s]", e.URL)
}

func (e *LightApp) ToString() string {
	return fmt.Sprintf("[轻应用 %s]", e.AppName)
}

func (e *XML) ToString() string {
	return fmt.Sprintf("[服务 %d]", e.ServiceID)
}

func (e *ForwardMessage) ToString() string {
	return "[转发消息]"
}

func (e *MarketFace) ToString() string {
	return fmt.Sprintf("[魔法表情 %s]", e.Summary)
}
