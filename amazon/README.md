# gotack-examples-amazon

## 亚马逊棋程序

- bin目录下有四个可执行程序
- Qtack2.0为在gotack1.2架构下的新版亚马逊棋，快速版，搜索深度为2跳3
- Stack2.0为在gotack1.2架构下的新版亚马逊棋，慢速版本，搜索深度为2跳4
- ZhangT1.0为学长版本亚马逊棋
- ZhangT1.0为学长版本亚马逊棋

## 启动

- 安装模块 `go get github.com/tongque0/gotack`
- 启动 `make run`
- 构建 `make build`

>此项启动用于调试，正常使用需要build后，载入棋盘UI中。

## 了解项目结构

- main.go - 程序入口，处理用户输入和游戏流程；根据UI中的通信引擎协议所编写
- amazon.go - 定义亚马逊棋的核心数据结构和基本操作
- value.go - 实现评估函数逻辑
- record.go - 处理游戏记录
- Zobrist.go - 实现 Zobrist 哈希（用于棋盘状态比较）
