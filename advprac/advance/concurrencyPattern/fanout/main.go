// 扇入/扇出模型
// 扇入: 多个生产者, 一个消费者, 聚合多数据源
// 扇出: 多个消费者, 一个生产者, 提高处理吞吐量
package main

import (
	"advance/concurrencyPattern/customType"
	"fmt"
	"sync"
	"time"
)

// fanOutProcessor 扇出
func fanOutProcessor(inputChan <-chan customType.Order, orderType int, wg *sync.WaitGroup) {
	defer wg.Done()
	for order := range inputChan {
		switch orderType {
		case 1: // 计算折扣
			order.Amount *= 0.8
			fmt.Printf("计算折扣: ID=%d, 状态=%s, 金额: %.2f\n", order.ID, order.Status, order.Amount)
		case 2: // 发送通知
			fmt.Printf("发送通知: ID=%d, 状态=%s\n", order.ID, order.Status)
		case 3: // 记录日志
			fmt.Printf("记录日志: ID=%d, 状态=%s\n", order.ID, order.Status)
		}
		time.Sleep(time.Millisecond * 100)
	}

}
func main() {
	// ================ 扇出模式 ================ //
	fanoutCh := make(chan customType.Order, 10)
	var fanoutWg sync.WaitGroup
	fanoutWg.Add(3)
	for i := 1; i <= 3; i++ {
		go fanOutProcessor(fanoutCh, i, &fanoutWg)
	}

	// 生产者
	go func() {
		for i := 1; i <= 6; i++ {
			fanoutCh <- customType.Order{ID: i, Amount: float64(i * 100), Status: "Pending"}
		}
		close(fanoutCh)
	}()
	fanoutWg.Wait()
}
