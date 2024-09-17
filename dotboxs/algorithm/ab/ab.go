package AB

import (
	"gotack/dotbot/algorithm/board"
	"log"
	"math"
	"math/rand"
	"runtime"
	"sync/atomic"
	"time"
)

type AB struct {
	stop      *int32
	leafCount int32
}

const (
	noCMark int     = 8
	INF     float64 = 1e100
)

var (
	numThread int = 1
)

func (self *AB) GetName() string {
	return "AB"
}

func (self *AB) MakeMove(b *board.Board, timeout uint, verbose bool) (h, v int32, err error) {
	var (
		enterTime         = time.Now()
		bestValue float64 = -INF
		bv        [60]float64
		moves     [60]*board.Moves
		exit      = make(chan int, numThread)
	)
	self.leafCount, self.stop = 0, new(int32)
	if verbose {
		defer func() {
			log.Println("Turn:", b.Turn, ", Elapse:", time.Since(enterTime).String(), ", Leaves:", self.leafCount,
				", BestValue:", bestValue)
			runtime.GC()
		}()
	}

	if moves, _ := b.Play(); moves != nil && (moves.H != 0 || moves.V != 0 || moves.M != 0) {
		h, v = moves.Moves2HV()
		b.UnMove(moves)
		return
	} else if b.IsEnd() != 0 {
		if m, _ := b.GetCMoves(); m != nil {
			h, v = m.Moves2HV()
			return
		}
		m, _ := b.PlayRandomOne()
		h, v = m.Moves2HV()
		b.UnMove(m)
		return
	}

	mt, _ := b.PlayRandomOne()
	h, v = mt.Moves2HV()
	b.UnMove(mt)
	depthChan := make(chan int, 32)
	for i := 1; i < 30; i++ {
		depthChan <- i
	}

	for t := 0; t < numThread; t++ {
		go func() {
			bb := board.NewBoard(b.H, b.V, b.S[0], b.S[1], b.Now, b.Turn)
		LOOP:
			for len(depthChan) > 0 {
				select {
				case dep := <-depthChan:
					bv[dep], moves[dep] = self.AB(bb, bb.Now, 0.0, 2.0, dep)
					if atomic.LoadInt32(self.stop) != 0 {
						bv[dep] = -bv[dep]
						break LOOP
					}
				default:
				}
			}
			exit <- 1
		}()
	}

	go func(ptr *int32) {
		time.Sleep(time.Duration(timeout) * time.Millisecond)
		atomic.AddInt32(ptr, 1)
	}(self.stop)

	for i := 0; i < numThread; i++ {
		<-exit
	}
	for d, m := range moves {
		if m != nil {
			if bv[d] > 0 {
				bestValue = bv[d]
				h, v = m.Moves2HV()
			} else {
				if -bv[d] > bestValue {
					bestValue = -bv[d]
					h, v = m.Moves2HV()
				}
			}
		}
	}
	return
}

func (self *AB) AB(b *board.Board, root int8, alpha, beta float64, depth int) (float64, *board.Moves) {
	if depth == 0 || b.IsEnd() != 0 {
		return self.Evaluate(b, root), nil // 如果达到最大深度或游戏结束，返回当前评分
	}

	var bestMove *board.Moves
	isMaximizingPlayer := b.Now == root // 判断当前是不是最大化玩家
	bestValue := -INF

	if isMaximizingPlayer {
		bestValue = -INF
	} else {
		bestValue = INF
	}

	moves, _, _ := b.GetMove() // 假设 GetMove 返回所有可能的移动
	for _, move := range moves {
		if err := b.Move(move); err != nil { // 执行移动
			log.Panic("AB move fail.")
		}

		value, _ := self.AB(b, root, alpha, beta, depth-1) // 递归调用
		b.UnMove(move)                                     // 撤销移动

		if isMaximizingPlayer {
			if value > bestValue {
				bestValue = value
				bestMove = move
			}
			alpha = math.Max(alpha, value) // 更新 alpha
		} else {
			if value < bestValue {
				bestValue = value
				bestMove = move
			}
			beta = math.Min(beta, value) // 更新 beta
		}

		if beta <= alpha {
			break // Alpha-Beta 剪枝
		}
	}
	return bestValue, bestMove
}

func (self *AB) Evaluate(b *board.Board, root int8) (val float64) {
	return 1.0
}

func SetNumThread(n int) {
	numThread = n
}

func init() {
	rand.Seed(time.Now().Unix())
	numThread = runtime.NumCPU()
}
