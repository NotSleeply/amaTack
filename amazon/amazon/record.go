package amazon

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type Record struct {
	FromX  rune
	FromY  int
	ToX    rune
	ToY    int
	ArrowX rune
	ArrowY int
}

var (
	recordSlice []Record
)

func AddRecord(fromX, fromY, toX, toY, arrowX, arrowY int) {
	recordSlice = append(recordSlice, Record{
		FromX:  rune(fromX),
		FromY:  fromY,
		ToX:    rune(toX),
		ToY:    toY,
		ArrowX: rune(arrowX),
		ArrowY: arrowY,
	})
}

func Save() {
	desktopPath := "D:\\code\\"

	// 创建文件名
	filename := fmt.Sprintf(desktopPath+"AM-先手队B vs 后手队W-先(后)手胜-%v 辽宁-2024.txt",
		time.Now().Format("2006年01月02日 15时04分"))
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString("#[AM][先手参赛队][后手参赛队][后手胜][" +
		time.Now().Format("2006.01.02 15:04") + " 辽宁][2024];\r\n")
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	for i, record := range recordSlice {
		if i%2 == 0 {
			_, err = writer.WriteString(fmt.Sprintf("%v ", i/2+1))
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
		}
		_, err = writer.WriteString(fmt.Sprintf("%c%d%c%d(%c%d)", record.FromX, record.FromY, record.ToX,
			record.ToY, record.ArrowX, record.ArrowY))
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
		if i%2 == 1 {
			_, err = writer.WriteString("\r\n")
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
		} else {
			_, err = writer.WriteString(" ")
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
		}
	}
	err = writer.Flush()
	if err != nil {
		fmt.Printf("Error flushing to file: %v\n", err)
		return
	}
	recordSlice = recordSlice[:0]
}
