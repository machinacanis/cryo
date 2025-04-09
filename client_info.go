package cryo

import (
	"github.com/go-json-experiment/json"
	"os"
)

// ClientInfo 客户端持久化信息，用于保存Bot客户端数据，以便于自动登录
type ClientInfo struct {
	Id        string `json:"id"`
	Platform  string `json:"platform"`
	Version   string `json:"version"`
	DeviceNum int    `json:"device_num"`
	Signature string `json:"signature"`
	Uin       uint32 `json:"uin"`
	Uid       string `json:"uid"`
}

// ReadClientInfos 从文件读取客户端信息
func ReadClientInfos() ([]ClientInfo, error) {
	data, err := os.ReadFile("client_infos.json")
	if err != nil {
		return nil, err
	}

	var clientInfos []ClientInfo
	err = json.Unmarshal(data, &clientInfos)
	if err != nil {
		return nil, err
	}

	return clientInfos, nil
}

// WriteClientInfos 写入客户端信息到文件
func WriteClientInfos(clientInfos []ClientInfo) error {
	data, err := json.Marshal(clientInfos)
	if err != nil {
		return err
	}

	err = os.WriteFile("client_infos.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// SaveClientInfo 保存客户端信息到文件
func SaveClientInfo(clientInfo ClientInfo) error {
	// 首先尝试读取现有的客户端信息
	clientInfos, err := ReadClientInfos()
	if err != nil { // 读取失败，可能是文件不存在
		if os.IsNotExist(err) {
			// 文件不存在，创建一个新的切片
			clientInfos = []ClientInfo{}
		} else {
			return err // 其他错误
		}
	}
	// 检测是否已经存在botid相同的客户端信息
	updateFlag := false
	for _, info := range clientInfos {
		if info.Id == clientInfo.Id {
			// 如果存在，则更新该信息
			updateFlag = true
			break
		}
	}
	if updateFlag {
		// 如果存在，则更新该信息
		for i, info := range clientInfos {
			if info.Id == clientInfo.Id {
				clientInfos[i] = clientInfo
				break
			}
		}
	} else {
		// 如果不存在，则添加新的信息
		clientInfos = append(clientInfos, clientInfo)
	}
	err = WriteClientInfos(clientInfos)
	return err
}

// RemoveClientInfo 从文件中删除指定ID的客户端信息
func RemoveClientInfo(botId string) error {
	// 读取现有的客户端信息
	clientInfos, err := ReadClientInfos()
	if err != nil {
		return err // 读取失败，可能是文件不存在
	}

	// 创建一个新的切片来存储更新后的客户端信息
	var updatedClientInfos []ClientInfo

	// 遍历现有的客户端信息，排除要删除的项
	for _, info := range clientInfos {
		if info.Id != botId {
			updatedClientInfos = append(updatedClientInfos, info)
		}
	}

	// 将更新后的客户端信息写回文件
	err = WriteClientInfos(updatedClientInfos)
	return err
}
