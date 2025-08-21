# DHF Configuration Manager

DHF耳机固件配置可视化管理工具 - 将复杂的DHF配置文件通过图形界面进行可视化编辑。

## 项目概述

DHF Configuration Manager 是为 DHF AC710N-V300P03 SDK 开发的轻量级配置管理工具，旨在简化复杂的耳机固件配置过程。

### 技术栈
- **语言**: Go 1.21+
- **GUI框架**: Fyne v2.4.3
- **配置格式**: YAML → DHF conf
- **打包方式**: 单exe可执行文件

## 已实现功能 ✅

### 核心架构
- [x] 基于YAML schema的动态UI生成系统
- [x] 模块化的代码架构 (MVC模式)
- [x] 配置文件解析和生成引擎
- [x] 跨平台Go应用框架

### 用户界面
- [x] 现代化GUI界面 (Fyne)
- [x] 响应式布局 (左侧导航树 + 右侧编辑器)
- [x] 工具栏 (新建、打开、保存、导出)
- [x] 状态栏显示
- [x] 窗口大小自适应 (900x650默认尺寸)

### 配置管理
- [x] 完整的DHF配置模型支持
- [x] 基础配置 (IC型号、VM操作、功放控制等)
- [x] 通话按键配置 (单击、双击、长按等)
- [x] 音乐按键配置 (左右耳、TWS状态等)
- [x] 特殊按键功能 (工厂重置、DUT模式等)
- [x] LED灯效配置 (TWS、蓝牙、系统事件、充电等)
- [x] 工厂重置高级配置

### 文件操作
- [x] 标准DHF conf格式输出
- [x] YAML用户配置保存/加载
- [x] 配置项验证和默认值支持
- [x] 多种配置控件 (下拉框、复选框、数字输入)

### 命令行版本
- [x] 完整功能的CLI版本 (`cmd-version.go`)
- [x] 交互式配置界面
- [x] 快速配置验证工具

## 待实现功能 🚧

### 界面优化
- [ ] 中文字符显示乱码问题修复
- [ ] 自定义字体和主题支持
- [ ] 图标和视觉设计优化
- [ ] 键盘快捷键支持

### 功能增强
- [ ] 配置文件导入功能 (conf → YAML)
- [ ] 配置模板系统
- [ ] 配置项搜索功能
- [ ] 撤销/重做操作
- [ ] 配置验证和错误提示
- [ ] 多配置文件对比功能

### 用户体验
- [ ] 多语言支持 (中文/英文)
- [ ] 在线帮助文档
- [ ] 配置向导模式
- [ ] 最近使用的配置文件
- [ ] 自动保存功能

### 技术改进
- [ ] 单元测试覆盖
- [ ] 配置schema版本管理
- [ ] 插件系统支持
- [ ] 批量配置处理
- [ ] 配置文件压缩

## 快速开始

### 环境要求
- Go 1.21 或更高版本
- Windows 10/11 
- TDM-GCC 10.3.0 (用于CGO编译)

### 编译运行

```bash
# 克隆项目
git clone <repository-url>
cd DHFConfigTool

# 编译GUI版本 (推荐)
build\build.bat

# 或者手动编译
go build -ldflags "-s -w -H windowsgui" -o build\dhf-config-manager.exe main.go

# 直接运行GUI (开发模式)
go run main.go

# 运行CLI版本
cd cmd && go run cli.go
```

### 使用方法

1. **GUI版本**: 运行 `build\dhf-config-manager.exe`
2. **命令行版本**: `cd cmd && go run cli.go`
3. **配置编辑**: 左侧选择配置分组，右侧编辑具体参数
4. **导出配置**: 点击 Export 按钮生成标准 DHF conf 文件

## 项目结构

```
dhf-config-manager/
├── main.go                 # GUI应用程序入口
├── go.mod                  # Go模块定义
├── customer.conf           # 真实配置文件参考
├── internal/               # 核心业务逻辑
│   ├── config/            # 配置解析器和生成器
│   ├── models/            # 数据模型和类型定义
│   └── ui/                # 用户界面
│       ├── app.go         # 主应用程序
│       ├── theme.go       # 界面主题
│       └── components/    # UI组件
│           ├── tree.go    # 配置树组件
│           ├── editor.go  # 配置编辑器
│           └── toolbar.go # 工具栏组件
├── assets/                # 静态资源
│   └── schemas/           # YAML配置模板
│       ├── dhf-real-schema.yaml    # 真实配置schema
│       ├── dhf-schema-en.yaml      # 英文版schema
│       └── dhf-schema.yaml         # 中文版schema
├── build/                 # 构建相关
│   └── build.bat          # Windows构建脚本
├── cmd/                   # 命令行工具
│   └── cli.go             # CLI版本
├── README.md              # 项目说明
├── CHANGELOG.md           # 版本更新日志
└── .gitignore             # Git忽略文件
```

## 技术亮点

- **零学习成本**: 图形化界面，无需了解conf文件语法
- **弹性扩展**: 基于YAML schema动态生成UI，易于添加新配置项
- **高效输出**: 直接生成标准conf格式，无缝集成现有构建流程
- **轻量快速**: 单exe文件，启动速度快 (约18MB)

## 贡献指南

欢迎提交Issue和Pull Request来改进项目。

### 开发环境设置
1. 安装Go 1.21+
2. 安装TDM-GCC编译器
3. 克隆项目并安装依赖
4. 运行测试确保环境正常

## 版本历史

### v0.3.2 (2025-08-21)
- ✅ 修复配置分组位置随机变化问题
- ✅ 统一排序逻辑，确保界面显示一致性
- ✅ 消除Go map遍历随机性的影响
- ✅ 实现稳定的分组显示优先级

### v0.3.1 (2025-08-21)
- ✅ 彻底解决树形控件闪烁和位置漂移问题
- ✅ 实现自定义树形控件，替换有问题的Fyne Tree
- ✅ 采用简约箭头设计，提升界面优雅度
- ✅ 优化布局间距，实现紧凑美观的视觉效果
- ✅ 流畅的展开/收缩动画，无闪烁的用户体验

### v0.3.0 (2025-08-21)
- ✅ 通用YAML配置支持，移除schema依赖
- ✅ 智能界面生成系统
- ✅ 智能保存逻辑，同时生成YAML和DHF conf文件

### v0.2.1 (2025-08-20)
- ✅ 完全解决中文字体显示问题
- ✅ 实现100%中文本地化界面

### v0.2.0 (2025-08-20)
- ✅ 现代化GUI布局重大改进
- ✅ 双面板设计和工具栏优化

### v0.1.0 (2025-01-19)
- ✅ 项目初始化和基础架构
- ✅ GUI界面框架完成
- ✅ 基于真实conf文件的完整配置支持
- ✅ YAML schema驱动的动态UI系统
- ✅ 标准DHF conf格式输出
- ✅ 命令行版本实现
- 🚧 中文显示问题待解决

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 联系方式

- 项目地址: [GitHub Repository]
- 问题反馈: [GitHub Issues]
- 技术支持: [Email]

---

*DHF Configuration Manager - 让固件配置变得简单* 🎧