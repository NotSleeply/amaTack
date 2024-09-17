package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"tamazon/amazon"

	"github.com/tongque0/gotack"
)

const INF = 0x3f3f3f3f

var (
	line  string
	step  int
	board *amazon.AmazonBoard
	color int
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line = sc.Text()
		if line == "name?" {
			fmt.Println("name MTack")
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
func runSearch() {
	var IsMaxPlayer = true
	var e *gotack.Evaluator
	if color == 2 {
		IsMaxPlayer = false
	}
	// 根据步数动态设置时间限制 调参
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
		gotack.AlphaBeta,
		gotack.NewEvaluatorOptions(
			gotack.WithBoard(board),
			gotack.WithDepth(searchDetph),
			gotack.WithIsMaxPlayer(IsMaxPlayer),
			gotack.WithStep(step),
			gotack.WithIsDetail(true),
		),
	)
	move := e.GetBestMove()
	m, ok := move[0].(amazon.AmazonMove)
	if !ok {
		return
	}
	board.Move(move[0])
	fmt.Printf("move %c%c%c%c%c%c\n", m.From.Y+'A', m.From.X+'A', m.To.Y+'A', m.To.X+'A', m.Put.Y+'A', m.Put.X+'A')
	amazon.AddRecord(m.From.Y+'a', 10-m.From.X, m.To.Y+'a', 10-m.To.X, m.Put.Y+'a', 10-m.Put.X)
	step++
}
