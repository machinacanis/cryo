package cryo

import (
	"encoding/base64"
	"fmt"
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/client/auth"
	"github.com/machinacanis/cryo/log"
	"os"
	"time"
)

// LagrangeClient cryo的Bot客户端封装
type LagrangeClient struct {
	Id        string
	Client    *client.QQClient
	Platform  string
	Version   string
	DeviceNum int
	Uin       uint32
	Uid       string
	Nickname  string

	initFlag bool   // 是否初始化完成
	conf     Config // 配置项
	bus      *EventBus
}

// NewLagrangeClient 创建一个新的LagrangeClient实例
func NewLagrangeClient() *LagrangeClient {
	return &LagrangeClient{}
}

// Init 初始化一个新的LagrangeClient客户端
func (c *LagrangeClient) Init(bus *EventBus, conf Config) {
	c.Id = newUUID() // 给Bot客户端分配一个唯一的UUID
	c.conf = conf
	c.bus = bus

	// 默认平台和版本
	if c.Platform == "" {
		c.Platform = "linux"
	}
	if c.Version == "" {
		c.Version = "3.2.15-30366"
	}

	appInfo := auth.AppList[c.Platform][c.Version]
	c.Client = client.NewClient(0, "")
	c.Client.SetLogger(log.PLogger) // 替换日志记录器，详见client/protocol_logger.go以及log/logger.go
	c.Client.UseVersion(appInfo)
	c.Client.AddSignServer(conf.SignServers...)
	c.DeviceNum = randomDeviceNumber()
	c.Client.UseDevice(auth.NewDeviceInfo(c.DeviceNum))
	c.Nickname = newNickname() // 生成一个默认的编号昵称

	c.initFlag = true
}

// Rebuild 重新构建LagrangeClient实例
func (c *LagrangeClient) Rebuild(clientInfo ClientInfo) bool {
	if !c.initFlag {
		log.Error("cryobot客户端没有完成初始化，请先调用Init()方法")
		return false
	}
	var sig string
	c.Id = clientInfo.Id
	c.Platform = clientInfo.Platform
	c.Version = clientInfo.Version
	c.DeviceNum = clientInfo.DeviceNum
	c.Uin = clientInfo.Uin
	c.Uid = clientInfo.Uid
	sig = clientInfo.Signature
	c.Client.UseDevice(auth.NewDeviceInfo(c.DeviceNum))
	c.Client.UseVersion(auth.AppList[c.Platform][c.Version])
	c.UseSignature(sig) // 使用指定的签名信息
	return true
}

// Save 将当前客户端的信息保存到文件中
func (c *LagrangeClient) Save() error {
	clientInfo := ClientInfo{
		Id:        c.Id,
		Signature: c.GetSignature(),
		Platform:  c.Platform,
		Version:   c.Version,
		DeviceNum: c.DeviceNum,
		Uin:       c.Uin,
		Uid:       c.Uid,
	}
	err := SaveClientInfo(clientInfo)
	return err
}

// GetSignature 获取当前客户端的签名信息
func (c *LagrangeClient) GetSignature() string {
	data, err := c.Client.Sig().Marshal()
	if err != nil {
		log.Error("序列化签名时出现错误：", err)
		return ""
	}
	// 将二进制的签名直接编码到字符串
	sig := base64.StdEncoding.EncodeToString(data)
	return sig
}

// UseSignature 使用指定的签名信息
func (c *LagrangeClient) UseSignature(sig string) {
	// 将字符串解码为二进制
	data, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		log.Error("解码签名时出现错误：", err)
		return
	}
	// 反序列化签名
	sigInfo, err := auth.UnmarshalSigInfo(data, true)
	if err != nil {
		log.Error("反序列化签名时出现错误：", err)
		return
	}
	c.Client.UseSig(sigInfo)
}

func (c *LagrangeClient) AfterLogin() {
	// 登录成功后，保存签名
	c.Uin = c.Client.Sig().Uin
	c.Uid = c.Client.Sig().UID
	c.sendBotConnectedEvent()        // 发送登录成功事件
	if c.conf.EnableClientAutoSave { // 如果启用了自动保存
		err := c.Save()
		if err != nil {
			log.Error("保存登录信息时出现错误：", err)
		} // 保存登录信息
	}

	// 订阅事件
	EventBind(c)
}

// GetQRCode 获取二维码信息
func (c *LagrangeClient) GetQRCode() ([]byte, string, error) {
	code, res, err := c.Client.FetchQRCodeDefault()
	// 这里获取到两个参数，第一个是字节形式的二维码图片，第二个是二维码指向的链接
	return code, res, err
}

// SaveQRCode 保存二维码图片
func (c *LagrangeClient) SaveQRCode(code []byte) bool {
	qrcodePath := fmt.Sprintf("QRCode_%s.png", c.Id)
	err := os.WriteFile(qrcodePath, code, 0644)
	if err != nil {
		log.Error("写入二维码图片时出现错误：", err)
		return false
	}
	log.Infof("登录二维码已保存到 %s", qrcodePath)
	return true
}

// PrintQRCode 打印二维码
func (c *LagrangeClient) PrintQRCode(url string) {
	// 打印二维码的链接
	fmt.Println(*getQRCodeString(url)) // 注意使用了指针
}

// SignatureLogin 使用签名快速登录
func (c *LagrangeClient) SignatureLogin() (ok bool) {
	sig := c.Client.Sig()
	if sig != nil {
		err := c.Client.FastLogin()
		if err == nil {
			// 通过保存的签名快速登录成功
			c.AfterLogin()
			return true
		}
	}
	return false
}

// QRCodeLogin 使用二维码登录
func (c *LagrangeClient) QRCodeLogin() bool {
	log.Info("正在使用二维码登录...")
	code, url, err := c.GetQRCode()
	if err != nil {
		log.Error("获取二维码时出现错误：", err)
		return false
	}
	// 保存二维码图片
	c.SaveQRCode(code)
	// 向终端输出二维码
	c.PrintQRCode(url)
	if !c.watingForLoginResult() { // 等待扫码登录
		log.Warn("扫码登录失败！")
		return false
	}
	c.AfterLogin()
	return true
}

// watingForLoginResult 等待扫码登录结果
func (c *LagrangeClient) watingForLoginResult() bool {
	//轮询登录状态
	for {
		retCode, err := c.Client.GetQRCodeResult()
		if err != nil {
			log.Error("获取二维码登录结果时出现错误：", err)
			return false
		}
		// 等待扫码
		if retCode.Waitable() {
			time.Sleep(1 * time.Second)
			continue
		}
		if !retCode.Success() {
			return false
		}
		break
	}
	_, err := c.Client.QRCodeLogin()
	if err != nil {
		log.Error("二维码登录时出现错误：", err)
		return false
	}
	return true
}

// sendBotConnectedEvent 发送bot连接事件
func (c *LagrangeClient) sendBotConnectedEvent() {
	// 发送bot连接事件
	c.bus.Publish(&BotConnectedEvent{
		UniEvent: UniEvent{
			payload:     nil,
			EventType:   BotConnectedEventType,
			EventId:     newUUID(),
			EventTags:   []string{"cryo", "bot_connected"},
			Time:        uint32(time.Now().Unix()),
			botClient:   c,
			BotId:       c.Id,
			BotNickname: c.Nickname,
			BotUin:      c.Uin,
			BotUid:      c.Uid,
			Platform:    c.Platform,
		},
		Version: c.Version,
	})
}

// SendPrivateMessage 发送私聊消息
func (c *LagrangeClient) SendPrivateMessage(userUin uint32, msg *Message) (ok bool, messageId uint32) {
	// 发送私聊消息
	message, err := c.Client.SendPrivateMessage(userUin, msg.ToIMessageElements())
	if err != nil {
		log.Errorf("向用户 %d 发送消息时出现错误：%v", userUin, err)
		return false, 0
	}
	return true, message.ID
}

// SendGroupMessage 发送群消息
func (c *LagrangeClient) SendGroupMessage(groupUin uint32, msg *Message) (ok bool, messageId uint32) {
	// 发送群消息
	message, err := c.Client.SendGroupMessage(groupUin, msg.ToIMessageElements())
	if err != nil {
		log.Errorf("向群 %d 发送消息时出现错误：%v", groupUin, err)
		return false, 0
	}
	return true, message.ID
}

// SendTempMessage 发送临时消息
func (c *LagrangeClient) SendTempMessage(groupUin, userUin uint32, msg *Message) (ok bool, messageId uint32) {
	// 发送临时消息
	message, err := c.Client.SendTempMessage(groupUin, userUin, msg.ToIMessageElements())
	if err != nil {
		log.Errorf("向与用户 %d 的临时会话发送消息时出现错误：%v", groupUin, err)
		return false, 0
	}
	return true, message.ID
}

func (c *LagrangeClient) Send(event MessageEvent, args ...interface{}) (ok bool, messageId uint32) {
	// 处理消息内容
	m := ProcessMessageContent(args...)
	// 根据传入的事件来发送消息
	switch event.GetEventType() {
	case PrivateMessageEventType:
		return c.SendPrivateMessage(event.GetUniMessageEvent().SenderUin, m)
	case GroupMessageEventType:
		return c.SendPrivateMessage(event.GetUniMessageEvent().SenderUin, m)
	case TempMessageEventType:
		return c.SendPrivateMessage(event.GetUniMessageEvent().SenderUin, m)
	case UniMessageEventType:
		me := event.GetUniMessageEvent()
		// 通过tag来判断消息类型
		if Contains(me.EventTags, "private_message") {
			return c.SendPrivateMessage(me.SenderUin, m)
		} else if Contains(me.EventTags, "group_message") {
			return c.SendGroupMessage(me.GroupUin, m)
		} else if Contains(me.EventTags, "temp_message") {
			return c.SendTempMessage(me.GroupUin, me.SenderUin, m)
		}
	default:
		log.Error("发送消息时传入了不支持的消息事件")
	}
	return false, 0
}

func (c *LagrangeClient) Reply(event MessageEvent, args ...interface{}) (ok bool, messageId uint32) {
	// 处理消息内容
	m := Message{}
	m.AddReply(event).Add(*ProcessMessageContent(args...)...)
	return c.Send(event, m)
}
