package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Order struct {
	ID        int
	UserID    string
	Amount    float64
	Status    string
	CreatedAt time.Time
}

// orderProduct 生产者, 生成订单
func orderProduct(orderChan chan<- Order, number int) {
	for i := 1; i <= number; i++ {
		order := Order{
			ID:        0,
			UserID:    fmt.Sprintf("user_%d", rand.Intn(100)),
			Amount:    rand.Float64() * 1000,
			Status:    "Pending",
			CreatedAt: time.Now(),
		}
		orderChan <- order
		fmt.Printf("生成订单: ID=%d, 用户ID=%s, 金额=%.2f\n", order.ID, order.UserID, order.Amount)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
	close(orderChan)
}

// orderProcessor 订单处理器, 处理订单
func orderProcessor(orderChan <-chan Order, resultChan chan<- Order) {
	for order := range orderChan {
		order.ID = rand.Intn(100000)
		order.Status = "Completed"
		resultChan <- order
	}
	close(resultChan)
}

// orderResultCollector 订单结果收集器, 收集订单处理结果
func orderResultCollector(resultChan <-chan Order, done chan<- bool) {
	for order := range resultChan {
		fmt.Printf("订单处理结果: ID=%d, 用户ID=%s, 状态=%s\n", order.ID, order.UserID, order.Status)
	}
	done <- true
}

func main() {
	orderChan := make(chan Order, 100)
	resultChan := make(chan Order, 100)
	done := make(chan bool)

	// 生产者:消费者 = 1:3
	go orderProduct(orderChan, 20)

	for i := 0; i < 3; i++ {
		go orderProcessor(orderChan, resultChan)
	}

	go orderResultCollector(resultChan, done)

	<-done
}
