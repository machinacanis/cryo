package cryo

import (
	"errors"
	"fmt"
	"github.com/machinacanis/cryo/log"
)

// Bot cryo 的Bot封装
//
// 提供了对Bot的操作和管理功能，可以通过 initFlag 来判断是否初始化完成
type Bot struct {
	initFlag         bool                       // 是否初始化完成
	connectedClients map[string]*LagrangeClient // 已连接的Bot客户端集合
	logger           log.CryoLogger             // 日志记录器
	bus              *EventBus                  // 事件总线
	conf             Config                     // 配置项
	plugin           []Plugin                   // 插件列表
}

// NewBot 创建一个新的CryoBot实例
func NewBot() *Bot {
	return &Bot{}
}

// Init 初始化cryobot
//
// 可以传入配置项来覆写默认配置，空的配置项会自动使用默认配置
//
// 如果本地配置文件存在，且没有传入配置项，则会自动加载本地配置文件
func (b *Bot) Init(logger log.CryoLogger, c ...Config) {
	defaultConfig := Config{
		SignServers:                  []string{DefaultSignServer},
		EnablePluginAutoLoad:         true,
		EnableClientAutoSave:         true,
		EnablePrintLogo:              true,
		EnableConnectPrintMiddleware: true,
		EnableMessagePrintMiddleware: true,
		EnableEventDebugMiddleware:   false,
	}
	b.logger = logger
	if len(c) == 0 { // 如果没有传入配置项，则尝试加载本地配置文件
		co, err := ReadCryoConfig()
		if err == nil {
			c = append(c, co)
			b.logger.Info("已正在加载本地配置文件")
		}
	}
	if len(c) > 0 {
		if c[0].SignServers != nil {
			defaultConfig.SignServers = c[0].SignServers
		}
		if c[0].EnablePluginAutoLoad {
			defaultConfig.EnablePluginAutoLoad = c[0].EnablePluginAutoLoad
		}
		if c[0].EnableClientAutoSave {
			defaultConfig.EnableClientAutoSave = c[0].EnableClientAutoSave
		}
		if c[0].EnablePrintLogo {
			defaultConfig.EnablePrintLogo = c[0].EnablePrintLogo
		}
		if c[0].EnableConnectPrintMiddleware {
			defaultConfig.EnableConnectPrintMiddleware = c[0].EnableConnectPrintMiddleware
		}
		if c[0].EnableMessagePrintMiddleware {
			defaultConfig.EnableMessagePrintMiddleware = c[0].EnableMessagePrintMiddleware
		}
		if c[0].EnableEventDebugMiddleware {
			defaultConfig.EnableEventDebugMiddleware = c[0].EnableEventDebugMiddleware
		}
	}
	b.conf = defaultConfig // 初始化配置

	// 初始化事件总线
	fmt.Print(log.Logo)
	b.logger.Infof("[Cryo] 🧊cryobot 正在初始化...")
	b.bus = NewEventBus() // 初始化事件总线
	// 初始化连接的客户端集合
	b.connectedClients = make(map[string]*LagrangeClient)
	// 设置连接打印中间件
	// setConnectPrintMiddleware()
	// 设置消息打印中间件
	// setMessagePrintMiddleware()
	// 设置事件调试中间件
	setDefaultMiddleware(b.bus, b.logger, b.conf)

	b.initFlag = true
}

// IsInit 判断是否初始化完成
func (b *Bot) IsInit() bool {
	return b.initFlag
}

// Start 启动cryobot
func (b *Bot) Start() error {
	if !b.initFlag {
		// 没有进行初始化
		b.logger.Error("cryobot 没有进行初始化，请先调用 Init() 函数进行初始化！")
		return errors.New("cryobot 没有进行初始化，请先调用 Init() 函数进行初始化！")
	}
	select {} // 阻塞主线程，运行事件循环
}

// AutoConnect 自动尝试建立连接，如果没有已保存的连接信息或已保存的连接信息无效，则尝试创建并连接新的bot客户端
//
// 如果已经有了连接过的bot客户端，则会跳过自动连接过程，你应该手动使用 ConnectNewClient 来新建连接
func (b *Bot) AutoConnect() error {
	if !b.initFlag {
		// 没有进行初始化
		b.logger.Error("cryobot 没有进行初始化，请先调用 Init() 函数进行初始化！")
		return errors.New("cryobot 没有进行初始化，请先调用 Init() 函数进行初始化！")
	}
	// 首先检测是否已经连接
	if len(b.connectedClients) > 0 {
		// 跳过自动连接
		return nil
	}
	// 尝试连接所有已保存的bot客户端
	b.ConnectAllSavedClient()
	// 如果没有连接成功，则尝试连接新的bot客户端
	retriedCount := 0
	for len(b.connectedClients) == 0 && retriedCount < 3 {
		b.ConnectNewClient()
		retriedCount++
	}
	if len(b.connectedClients) == 0 {
		b.logger.Error("达到最大重试次数，cryobot 无法连接到bot客户端，请检查网络或配置文件")
		return errors.New("达到最大重试次数，cryobot 无法连接到bot客户端，请检查网络或配置文件")
	}
	return nil
}

// ConnectSavedClient 尝试查询并连接到指定的bot客户端
func (b *Bot) ConnectSavedClient(info ClientInfo) bool {
	c := NewLagrangeClient()
	c.Init(b.bus, b.logger, b.conf)
	if !c.Rebuild(info) {
		return false
	}
	b.logger.Infof("[Cryo] 正在连接 %s：%s (%d)", c.Nickname, c.Id, c.Uin)
	if !c.SignatureLogin() {
		return false
	}
	b.connectedClients[c.Id] = c
	return true
}

// ConnectNewClient 尝试连接一个新的bot客户端
func (b *Bot) ConnectNewClient() bool {
	c := NewLagrangeClient()
	c.Init(b.bus, b.logger, b.conf)
	b.logger.Infof("[Cryo] 正在连接 %s：%s (%d)", c.Nickname, c.Id, c.Uin)
	if !c.QRCodeLogin() {
		return false
	}
	b.connectedClients[c.Id] = c
	return true
}

// ConnectAllSavedClient 尝试连接所有已保存的bot客户端
func (b *Bot) ConnectAllSavedClient() {
	// 读取历史连接的客户端
	clientInfos, err := ReadClientInfos()
	if err != nil {
		b.logger.Error("读取Bot信息时出现错误：", err)
		return
	}
	if len(clientInfos) == 0 {
		b.logger.Info("没有找到Bot信息")
		return
	}
	for _, info := range clientInfos {
		if !b.ConnectSavedClient(info) {
			b.logger.Error("通过历史记录连接Bot客户端失败")
			b.logger.Error("已自动清除失效的客户端信息，请重新登录")
		}
	}
}

// GetClientById 获取指定ID的bot客户端
func (b *Bot) GetClientById(id string) *LagrangeClient {
	if client, ok := b.connectedClients[id]; ok {
		return client
	}
	return nil
}

// GetClientByUin 获取指定Uin的bot客户端
func (b *Bot) GetClientByUin(uin uint32) *LagrangeClient {
	for _, client := range b.connectedClients {
		if client.Uin == uin {
			return client
		}
	}
	return nil
}

// GetClientByUid 获取指定Uid的bot客户端
func (b *Bot) GetClientByUid(uid string) *LagrangeClient {
	for _, client := range b.connectedClients {
		if client.Uid == uid {
			return client
		}
	}
	return nil
}

// GetClient 获取指定事件对应的bot客户端
func (b *Bot) GetClient(event Event) *LagrangeClient {
	return b.GetClientById(event.GetUniEvent().BotId)
}

// GetBus 获取事件总线
func (b *Bot) GetBus() *EventBus {
	return b.bus
}

// AddPlugin 添加插件
func (b *Bot) AddPlugin(plugin ...Plugin) {
	for _, p := range plugin {
		err := p.Init(b)
		if err != nil {
			b.logger.Errorf("[Cryo] 插件 %s 初始化失败：%v", p.GetPluginName(), err)
		}
		if b.conf.EnablePluginAutoLoad { // 如果启用自动加载插件
			b.logger.Successf("[Cryo] 插件 %s 已成功加载", p.GetPluginName())
			p.Enable()
		}
		b.plugin = append(b.plugin, p)
	}
}

// GetPlugin 获取插件
func (b *Bot) GetPlugin(name string) []Plugin {
	var plugins []Plugin
	for _, p := range b.plugin {
		if p.GetPluginName() == name {
			plugins = append(plugins, p)
		}
	}
	return plugins
}

// GetEnabledPlugin 获取启用的插件
func (b *Bot) GetEnabledPlugin() []Plugin {
	var plugins []Plugin
	for _, p := range b.plugin {
		if p.IsEnable() {
			plugins = append(plugins, p)
		}
	}
	return plugins
}

// GetDisabledPlugin 获取禁用的插件
func (b *Bot) GetDisabledPlugin() []Plugin {
	var plugins []Plugin
	for _, p := range b.plugin {
		if !p.IsEnable() {
			plugins = append(plugins, p)
		}
	}
	return plugins
}

// GetAllPlugin 获取所有插件
func (b *Bot) GetAllPlugin() []Plugin {
	return b.plugin
}

// RemovePlugin 移除插件
func (b *Bot) RemovePlugin(plugin ...Plugin) {
	for _, p := range plugin {
		for i, pl := range b.plugin {
			if pl.GetPluginName() == p.GetPluginName() {
				b.plugin = append(b.plugin[:i], b.plugin[i+1:]...)
				b.logger.Infof("[Cryo] 插件 %s 已卸载", p.GetPluginName())
				break
			}
		}
	}
}

// GetLogger 获取日志记录器
func (b *Bot) GetLogger() log.CryoLogger {
	return b.logger
}

// GetConfig 获取配置项
func (b *Bot) GetConfig() Config {
	return b.conf
}

// Send 快速根据事件内容发送消息
func (b *Bot) Send(event MessageEvent, args ...interface{}) (ok bool, messageId uint32) {
	// 根据事件获取对应的bot客户端
	return b.GetClient(event).Send(event, args...)
}

// Reply 快速根据事件内容回复消息
func (b *Bot) Reply(event MessageEvent, args ...interface{}) (ok bool, messageId uint32) {
	// 根据事件获取对应的bot客户端
	return b.GetClient(event).Reply(event, args...)
}
