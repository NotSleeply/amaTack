// 定义评估函数，用于AI决策中的评分机制。
package amazon

import (
	"github.com/tongque0/gotack"
)

// 将通用的评估选项转换为特定的 AmazonBoard 评估值，调用 CalculateEvaluationValue 方法计算评分
func EvaluateFunc(opt *gotack.EvalOptions) float64 {
	// 尝试将board转换为*AmazonBoard类型
	amazonBoard, ok := opt.Board.(*AmazonBoard)
	if !ok {
		return 0.0 // 或者处理这种情况的其他方式
	}

	// 假设我们使用AmazonBoard的TurnID属性和isMaxPlayer来调用CalculateEvaluationValue
	value := amazonBoard.CalculateEvaluationValue(opt.Step, opt.IsMaxPlayer)
	// 返回评估分数
	return value
}
