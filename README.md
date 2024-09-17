
# gotack-examples

欢迎来到 gotack 包的示例仓库！本仓库包含了辽宁科技大学部分棋种计算机博弈的实现示例，适用于研究和教学目的。

## 项目概述

gotack-examples 旨在提供一个参考框架，帮助研究人员和开发者快速入门棋类博弈算法的开发和实验。目前，本仓库涵盖以下棋种的实现：

- 点格棋（采用机器学习+UCT算法）
- 亚马逊棋（采用简单的 Alpha-Beta 剪枝算法）
- 不围棋（基于价值评估的递归算法）

## 快速启动

如果没有安装 CMake，可以手动进入 `ui` 目录并启动可执行文件 `SAU Game Platform.exe`：

### 亚马逊棋

1. 克隆仓库到本地：
   ```
   git clone https://example.com/gotack-examples.git
   cd gotack-examples
   ```
2. 编译并启动 UI：
   ```
   make ui
   ```
3. 加载 `amazon/bin` 目录下的引擎开始博弈。

### 不围棋

1. 克隆仓库到本地：
   ```
   git clone https://example.com/gotack-examples.git
   cd gotack-examples
   ```
2. 编译并启动 UI：
   ```
   make ui
   ```
3. 加载 `nogo/bin` 目录下的引擎开始博弈。

## 技术栈

本项目主要使用以下技术实现：

- 机器学习
- UCT算法
- Alpha-Beta 剪枝
- 递归算法


## 许可

该项目采用 [MIT 许可证](LICENSE)。

# UI

- UI文件中为棋谱
