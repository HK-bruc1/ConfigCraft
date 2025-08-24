# ConfigCraft 🛠️

<div align="center">

**通用配置管理可视化工具**  
*将复杂配置文件转换为用户友好的图形界面*

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Platform](https://img.shields.io/badge/Platform-Windows-lightgrey.svg)](https://github.com/ConfigCraft/configcraft)
[![Release](https://img.shields.io/github/v/release/ConfigCraft/configcraft)](https://github.com/ConfigCraft/configcraft/releases)

[English](README.md) | **中文**

</div>

## 🚀 ConfigCraft 是什么？

ConfigCraft 是一款强大的配置管理可视化工具，将复杂的YAML配置文件转换为直观的图形化界面。最初为固件配置管理而开发，现已演化为适用于任何结构化配置工作流程的通用解决方案。

**核心特性：**
- 📊 **Schema驱动界面**：根据YAML结构自动生成表单控件
- 🎯 **通用支持**：适用于任何基于YAML的配置结构  
- 🖥️ **原生体验**：精美的跨平台GUI界面，原生文件对话框
- 🔄 **双重输出**：维护YAML配置的同时生成自定义格式文件
- ⚡ **零学习成本**：直观界面无需技术背景
- 🎨 **现代设计**：简洁专业的界面设计，智能状态显示

## 🎯 应用场景

- **固件配置管理**：硬件参数可视化配置与验证
- **应用程序设置**：将复杂配置文件转换为用户友好的表单
- **DevOps工具**：为非技术用户简化部署配置
- **配置模板系统**：创建可重用的配置模式
- **多格式输出**：从单一源文件生成多种配置格式

## 🚀 快速开始

### 环境要求
- **Windows 10/11**（主要平台）
- **Go 1.21+**（从源码构建需要）
- **TDM-GCC 10.3.0**（CGO编译需要）

### 安装方式

**方式一：下载发布版本**
```bash
# 从GitHub下载最新版本
curl -LO https://github.com/ConfigCraft/configcraft/releases/latest/download/configcraft.exe
```

**方式二：源码构建**
```bash
# 克隆仓库
git clone https://github.com/ConfigCraft/configcraft.git
cd configcraft

# 快速构建（推荐）
make build

# 或使用构建脚本
build\build.bat

# 或手动构建
go build -ldflags "-s -w -H windowsgui" -o build\configcraft.exe main.go
```

### 基本使用

1. **启动应用程序**
   ```bash
   # 运行GUI版本
   .\build\configcraft.exe
   
   # 或运行CLI版本（用于自动化）
   cd cmd && go run cli.go
   ```

2. **加载配置文件**
   - 点击"打开配置"选择YAML文件
   - ConfigCraft自动识别schema文件与配置文件
   - 使用左侧树形导航浏览配置分组

3. **编辑配置**
   - 从树形导航选择配置分组
   - 使用生成的表单控件修改数值
   - 查看实时验证和帮助信息

4. **保存结果**
   - 点击"保存配置"保存更改
   - 同时生成YAML配置文件和自定义输出格式
   - 文件命名保持一致：`config.yaml` + `config.conf`

## 📁 项目结构

```
configcraft/
├── internal/
│   ├── version/           # 版本管理
│   ├── config/           # 解析器和生成器引擎  
│   ├── models/           # 数据结构和类型
│   └── ui/               # GUI组件和逻辑
│       └── components/   # 自定义UI控件
├── assets/schemas/       # 示例schema文件
├── build/               # 构建产物和脚本
├── docs/                # 附加文档
├── cmd/                 # CLI版本
└── main.go              # 应用程序入口点
```

## 🛠️ 配置Schema格式

ConfigCraft使用YAML schema来定义配置结构：

```yaml
sections:
  section_name:
    name: "显示名称"
    groups:
      group_name:
        name: "分组显示名称"
        fields:
          field_name:
            type: "select"  # select, combo, number, boolean, text
            label: "字段标签"
            description: "字段下方显示的帮助文本"
            tooltip: "弹出窗口中的详细信息"
            placeholder: "输入提示文本"
            options:
              - value: "option1"
                label: "选项1"
              - value: "option2" 
                label: "选项2"
            default: "option1"
            required: true
```

**支持的字段类型：**
- `select`: 预定义选项的下拉菜单
- `combo`: 可编辑下拉菜单（预设+自定义输入）
- `number`: 带验证的数字输入
- `boolean`: 复选框控件
- `text`: 自由文本输入

## 🔧 开发指南

### 开发环境搭建

```bash
# 克隆并设置项目
git clone https://github.com/ConfigCraft/configcraft.git
cd configcraft

# 安装依赖和设置
make deps

# 开发模式运行
make dev

# 运行测试
make test
```

### 关键开发准则

- **版本更新**：仅修改 `internal/version/version.go`
- **UI组件**：遵循 `internal/ui/components/` 中的现有模式
- **配置逻辑**：在 `internal/config/parser.go` 中扩展新格式支持
- **错误处理**：提供有意义的错误信息和上下文
- **文档维护**：为任何API更改更新相关文档

## 🤝 参与贡献

我们欢迎贡献！请查看我们的[贡献指南](CONTRIBUTING.md)了解详情。

**贡献步骤：**
1. Fork 仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

**需要帮助的领域：**
- [ ] 跨平台测试与支持
- [ ] 附加配置格式支持
- [ ] 大型配置文件性能优化
- [ ] 高级验证功能
- [ ] 国际化支持

## 📊 技术特色

- **自定义树形导航**：解决Fyne框架树形控件闪烁问题，使用VBox实现
- **原生文件对话框**：集成zenity库，提供Windows原生文件选择体验
- **智能路径显示**：相对路径智能截取显示（最多2级目录）
- **版本同步管理**：跨所有UI元素的集中版本管理
- **Schema驱动架构**：从YAML定义零配置生成UI界面
- **跨平台基础**：基于Go和Fyne构建，为未来平台扩展奠定基础

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- **Fyne框架**：跨平台GUI开发支持
- **Zenity库**：原生对话框集成
- **Go社区**：优秀的工具链和生态系统
- **所有贡献者**：感谢让ConfigCraft变得更好！

---

<div align="center">

**🌟 如果这个项目对您有帮助，请给我们一个Star！**

[🐛 报告问题](https://github.com/ConfigCraft/configcraft/issues) | 
[💡 功能请求](https://github.com/ConfigCraft/configcraft/issues) | 
[📖 文档](https://github.com/ConfigCraft/configcraft/wiki) |
[🔄 更新日志](./CHANGELOG.md)

*ConfigCraft - 让配置管理变得简单高效* 

</div>