// 流水线
// 每个阶段处理单一职责
package main

import (
	"advance/concurrencyPattern/customType"
	"fmt"
	"strconv"
	"time"
)

// 第一阶段: validate 订单校验
func validate(input <-chan customType.Order) <-chan customType.Order {
	outputCh := make(chan customType.Order, 10)
	go func() {
		defer close(outputCh)
		for order := range input {
			time.Sleep(50 * time.Millisecond)
			if order.Amount > 0 {
				order.Status = "validated"
				outputCh <- order
				fmt.Printf("验证通过: 订单: %d\n", order.ID)
			} else {
				fmt.Printf("验证失败: 订单: %d\n", order.ID)
			}
		}
	}()
	return outputCh
}

// 第二阶段: pay 支付
func pay(input <-chan customType.Order) <-chan customType.Order {
	outputCh := make(chan customType.Order, 10)
	go func() {
		defer close(outputCh)
		for order := range input {
			time.Sleep(100 * time.Millisecond)
			order.Status = "paid"
			fmt.Printf("支付成功: 订单: %d\n", order.ID)
			outputCh <- order

		}
	}()
	return outputCh
}

// 第三阶段: ship 发货
func ship(input <-chan customType.Order) <-chan customType.Order {
	outputCh := make(chan customType.Order, 10)
	go func() {
		defer close(outputCh)
		for order := range input {
			time.Sleep(150 * time.Millisecond)
			order.Status = "shipped"
			fmt.Printf("发货完成: 订单: %d\n", order.ID)
			outputCh <- order

		}
	}()
	return outputCh
}
func main() {
	pipelineInput := make(chan customType.Order, 10)

	// 构造流水线
	validateOrders := validate(pipelineInput)
	payOrders := pay(validateOrders)
	shipOrders := ship(payOrders)

	go func() {

		for i := 0; i <= 3; i++ {
			pipelineInput <- customType.Order{ID: i, UserID: strconv.Itoa(i * 10), Amount: float64(i * 100), Status: "new"}
		}
		close(pipelineInput)
	}()

	fmt.Println("流水线处理结果")
	for order := range shipOrders {
		fmt.Printf("订单处理结果: ID=%d, 状态=%s\n", order.ID, order.Status)
	}
}
