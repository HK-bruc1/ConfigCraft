package version

// Version 应用程序版本信息
const Version = "0.3.6"

// AppName 应用程序名称
const AppName = "ConfigCraft"

// GetVersionString 获取完整版本字符串
func GetVersionString() string {
	return AppName + " v" + Version
}