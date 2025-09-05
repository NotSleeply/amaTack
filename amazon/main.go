package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"tamazon/amazon"

	"github.com/tongque0/gotack"
)

const INF = 0x3f3f3f3f  // 表示"无穷大"的常量，常用于最大最小值初始化
const Name = "MTackTao" // 程序名称

var (
	line  string              // 存储输入的行
	step  int                 // 当前步数
	board *amazon.AmazonBoard // 棋盘
	color int                 // 当前颜色
)

/*
 * main
 * 通过命令行输入实现前端UI交互协议
 * 输入"new black"或"new white"开始新游戏
 * 输入"move A1B2C3"进行移动，格式为"move from to put"
 * 输入"end"保存游戏记录
 */
func main() {
	fmt.Printf("-------------欢迎使用%s-----------------\n", Name)
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line = sc.Text()
		if line == "name?" {
			fmt.Printf("name %s\n", Name)
		} else if line == "quit" {
			os.Exit(0)
		} else if strings.HasPrefix(line, "new") {
			step = 1
			words := strings.Split(line, " ")
			board = amazon.NewBoard()
			if words[1] == "black" {
				color = amazon.Black
				runSearch()
			} else {
				color = amazon.White
			}
		} else if strings.HasPrefix(line, "move") {
			words := strings.Split(line, " ")
			move := words[1]
			board[move[3]-'A'][move[2]-'A'] = board[move[1]-'A'][move[0]-'A']
			board[move[1]-'A'][move[0]-'A'] = amazon.Empty
			board[move[5]-'A'][move[4]-'A'] = amazon.Arrow
			step++
			if !board.IsGameOver() {
				runSearch()
			}
		} else if line == "end" {
			amazon.Save()
			continue
		} else {
			amazon.Save()
			continue
		}
	}
}

/*
 * runSearch
 * 运行搜索算法，寻找最佳移动
 */

func runSearch() {
	var IsMaxPlayer = true
	var e *gotack.Evaluator
	if color == 2 { // 白方
		IsMaxPlayer = false
	}
	// 根据步数动态设置搜索深度
	var searchDetph int
	switch {
	case step < 23:
		searchDetph = 2
	case step < 50:
		searchDetph = 3
	case step < 70:
		searchDetph = 4
	default:
		searchDetph = 5
	}

	// 创建评估器
	e = gotack.NewEvaluator(
		gotack.AlphaBeta, // 使用Alpha-Beta剪枝算法
		gotack.NewEvaluatorOptions(
			gotack.WithBoard(board),             // 当前棋盘
			gotack.WithDepth(searchDetph),       // 搜索深度
			gotack.WithIsMaxPlayer(IsMaxPlayer), //最大玩家
			gotack.WithStep(step),               // 当前步数
			gotack.WithIsDetail(true),           // 详细输出
		),
	)
	// 获取最佳移动
	move := e.GetBestMove()
	m, ok := move[0].(amazon.AmazonMove)
	if !ok {
		return
	}
	// 执行最佳移动
	board.Move(move[0])
	// 输出移动信息
	fmt.Printf("move %c%c%c%c%c%c\n", m.From.Y+'A', m.From.X+'A', m.To.Y+'A', m.To.X+'A', m.Put.Y+'A', m.Put.X+'A')
	// 记录游戏
	amazon.AddRecord(m.From.Y+'a', 10-m.From.X, m.To.Y+'a', 10-m.To.X, m.Put.Y+'a', 10-m.Put.X)
	// 更新步数
	step++
}
