package cryo

// Plugin 是 cryo 的插件接口
//
// 实现这个接口就可以创建一个插件，通过插件来扩展功能管理起来更方便
type Plugin interface {
	Init(bot *Bot) error          // 初始化插件，向插件中传递bot实例
	GetPluginName() string        // 获取插件名称信息
	GetPluginVersion() string     // 获取插件版本号信息
	GetPluginDescription() string // 获取插件描述信息
	GetPluginAuthor() string      // 获取插件作者信息
	Enable()                      // 启用插件
	Disable()                     // 禁用插件
	IsEnable() bool               // 是否已启用插件
}
