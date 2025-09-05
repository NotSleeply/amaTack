// 定义游戏的核心数据结构和基本操作。
package amazon

import (
	"fmt"
	"sync"

	"github.com/tongque0/gotack"
)

const (
	Empty = iota // 值为0，表示空位  使用 iota 自动递增生成常量值：
	Black        // 值为1，表示黑方棋子
	White        // 值为2，表示白方棋子
	Arrow        // 值为3，表示障碍
)

// 方向数组 左上、上、右上、右、右下、下、左下、左
var (
	dir = [8][2]int{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
		{1, 0},
		{1, -1},
		{0, -1},
	}
)

// 棋子的位置
type Position struct {
	X int
	Y int
}

// 一次完整的移动，包括起点、终点和放置箭的位置。
type AmazonMove struct {
	From Position
	To   Position
	Put  Position
}

// 定义了一个10x10的棋盘，使用二维数组表示棋盘状态。
type AmazonBoard [10][10]int

// 初始化新的棋盘。
func NewBoard() *AmazonBoard {
	board := &AmazonBoard{}

	// 初始化棋盘，设置棋子的初始位置
	// 白棋 位置
	board[0][3] = White
	board[0][6] = White
	board[3][0] = White
	board[3][9] = White
	// 黑棋 位置
	board[6][0] = Black
	board[6][9] = Black
	board[9][3] = Black
	board[9][6] = Black

	return board
}

/*
* 打印棋盘状态
* 空位显示为 "."
* 黑棋显示为 "B"
* 白棋显示为 "W"
* 障碍显示为 "X"
 */
func (b *AmazonBoard) Print() {
	for i := 0; i < len(*b); i++ {
		for j := 0; j < len((*b)[i]); j++ {
			switch (*b)[i][j] {
			case Empty:
				fmt.Print(". ")
			case Black:
				fmt.Print("B ")
			case White:
				fmt.Print("W ")
			case Arrow:
				fmt.Print("X ")
			}
		}
		fmt.Println()
	}
}

// 打印移动信息
func (m AmazonMove) String() string {
	return fmt.Sprintf("From (%d,%d)\tTo (%d,%d)\tPut (%d,%d)\n", m.From.X, m.From.Y, m.To.X, m.To.Y, m.Put.X, m.Put.Y)
}

// 步法棋盘
func (b *AmazonBoard) PrintMoveBoard() {
	for i := 0; i < len(*b); i++ {
		for j := 0; j < len((*b)[i]); j++ {
			if (*b)[i][j] == 100 { // 假设使用100表示不可达的位置
				fmt.Print(" . ") // 注意这里有两个空格，与下面两位数的步数占位保持一致
			} else {
				// 如果步数小于10，则在前面添加一个空格来保持对齐
				if (*b)[i][j] < 10 {
					fmt.Printf(" %d ", (*b)[i][j])
				} else {
					fmt.Printf("%d ", (*b)[i][j])
				}
			}
		}
		fmt.Println()
	}
}

/*
* 生成所有合法移动
* 根据当前玩家确定棋子颜色
* 找出该颜色的所有棋子
* 为每个棋子并发生成所有可能的移动
* 收集并返回所有合法移动
 */
func (b *AmazonBoard) GetAllMoves(IsMaxPlayer bool) []gotack.Move {
	var moves []gotack.Move
	var color = 1
	if !IsMaxPlayer {
		color = 2
	}
	// 获取指定颜色的所有棋子位置
	allChess := b.getAllChess(color)

	// 创建带缓冲通道，用于存储生成的合法移动
	moveChan := make(chan AmazonMove, len(allChess))

	// 使用等待组，以便等待所有并发协程完成
	var wg sync.WaitGroup

	// 遍历所有棋子  核心
	for _, chess := range allChess {
		// 启动并发协程，为每个棋子生成合法移动
		wg.Add(1)
		go func(chess Position) {
			defer wg.Done()
			b.generateMovesForChess(chess, moveChan)
		}(chess)
	}

	// 在单独的协程中等待所有并发协程完成，并关闭通道
	go func() {
		wg.Wait()
		close(moveChan)
	}()

	// 从通道中读取所有生成的合法移动，并添加到结果切片中
	for move := range moveChan {
		moves = append(moves, move)
	}

	return moves
}

// 执行移动操作
func (b *AmazonBoard) Move(move gotack.Move) {
	m, ok := move.(AmazonMove)
	if !ok {
		fmt.Println("Invalid move type")
		return
	}

	b[m.To.X][m.To.Y] = b[m.From.X][m.From.Y] // 移动棋子
	b[m.From.X][m.From.Y] = Empty             // 清空原位置
	b[m.Put.X][m.Put.Y] = Arrow               // 放置障碍
}

// 撤销移动操作
func (b *AmazonBoard) UndoMove(move gotack.Move) {
	m, ok := move.(AmazonMove)
	if !ok {
		fmt.Println("Invalid move type")
		return
	}
	b[m.From.X][m.From.Y] = b[m.To.X][m.To.Y] // 恢复棋子位置
	b[m.To.X][m.To.Y] = Empty                 // 清空移动后位置
	if b[m.From.X] != b[m.Put.X] || b[m.From.Y] != b[m.Put.Y] {
		b[m.Put.X][m.Put.Y] = Empty // 移除障碍
	}
}

// 检查游戏是否结束
func (b *AmazonBoard) IsGameOver() bool {
	// 检查黑方是否还有合法的移动
	blackMoves := b.GetAllMoves(true) // 假设true代表黑方
	if len(blackMoves) == 0 {
		return true // 黑方没有合法的移动，游戏结束
	}

	// 检查白方是否还有合法的移动
	whiteMoves := b.GetAllMoves(false) // 假设false代表白方

	// 直接返回白方移动列表长度是否为0的结果
	return len(whiteMoves) == 0
}

// 检查位置是否合法
func (b *AmazonBoard) legal(x, y int) bool {
	return x >= 0 && y >= 0 && x < 10 && y < 10
}

// 获取指定颜色的所有棋子位置
func (b *AmazonBoard) getAllChess(color int) []Position {
	var positions []Position
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if (*b)[i][j] == color {
				positions = append(positions, Position{i, j})
			}
		}
	}
	return positions
}

// 为单个棋子生成所有合法移动
func (b *AmazonBoard) generateMovesForChess(chess Position, moveChan chan AmazonMove) {
	// 遍历所有方向
	for j := 0; j < 8; j++ {
		// 初始方向
		x, y := chess.X+dir[j][0], chess.Y+dir[j][1]
		// 沿着当前方向一直移动，直到碰到边界或非空位置
		for b.legal(x, y) && (*b)[x][y] == Empty {
			// 从当前位置，遍历8个方向放置障碍箭
			for k := 0; k < 8; k++ {
				ax, ay := x+dir[k][0], y+dir[k][1]
				// 沿着当前方向一直移动，寻找可放置箭的位置
				for b.legal(ax, ay) && ((*b)[ax][ay] == Empty || ax == chess.X && ay == chess.Y) {
					// 创建合法移动对象，并发送到通道中
					move := AmazonMove{
						From: Position{X: chess.X, Y: chess.Y},
						To:   Position{X: x, Y: y},
						Put:  Position{X: ax, Y: ay},
					}
					moveChan <- move
					// 继续沿着当前方向
					ax += dir[k][0]
					ay += dir[k][1]
				}
			}
			// 继续沿原方向移动棋子
			x += dir[j][0]
			y += dir[j][1]
		}
	}
}

// 评估函数接口
func (b *AmazonBoard) EvaluateFunc(opts gotack.EvalOptions) float64 {
	return EvaluateFunc(&opts)
}

// 计算当前棋盘的哈希值 用于检测重复局面。
func (b *AmazonBoard) Hash() uint64 {
	var hash uint64 = 0
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			piece := b[i][j]
			hash ^= zobristTable[i][j][piece]
		}
	}
	return hash
}

// 克隆棋盘
func (b *AmazonBoard) Clone() gotack.Board {
	clone := AmazonBoard{} // 创建一个新的 AmazonBoard 实例

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			clone[i][j] = b[i][j] // 复制每个位置的棋子状态
		}
	}

	return &clone // 返回克隆的棋盘指针
}
