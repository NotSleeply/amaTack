package main

import (
	"bufio"
	"fmt"
	"nogo/nogo"
	"os"
	"strings"
)

var (
	board       [9][9]int
	dfsAirVisit [9][9]bool
	cx          = []int{-1, 0, 1, 0}
	cy          = []int{0, -1, 0, 1}
)
var (
	line        string
	step        int
	IsMaxPlayer bool
)

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type GameData struct {
	Requests  []Position `json:"requests"`
	Responses []Position `json:"responses"`
}

func inBorder(x, y int) bool {
	return x >= 0 && y >= 0 && x < 9 && y < 9
}

func dfsAir(fx, fy int) bool {
	dfsAirVisit[fx][fy] = true
	flag := false
	for dir := 0; dir < 4; dir++ {
		dx := fx + cx[dir]
		dy := fy + cy[dir]
		if inBorder(dx, dy) {
			if board[dx][dy] == 0 {
				flag = true
			}
			if board[dx][dy] == board[fx][fy] && !dfsAirVisit[dx][dy] {
				if dfsAir(dx, dy) {
					flag = true
				}
			}
		}
	}
	return flag
}

func judgeAvailable(fx, fy, col int) bool {
	if board[fx][fy] != 0 {
		return false
	}
	board[fx][fy] = col
	clearVisit()
	if !dfsAir(fx, fy) {
		board[fx][fy] = 0
		return false
	}
	for dir := 0; dir < 4; dir++ {
		dx := fx + cx[dir]
		dy := fy + cy[dir]
		if inBorder(dx, dy) {
			if board[dx][dy] != 0 && !dfsAirVisit[dx][dy] {
				if !dfsAir(dx, dy) {
					board[fx][fy] = 0
					return false
				}
			}
		}
	}
	board[fx][fy] = 0
	return true
}

func clearVisit() {
	for i := range dfsAirVisit {
		for j := range dfsAirVisit[i] {
			dfsAirVisit[i][j] = false
		}
	}
}

func valuePoint(x, y int) int {
	value := 0
	if judgeAvailable(x, y, -1) {
		board[x][y] = -1
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if board[i][j] == 0 {
					if !judgeAvailable(i, j, 1) {
						value++
					}
				}
			}
		}
	}
	if judgeAvailable(x, y, 1) {
		board[x][y] = 1
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if board[i][j] == 0 {
					if !judgeAvailable(i, j, -1) {
						value++
					}
				}
			}
		}
	}
	board[x][y] = 0
	return value
}

func findMaxValuePoint(availableList []int) ([]int, int) {
	maxValue := -1
	var waitList []int
	for _, p := range availableList {
		x, y := p/9, p%9
		value := valuePoint(x, y)
		if value > maxValue {
			maxValue = value
			waitList = []int{p}
		} else if value == maxValue {
			waitList = append(waitList, p)
		}
	}
	return waitList, maxValue
}

func getScatterPoint(availableList []int) int {
	maxDis := -1
	result := -1
	for _, p := range availableList {
		x, y := p/9, p%9
		minDis := 100
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if board[i][j] == -1 {
					if d := abs(x-i) + abs(y-j); d < minDis {
						minDis = d
					}
				}
			}
		}
		if minDis > maxDis {
			maxDis = minDis
			result = p
		}
	}
	return result
}
func getMaxValuePoint(waitList []int) int {
	if len(waitList) == 0 {
		return -1
	}
	if len(waitList) == 1 {
		return waitList[0]
	}

	var result struct {
		Index int
		Value int
	}
	for _, p := range waitList {
		x, y := p/9, p%9
		board[x][y] = -1
		_, maxValue := findMaxValuePoint(waitList)
		if maxValue > result.Value {
			result.Index = p
			result.Value = maxValue
		}
		board[x][y] = 0
	}
	return result.Index
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func getAllPossibleMoves(col int) []int {
	var moves []int
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if judgeAvailable(i, j, col) {
				moves = append(moves, i*9+j)
			}
		}
	}
	return moves
}
func aiMove() {
	moves := getAllPossibleMoves(-1) // Assume AI uses -1
	if len(moves) == 0 {
		fmt.Println("No moves available for AI.")
		return
	}
	waitList, _ := findMaxValuePoint(moves)
	if len(waitList) == 0 {
		fmt.Println("AI cannot find a valid move.")
		return
	}
	bestMove := getMaxValuePoint(waitList)
	if bestMove == -1 {
		fmt.Println("AI cannot decide the best move.")
		return
	}
	x, y := bestMove/9, bestMove%9
	board[x][y] = -1 // Execute AI move
	fmt.Printf("move %c%c\n", 'A'+x, 'A'+y)
	nogo.AddMove(IsMaxPlayer, x, y)
}

func main() {
	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		line = sc.Text()
		line = strings.TrimSpace(line)
		if line == "name?" {
			fmt.Println("name Tack")
		} else if line == "quit" {
			fmt.Println("Quitting game.")
			os.Exit(0)
		} else if strings.HasPrefix(line, "new") {
			resetBoard()
			args := strings.Split(line, " ")
			if len(args) > 1 && args[1] == "black" {
				IsMaxPlayer = true
				// board[4][4] = 1
				// fmt.Printf("move %c%c\n", 'A'+4, 'A'+4)
				aiMove()
			} else {
				IsMaxPlayer = false
			}
			step = 1
		} else if strings.HasPrefix(line, "move") {
			words := strings.Split(line, " ")
			move := words[1]
			X := move[0] - 'A'
			Y := move[1] - 'A'
			board[X][Y] = 1
			step++
			nogo.AddMove(!IsMaxPlayer, int(X), int(Y))
			aiMove()
		} else if line == "end" {
			nogo.Save()
			fmt.Println("Game over.")
			resetBoard()
			continue
		} else {
			fmt.Println("Unknown command. Available commands: 'new [color]', 'move x y', 'end', 'quit'.")
		}
	}

	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "reading standard input: %s\n", err)
	}
}

func resetBoard() {
	for i := range board {
		for j := range board[i] {
			board[i][j] = 0
		}
	}
	clearVisit()
}
