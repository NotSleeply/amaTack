package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	board       [9][9]int            // 9x9的围棋棋盘，0表示空位，1表示白棋，-1表示黑棋
	dfsAirVisit [9][9]bool           // DFS访问标记数组，用于判断棋子是否有气
	cx          = []int{-1, 0, 1, 0} // 上下左右四个方向的x坐标偏移
	cy          = []int{0, -1, 0, 1} // 上下左右四个方向的y坐标偏移
)
var (
	line        string // 输入的命令行
	step        int    // 当前步数
	IsMaxPlayer bool   // 是否为先手玩家
)

// Position 表示棋盘上的一个位置
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// GameData 存储游戏数据，包括请求和响应的位置
type GameData struct {
	Requests  []Position `json:"requests"`
	Responses []Position `json:"responses"`
}

// inBorder 检查坐标(x,y)是否在棋盘范围内
func inBorder(x, y int) bool {
	return x >= 0 && y >= 0 && x < 9 && y < 9
}

// dfsAir 使用深度优先搜索判断从位置(fx,fy)开始的连通棋子组是否有气
// 返回true表示有气，false表示没气（被围）
func dfsAir(fx, fy int) bool {
	dfsAirVisit[fx][fy] = true // 标记当前位置已访问
	flag := false              // 是否找到气的标志

	// 检查四个方向
	for dir := 0; dir < 4; dir++ {
		dx := fx + cx[dir]
		dy := fy + cy[dir]
		if inBorder(dx, dy) {
			// 如果相邻位置是空位，则有气
			if board[dx][dy] == 0 {
				flag = true
			}
			// 如果相邻位置是同色棋子且未访问过，递归检查
			if board[dx][dy] == board[fx][fy] && !dfsAirVisit[dx][dy] {
				if dfsAir(dx, dy) {
					flag = true
				}
			}
		}
	}
	return flag
}

// judgeAvailable 判断在位置(fx,fy)放置颜色为col的棋子是否合法
// col: 1表示白棋，-1表示黑棋
func judgeAvailable(fx, fy, col int) bool {
	// 如果位置已被占用，不能下棋
	if board[fx][fy] != 0 {
		return false
	}

	// 临时放置棋子
	board[fx][fy] = col
	clearVisit()

	// 检查放置后自己是否有气
	if !dfsAir(fx, fy) {
		board[fx][fy] = 0 // 恢复棋盘状态
		return false      // 自杀式下法，不合法
	}

	// 检查是否能吃掉对方的棋子
	for dir := 0; dir < 4; dir++ {
		dx := fx + cx[dir]
		dy := fy + cy[dir]
		if inBorder(dx, dy) {
			// 如果相邻位置有对方棋子且未检查过
			if board[dx][dy] != 0 && !dfsAirVisit[dx][dy] {
				// 如果对方棋子没气，则可以吃掉
				if !dfsAir(dx, dy) {
					board[fx][fy] = 0 // 恢复棋盘状态
					return false      // 这种情况下不合法（可能是特殊规则）
				}
			}
		}
	}

	board[fx][fy] = 0 // 恢复棋盘状态
	return true
}

// clearVisit 清空DFS访问标记数组
func clearVisit() {
	for i := range dfsAirVisit {
		for j := range dfsAirVisit[i] {
			dfsAirVisit[i][j] = false
		}
	}
}

// valuePoint 评估在位置(x,y)下棋的价值
// 通过计算下棋后能限制对方多少个可下位置来评估价值
func valuePoint(x, y int) int {
	value := 0

	// 尝试AI（黑棋-1）在此位置下棋
	if judgeAvailable(x, y, -1) {
		board[x][y] = -1
		// 计算对方（白棋1）还有多少合法位置
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if board[i][j] == 0 {
					if !judgeAvailable(i, j, 1) {
						value++ // 限制了对方一个位置
					}
				}
			}
		}
	}

	// 尝试对方（白棋1）在此位置下棋
	if judgeAvailable(x, y, 1) {
		board[x][y] = 1
		// 计算AI（黑棋-1）还有多少合法位置
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if board[i][j] == 0 {
					if !judgeAvailable(i, j, -1) {
						value++ // 限制了AI一个位置
					}
				}
			}
		}
	}

	board[x][y] = 0 // 恢复棋盘状态
	return value
}

// findMaxValuePoint 从可用位置列表中找到价值最高的位置
// 返回价值最高的位置列表和最高价值
func findMaxValuePoint(availableList []int) ([]int, int) {
	maxValue := -1
	var waitList []int

	for _, p := range availableList {
		x, y := p/9, p%9          // 将一维坐标转换为二维坐标
		value := valuePoint(x, y) // 计算该位置的价值

		if value > maxValue {
			maxValue = value
			waitList = []int{p} // 重新开始等待列表
		} else if value == maxValue {
			waitList = append(waitList, p) // 添加到等待列表
		}
	}
	return waitList, maxValue
}

// getMaxValuePoint 从等价值的位置中选择最佳位置
// 使用更深层的分析来选择最优位置
func getMaxValuePoint(waitList []int) int {
	if len(waitList) == 0 {
		return -1
	}
	if len(waitList) == 1 {
		return waitList[0]
	}

	// 用于存储最佳选择的结构
	var result struct {
		Index int // 最佳位置的索引
		Value int // 最佳位置的价值
	}

	// 对每个候选位置进行更深层分析
	for _, p := range waitList {
		x, y := p/9, p%9
		board[x][y] = -1 // 临时放置AI棋子

		// 重新评估在此位置下棋后的最大价值
		_, maxValue := findMaxValuePoint(waitList)
		if maxValue > result.Value {
			result.Index = p
			result.Value = maxValue
		}

		board[x][y] = 0 // 恢复棋盘状态
	}
	return result.Index
}

// abs 返回整数的绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// getAllPossibleMoves 获取指定颜色棋子的所有可能下棋位置
// col: 1表示白棋，-1表示黑棋
func getAllPossibleMoves(col int) []int {
	var moves []int
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if judgeAvailable(i, j, col) {
				moves = append(moves, i*9+j) // 将二维坐标转换为一维索引
			}
		}
	}
	return moves
}

// aiMove AI进行一步棋的决策和执行
func aiMove() {
	moves := getAllPossibleMoves(-1) // 获取AI（黑棋）所有可能的下棋位置
	if len(moves) == 0 {
		fmt.Println("No moves available for AI.")
		return
	}

	// 找到价值最高的位置列表
	waitList, _ := findMaxValuePoint(moves)
	if len(waitList) == 0 {
		fmt.Println("AI cannot find a valid move.")
		return
	}

	// 从等价值位置中选择最佳位置
	bestMove := getMaxValuePoint(waitList)
	if bestMove == -1 {
		fmt.Println("AI cannot decide the best move.")
		return
	}

	// 执行下棋
	x, y := bestMove/9, bestMove%9
	board[x][y] = -1                        // 在选定位置放置AI棋子
	fmt.Printf("move %c%c\n", 'A'+x, 'A'+y) // 输出下棋位置（用字母表示）
}

// main 主函数，处理游戏命令和逻辑
func main() {
	sc := bufio.NewScanner(os.Stdin)

	// 循环读取命令
	for sc.Scan() {
		line = sc.Text()
		line = strings.TrimSpace(line)

		if line == "name?" {
			// 返回AI名称
			fmt.Println("name Tack")
		} else if line == "quit" {
			// 退出游戏
			fmt.Println("Quitting game.")
			os.Exit(0)
		} else if strings.HasPrefix(line, "new") {
			// 开始新游戏
			resetBoard()
			args := strings.Split(line, " ")
			if len(args) > 1 && args[1] == "black" {
				// AI执黑棋，先手
				IsMaxPlayer = true
				aiMove() // AI先下
			} else {
				// AI执白棋，后手
				IsMaxPlayer = false
			}
			step = 1
		} else if strings.HasPrefix(line, "move") {
			// 处理对方下棋
			words := strings.Split(line, " ")
			move := words[1]
			X := move[0] - 'A' // 将字母坐标转换为数字
			Y := move[1] - 'A'
			board[X][Y] = 1 // 在棋盘上放置对方棋子（白棋）
			step++
			aiMove() // AI应对
		} else if line == "end" {
			// 游戏结束
			fmt.Println("Game over.")
			resetBoard()
			continue
		} else {
			// 未知命令
			fmt.Println("Unknown command. Available commands: 'new [color]', 'move x y', 'end', 'quit'.")
		}
	}

	// 处理输入错误
	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "reading standard input: %s\n", err)
	}
}

// resetBoard 重置棋盘，将所有位置清空
func resetBoard() {
	for i := range board {
		for j := range board[i] {
			board[i][j] = 0 // 0表示空位
		}
	}
	clearVisit() // 清空访问标记
}
