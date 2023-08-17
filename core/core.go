package core

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
)

func Start() (err error) {
	// TODO: 任务启动
	return nil
}

func ScheduledTask(ctx context.Context, f func(), schedule string) {
	// TODO: 定时任务 | 函数, 任务时间(时间间隔 | 时间标志[day, min, second])
	c := cron.New()
	_, err := c.AddFunc(schedule, f)
	if err != nil {
		fmt.Println("Error adding scheduled task:", err)
		return
	}
	c.Start()

	done := make(chan struct{}) // 创建通道

	go func() {
		<-ctx.Done() // 等待上下文取消
		c.Stop()
		close(done) // 关闭通道
	}()

	<-done // 等待通道关闭
}

func TestScheduledTask() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ScheduledTask(ctx, func() {
		fmt.Println("定时任务正在执行...")
	}, `*/5 * * * * *`) // 每隔 5 秒

	select {}
}
