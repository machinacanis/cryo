package cryo

import (
	"fmt"
	lgrMsg "github.com/LagrangeDev/LagrangeGo/message"
	"github.com/go-json-experiment/json"
	uuid "github.com/satori/go.uuid"
	"github.com/skip2/go-qrcode"
	"math/rand"
	"strconv"
	"time"
)

var botClientCount = 0

func randomDeviceNumber() int {
	return rand.Intn(9999999-1000000+1) + 1000000
}

func newUUID() string {
	return uuid.NewV4().String()
}

func newNickname() string {
	botClientCount++
	return fmt.Sprintf("Bot%d", botClientCount-1)
}

// getQRCodeString 生成二维码字符串
//
// 基于 https://github.com/Baozisoftware/qrcode-terminal-go 修改而来
func getQRCodeString(content string) (result *string) {
	var qr *qrcode.QRCode
	var err error
	qr, err = qrcode.New(content, qrcode.Low)
	if err != nil {
		return nil
	}
	data := qr.Bitmap()

	str := ""
	for ir, row := range data {
		lr := len(row)
		if ir == 0 || ir == 1 || ir == 2 ||
			ir == lr-1 || ir == lr-2 || ir == lr-3 {
			continue
		}
		for ic, col := range row {
			lc := len(data)
			if ic == 0 || ic == 1 || ic == 2 ||
				ic == lc-1 || ic == lc-2 || ic == lc-3 {
				continue
			}
			if col {
				str += fmt.Sprint("\033[48;5;0m  \033[0m") // 前景色
			} else {
				str += fmt.Sprint("\033[48;5;7m  \033[0m") // 背景色
			}
		}
		str += fmt.Sprintln()
	}
	return &str
}

func IMessageElementsToString(elements []lgrMsg.IMessageElement) string {
	// 将LagrangeGo的消息元素列表转换为字符串
	var result string
	for _, element := range elements {
		switch e := element.(type) {
		case *lgrMsg.TextElement:
			result += e.Content
		case *lgrMsg.AtElement:
			result += fmt.Sprintf("[@%d %s]", e.TargetUin, e.Display)
		case *lgrMsg.FaceElement:
			result += fmt.Sprintf("[表情 %d]", e.FaceID)
		case *lgrMsg.ReplyElement:
			result += fmt.Sprintf(
				"[回复 %d 于 %s 发送的消息 %d]",
				e.SenderUin, time.Unix(int64(e.Time), 0).Format("2006-01-02 15:04:05"),
				e.Elements,
			)
		case *lgrMsg.VoiceElement:
			result += fmt.Sprintf("[语音 %ds]", e.Duration)
		case *lgrMsg.ImageElement:
			result += fmt.Sprintf("[图片 %s]", e.URL)
		case *lgrMsg.FileElement:
			result += fmt.Sprintf("[文件 %s]", e.FileURL)
		case *lgrMsg.ShortVideoElement:
			result += fmt.Sprintf("[视频 %s]", e.URL)
		case *lgrMsg.LightAppElement:
			result += fmt.Sprintf("[轻应用 %s]", e.AppName)
		case *lgrMsg.XMLElement:
			result += fmt.Sprintf("[服务 %d]", e.ServiceID)
		case *lgrMsg.ForwardMessage:
			result += "[转发消息]"
		case *lgrMsg.MarketFaceElement:
			result += fmt.Sprintf("[魔法表情 %s]", e.Summary)
		default:
			continue
		}
	}
	return result
}

// ToJson 将数据转换为JSON格式
func ToJson(e any) []byte {
	res, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return res
}

// ToJsonString 将数据转换为JSON格式的字符串
func ToJsonString(e any) string {
	res, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return string(res)
}

// ProcessMessageContent 处理消息内容
func ProcessMessageContent(args ...interface{}) *Message {
	result := Message{} // 创建一个新的消息对象
	for _, arg := range args {
		// 遍历参数并根据其类型进行处理
		switch v := arg.(type) {
		case string:
			// 如果参数是字符串，则将其添加到消息元素中
			result.AddText(v)
		case int:
			result.AddText(strconv.Itoa(v))
		case int8:
			result.AddText(strconv.FormatInt(int64(v), 10))
		case int16:
			result.AddText(strconv.FormatInt(int64(v), 10))
		case int32:
			result.AddText(strconv.FormatInt(int64(v), 10))
		case int64:
			result.AddText(strconv.FormatInt(v, 10))
		case uint32:
			result.AddText(strconv.FormatUint(uint64(v), 10))
		case uint64:
			result.AddText(strconv.FormatUint(v, 10))
		case Message:
			// 如果参数是CryoMessage，则将其添加到消息元素中
			result.Add(v...)
		}
	}
	return &result
}

func Contains[T string | uint32 | int | EventType](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
