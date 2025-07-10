package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type GameMessage struct {
	Board [][]int `json:"board"`
	Time  string  `json:"time"`
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	colMap = map[byte]int{
		'A': 0, 'B': 1, 'C': 2, 'D': 3, 'E': 4,
		'F': 5, 'G': 6, 'H': 7, 'J': 8,
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("连接升级失败:", err)
		return
	}
	defer conn.Close()

	// 发送连接确认
	sendJSON(conn, map[string]string{"message": "连接成功，欢迎进入游戏！"})

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取消息出错:", err)
			break
		}

		var gm GameMessage
		if err := json.Unmarshal(msg, &gm); err != nil {
			log.Println("JSON解析错误:", err)
			continue
		}

		board9, ok := convertBoard(gm.Board)
		if !ok {
			log.Println("非法棋盘数据，必须为9x9大小")
			continue
		}
		if gm.Time == "" {
			gm.Time = "3" // 默认时间为10秒
		}
		move := bestMove(board9, gm.Time)
		var response [2]int

		if move == "resign" {
			response = [2]int{-1, -1}
		} else {
			response = [2]int{colMap[move[0]], int(move[1] - '1')}
		}

		sendJSON(conn, response)
		fmt.Println("当前最佳走法:", move)
	}
}

func sendJSON(conn *websocket.Conn, data interface{}) {
	if jsonData, err := json.Marshal(data); err == nil {
		conn.WriteMessage(websocket.TextMessage, jsonData)
	}
}

func convertBoard(input [][]int) ([9][9]int, bool) {
	var out [9][9]int
	if len(input) != 9 {
		return out, false
	}

	for i := 0; i < 9; i++ {
		if len(input[i]) != 9 {
			return out, false
		}
		copy(out[i][:], input[i])
	}
	return out, true
}

func bestMove(board [9][9]int, time string) string {
	// 构建AI引擎命令
	commands := buildCommands(board)

	// 执行AI计算
	return executeAI(commands, time)
}

func buildCommands(board [9][9]int) []string {
	var commands []string
	var blackCount, whiteCount int

	// 打印棋盘状态
	printBoard(board)

	// 构建棋盘状态命令
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			switch board[x][y] {
			case 1:
				commands = append(commands, fmt.Sprintf("play b %c%d", 'A'+x, y+1))
				blackCount++
			case -1:
				commands = append(commands, fmt.Sprintf("play w %c%d", 'A'+x, y+1))
				whiteCount++
			}
		}
	}

	// 决定下一步颜色
	if blackCount <= whiteCount {
		commands = append(commands, "genmove b")
	} else {
		commands = append(commands, "genmove w")
	}

	return commands
}

func printBoard(board [9][9]int) {
	fmt.Println("棋盘当前状态:")
	symbols := map[int]string{1: "●", -1: "○", 0: "+"}

	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			fmt.Print(symbols[board[x][y]])
		}
		fmt.Println()
	}
}

func executeAI(commands []string, times string) string {
	cmd := exec.Command("./trendy.exe", times)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("获取输入管道失败: %v", err)
		return "resign"
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("获取输出管道失败: %v", err)
		return "resign"
	}

	if err := cmd.Start(); err != nil {
		log.Printf("启动命令失败: %v", err)
		return "resign"
	}

	// 发送命令
	go func() {
		defer stdin.Close()
		time.Sleep(100 * time.Millisecond)

		for _, command := range commands {
			stdin.Write([]byte(command + "\n"))
			time.Sleep(10 * time.Millisecond)
		}

		time.Sleep(time.Second)
		stdin.Write([]byte("quit\n"))
	}()

	// 读取输出
	bestMove := ""
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		output := scanner.Text()
		if strings.HasPrefix(output, "=") && len(output) > 1 {
			bestMove = strings.TrimSpace(output[1:])
		}
		if strings.HasPrefix(output, "*") {
			fmt.Println("本次模拟次数:", strings.TrimSpace(output[1:]))
		}
		if strings.HasPrefix(output, "%") {
			fmt.Println("当前胜率:", strings.TrimSpace(output[1:]))
		}

	}

	cmd.Wait()

	if bestMove == "" {
		return "resign"
	}
	return bestMove
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("name:go-neat version:1.0.0")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
