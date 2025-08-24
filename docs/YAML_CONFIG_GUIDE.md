# YAML配置文件详细指南

本文档详细说明ConfigCraft中YAML配置文件的格式、解析规则，以及如何转换为配置文件格式。

## 📋 目录
- [YAML配置文件结构](#yaml配置文件结构)
- [配置项命名规则](#配置项命名规则)
- [YAML到Conf的映射规则](#yaml到conf的映射规则)
- [配置分组详解](#配置分组详解)
- [示例与对照](#示例与对照)
- [手动维护指南](#手动维护指南)

## YAML配置文件结构

### 基础结构
```yaml
values:
  section.subsection.field_name: value
  section.field_name: value
  simple_field: value
```

### 实际示例
```yaml
values:
  # 基础配置
  basic.ic_model: 1
  basic.pa_control: false
  basic.vm_operation: "erase"
  basic.low_power_warn_time: 600000
  
  # 按键配置
  key_actions.call_scenario.active_click: "APP_MSG_NULL"
  key_actions.music_scenario.tws_connected_phone_connected_click_l: "APP_MSG_VOL_DOWN"
  
  # LED配置
  led_config.system_events.power_on: "LED_BLUE_ON"
  led_config.connection_status.bluetooth_connected: "LED_BLUE_ON"
  
  # 工厂设置
  factory.dut_mode_enable: false
  factory.reset_mode: "FACTORY_RESET_ALL"
  
  # 高级设置
  advanced.auto_power_on: true
  advanced.noise_cancellation: true
```

## 配置项命名规则

### 命名层级结构
```
section.group.field_name
└─────┘ └───┘ └────────┘
   │      │       │
   │      │       └─ 具体字段名
   │      └─────────── 子分组（可选）
   └────────────────── 主分组
```

### 支持的分组类型

| 主分组 | 子分组 | 说明 |
|--------|--------|------|
| `basic` | - | 基础配置项 |
| `key_actions` | `call_scenario`, `music_scenario` | 按键配置 |
| `led_config` | `system_events`, `connection_status`, `call_events` | LED灯效 |
| `factory` | - | 工厂设置 |
| `advanced` | - | 高级设置 |

## YAML到Conf的映射规则

### 转换规则说明

1. **键名转换**：所有点号(`.`)替换为下划线(`_`)
2. **大写转换**：整个键名转换为大写
3. **前缀添加**：添加`_`前缀
4. **值处理**：根据数据类型进行相应处理

### 映射算法
```go
// YAML: section.group.field_name = value
// Conf: _SECTION_GROUP_FIELD_NAME=value

func yamlToConfKey(yamlKey string) string {
    // 1. 替换点号为下划线
    confKey := strings.ReplaceAll(yamlKey, ".", "_")
    // 2. 转换为大写
    confKey = strings.ToUpper(confKey)
    // 3. 添加前缀
    confKey = "_" + confKey
    return confKey
}
```

### 实际映射示例

| YAML键名 | Conf键名 | 示例值 |
|----------|----------|--------|
| `basic.ic_model` | `_BASIC_IC_MODEL` | `1` |
| `basic.pa_control` | `_BASIC_PA_CONTROL` | `false` |
| `key_actions.call_scenario.active_click` | `_KEY_ACTIONS_CALL_SCENARIO_ACTIVE_CLICK` | `APP_MSG_NULL` |
| `led_config.system_events.power_on` | `_LED_CONFIG_SYSTEM_EVENTS_POWER_ON` | `LED_BLUE_ON` |

## 配置分组详解

### 1. Basic (基础配置)
```yaml
values:
  basic.ic_model: 1                    # IC型号: 0=AC7106A8, 1=AC7103D4
  basic.vm_operation: "erase"          # VM操作: "erase"或"keep"
  basic.pa_control: false              # 功放控制: true/false
  basic.low_power_warn_time: 600000    # 低电量警告时间(毫秒)
```

**生成的Conf**:
```conf
_BASIC_IC_MODEL=1
_BASIC_VM_OPERATION=erase
_BASIC_PA_CONTROL=false
_BASIC_LOW_POWER_WARN_TIME=600000
```

### 2. Key Actions (按键配置)

#### 通话场景 (call_scenario)
```yaml
values:
  # 通话中按键
  key_actions.call_scenario.active_click: "APP_MSG_NULL"
  key_actions.call_scenario.active_double_click: "APP_MSG_CALL_ANSWER"
  key_actions.call_scenario.active_long: "APP_MSG_NULL"
  
  # 来电按键
  key_actions.call_scenario.incoming_click: "APP_MSG_NULL"
  key_actions.call_scenario.incoming_double_click: "APP_MSG_CALL_ANSWER"
  key_actions.call_scenario.incoming_long: "APP_MSG_CALL_HANGUP"
  key_actions.call_scenario.incoming_triple_click: "APP_MSG_NULL"
  
  # 外拨按键
  key_actions.call_scenario.outgoing_click: "APP_MSG_NULL"
  
  # Siri按键
  key_actions.call_scenario.siri_click: "APP_MSG_NULL"
  key_actions.call_scenario.siri_double_click: "APP_MSG_NULL"
```

#### 音乐场景 (music_scenario)
```yaml
values:
  # TWS连接+手机连接状态下的按键
  key_actions.music_scenario.tws_connected_phone_connected_click_l: "APP_MSG_VOL_DOWN"
  key_actions.music_scenario.tws_connected_phone_connected_click_r: "APP_MSG_VOL_UP"
  key_actions.music_scenario.tws_connected_phone_connected_double_click_l: "APP_MSG_MUSIC_PP"
  key_actions.music_scenario.tws_connected_phone_connected_double_click_r: "APP_MSG_MUSIC_PP"
  key_actions.music_scenario.tws_connected_phone_connected_long_l: "APP_MSG_OPEN_SIRI"
  key_actions.music_scenario.tws_connected_phone_connected_long_r: "APP_MSG_OPEN_SIRI"
  key_actions.music_scenario.tws_connected_phone_connected_triple_click_l: "APP_MSG_MUSIC_PREV"
  key_actions.music_scenario.tws_connected_phone_connected_triple_click_r: "APP_MSG_MUSIC_NEXT"
  
  # TWS断开+手机连接状态下的按键
  key_actions.music_scenario.tws_disconnected_phone_connected_click_l: "APP_MSG_VOL_DOWN"
  key_actions.music_scenario.tws_disconnected_phone_connected_click_r: "APP_MSG_VOL_UP"
  # ... 其他tws_disconnected配置
```

**按键值含义**:
| 值 | 功能 |
|----|------|
| `APP_MSG_NULL` | 无操作 |
| `APP_MSG_CALL_ANSWER` | 接听电话 |
| `APP_MSG_CALL_HANGUP` | 挂断电话 |
| `APP_MSG_VOL_UP` | 音量+ |
| `APP_MSG_VOL_DOWN` | 音量- |
| `APP_MSG_MUSIC_PP` | 播放/暂停 |
| `APP_MSG_MUSIC_PREV` | 上一首 |
| `APP_MSG_MUSIC_NEXT` | 下一首 |
| `APP_MSG_OPEN_SIRI` | 打开Siri |

### 3. LED Config (LED灯效配置)

#### 系统事件 (system_events)
```yaml
values:
  led_config.system_events.power_on: "LED_BLUE_ON"          # 开机
  led_config.system_events.power_off: "LED_OFF"            # 关机
  led_config.system_events.charging: "LED_RED_ON"          # 充电中
  led_config.system_events.charge_complete: "LED_GREEN_ON" # 充电完成
  led_config.system_events.low_battery: "LED_RED_FLASH"    # 低电量
```

#### 连接状态 (connection_status)
```yaml
values:
  led_config.connection_status.bluetooth_connected: "LED_BLUE_ON"      # 蓝牙已连接
  led_config.connection_status.bluetooth_disconnected: "LED_RED_SLOW"  # 蓝牙断开
  led_config.connection_status.pairing_mode: "LED_BLUE_FAST"          # 配对模式
  led_config.connection_status.tws_connected: "LED_GREEN_ON"          # TWS已连接
  led_config.connection_status.tws_disconnected: "LED_ORANGE_SLOW"     # TWS断开
```

#### 通话事件 (call_events)
```yaml
values:
  led_config.call_events.incoming_call: "LED_BLUE_FAST"  # 来电
  led_config.call_events.active_call: "LED_BLUE_ON"     # 通话中
  led_config.call_events.call_end: "LED_OFF"            # 通话结束
```

**LED灯效值含义**:
| 值 | 效果 |
|----|------|
| `LED_OFF` | 关闭 |
| `LED_BLUE_ON` | 蓝灯常亮 |
| `LED_RED_ON` | 红灯常亮 |
| `LED_GREEN_ON` | 绿灯常亮 |
| `LED_BLUE_FAST` | 蓝灯快闪 |
| `LED_RED_SLOW` | 红灯慢闪 |
| `LED_RED_FLASH` | 红灯闪烁 |
| `LED_ORANGE_SLOW` | 橙灯慢闪 |

### 4. Factory (工厂设置)
```yaml
values:
  factory.dut_mode_enable: false              # DUT测试模式: true/false
  factory.reset_mode: "FACTORY_RESET_ALL"     # 重置模式
  factory.test_mode_timeout: 30               # 测试模式超时(秒)
```

### 5. Advanced (高级设置)
```yaml
values:
  advanced.auto_power_on: true           # 自动开机: true/false
  advanced.low_latency_mode: false       # 低延迟模式: true/false
  advanced.noise_cancellation: true      # 降噪: true/false
  advanced.three_call_enable: true       # 三方通话: true/false
```

## 示例与对照

### 完整YAML示例
```yaml
values:
  # === 基础配置 ===
  basic.ic_model: 1
  basic.pa_control: false
  basic.vm_operation: "erase"
  basic.low_power_warn_time: 600000
  
  # === 按键配置 ===
  # 通话场景
  key_actions.call_scenario.active_click: "APP_MSG_NULL"
  key_actions.call_scenario.incoming_double_click: "APP_MSG_CALL_ANSWER"
  key_actions.call_scenario.incoming_long: "APP_MSG_CALL_HANGUP"
  
  # 音乐场景
  key_actions.music_scenario.tws_connected_phone_connected_click_l: "APP_MSG_VOL_DOWN"
  key_actions.music_scenario.tws_connected_phone_connected_click_r: "APP_MSG_VOL_UP"
  
  # === LED配置 ===
  led_config.system_events.power_on: "LED_BLUE_ON"
  led_config.connection_status.bluetooth_connected: "LED_BLUE_ON"
  led_config.call_events.incoming_call: "LED_BLUE_FAST"
  
  # === 工厂设置 ===
  factory.dut_mode_enable: false
  factory.reset_mode: "FACTORY_RESET_ALL"
  
  # === 高级设置 ===
  advanced.auto_power_on: true
  advanced.noise_cancellation: true
```

### 对应的Conf文件
```conf
#***************************************************************************
#                       Configuration Settings
#***************************************************************************

# 基础配置 (Basic Configuration)
#------------------------------------
_BASIC_IC_MODEL=1
_BASIC_PA_CONTROL=false
_BASIC_VM_OPERATION=erase
_BASIC_LOW_POWER_WARN_TIME=600000

# 按键配置 (Key Actions)
#----------------------------
_KEY_ACTIONS_CALL_SCENARIO_ACTIVE_CLICK=APP_MSG_NULL
_KEY_ACTIONS_CALL_SCENARIO_INCOMING_DOUBLE_CLICK=APP_MSG_CALL_ANSWER
_KEY_ACTIONS_CALL_SCENARIO_INCOMING_LONG=APP_MSG_CALL_HANGUP
_KEY_ACTIONS_MUSIC_SCENARIO_TWS_CONNECTED_PHONE_CONNECTED_CLICK_L=APP_MSG_VOL_DOWN
_KEY_ACTIONS_MUSIC_SCENARIO_TWS_CONNECTED_PHONE_CONNECTED_CLICK_R=APP_MSG_VOL_UP

# LED配置 (LED Configuration)
#-------------------------------
_LED_CONFIG_SYSTEM_EVENTS_POWER_ON=LED_BLUE_ON
_LED_CONFIG_CONNECTION_STATUS_BLUETOOTH_CONNECTED=LED_BLUE_ON
_LED_CONFIG_CALL_EVENTS_INCOMING_CALL=LED_BLUE_FAST

# 工厂设置 (Factory Settings)
#---------------------------------
_FACTORY_DUT_MODE_ENABLE=false
_FACTORY_RESET_MODE=FACTORY_RESET_ALL

# 高级设置 (Advanced Settings)
#----------------------------------
_ADVANCED_AUTO_POWER_ON=true
_ADVANCED_NOISE_CANCELLATION=true
```

## 手动维护指南

### 1. 添加新配置项
```yaml
# 在对应的分组下添加
values:
  section.new_field: value
  # 或者带子分组
  section.subsection.new_field: value
```

### 2. 修改现有配置
```yaml
# 直接修改对应的值
values:
  basic.ic_model: 0  # 从1改为0
  led_config.system_events.power_on: "LED_OFF"  # 从LED_BLUE_ON改为LED_OFF
```

### 3. 删除配置项
```yaml
# 直接删除对应行即可
# 删除: basic.pa_control: false
```

### 4. 注意事项

**键名规则**:
- 使用小写字母和下划线
- 分层结构用点号(`.`)分隔
- 不要使用空格或特殊字符

**值的格式**:
- 字符串值用引号包围: `"APP_MSG_NULL"`
- 布尔值直接写: `true` / `false`
- 数字值直接写: `600000`

**常见错误**:
```yaml
# ❌ 错误示例
values:
  basic.ic model: 1              # 键名不能有空格
  basic.pa_control: "false"      # 布尔值不要加引号
  LED_CONFIG.POWER_ON: value     # 不要用大写和下划线

# ✅ 正确示例
values:
  basic.ic_model: 1
  basic.pa_control: false
  led_config.system_events.power_on: "LED_BLUE_ON"
```

### 5. 验证配置
使用ConfigCraft加载YAML文件，查看：
1. 是否正确识别所有配置项
2. 分组是否按预期显示
3. 生成的conf文件是否格式正确

这样就能确保手动维护的YAML文件与工具完全兼容！