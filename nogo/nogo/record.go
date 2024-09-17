package nogo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Move struct {
	Color string // "B" for black, "W" for white
	Pos   string // Position like "A1", "B2", etc.
}

var moves []Move

func AddMove(IsMaxPlayer bool, x int, y int) {
	if IsMaxPlayer {
		moves = append(moves, Move{"B", fmt.Sprintf("%c%d", 'A'+x, y+1)})
	} else {
		moves = append(moves, Move{"W", fmt.Sprintf("%c%d", 'A'+x, y+1)})
	}
}

// Generates a Go game record in a specific format.
func generateGameRecord() string {
	var sb strings.Builder
	sb.WriteString("([NG][A队][B队][先手/后手胜][2020.8.23.840线上][2020CCGC];\n")
	for _, move := range moves {
		sb.WriteString(fmt.Sprintf("%s:[%s];\n", move.Color, move.Pos))
	}
	sb.WriteString(")")
	return sb.String()
}

// SaveGameRecord saves the game record to a file on the desktop.
func Save() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v\n", err)
		return
	}

	desktopPath := filepath.Join(homeDir, "Desktop")
	filename := filepath.Join(desktopPath, "GameRecord.txt")

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	record := generateGameRecord()
	if _, err := writer.WriteString(record); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	if err := writer.Flush(); err != nil {
		fmt.Printf("Error flushing to file: %v\n", err)
	}
	fmt.Printf("Game record saved to %s\n", filename)
}
