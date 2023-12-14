package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// 创建一个context和cancel函数
	ctx, cancel := context.WithCancel(context.Background())

	// 使用WaitGroup来跟踪所有后台任务
	var wg sync.WaitGroup

	// 启动网络服务
	wg.Add(1)
	go func() {
		defer wg.Done()
		startHTTPServer(ctx)
	}()

	// 监听终止信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// 等待接收终止信号
	<-signalChan

	// 收到终止信号，开始优雅退出
	log.Println("Received termination signal. Gracefully shutting down...")

	// 通知所有任务退出
	cancel()

	// 等待所有后台任务完成
	wg.Wait()

	log.Println("Graceful shutdown completed.")
}

// CustomHandler 是一个实现了 http.Handler 接口的自定义处理器
type CustomHandler struct {
	Message string
}

// ServeHTTP 是自定义处理器的方法，用于处理HTTP请求
func (h *CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m := 120
	fmt.Fprintf(w, "CustomHandler says: %s\n", h.Message)
	fmt.Fprintf(w, "等待start wait %v seconds\n", m)
	w.(http.Flusher).Flush()
	breakflag := false
	for i := 0; i < m; i++ {

		// 检查连接是否已经关闭
		closeNotify := w.(http.CloseNotifier).CloseNotify()
		select {
		case <-closeNotify:
			fmt.Println("Connection closed by client.")
			// 在这里执行连接关闭时的逻辑
			breakflag = true
		default:
			fmt.Fprintf(w, "count: %v \n", i)
			w.(http.Flusher).Flush()
			time.Sleep(time.Second)
		}
		if breakflag {
			break
		}
		// time.Sleep(time.Second * time.Duration(m))
	}
	if !breakflag {
		fmt.Fprintf(w, "end wait %v seconds\n", m)
	}
}

func startHTTPServer(ctx context.Context) {
	h := CustomHandler{
		Message: "aaa",
	}
	server := http.Server{
		Addr:    ":8080",
		Handler: &h,
	}
	// 启动HTTP服务
	go func() {
		log.Println("HTTP server started.")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// 等待终止信号或接收到取消信号
	select {
	case <-ctx.Done():
		log.Println("HTTP server received cancel signal. Shutting down...")
		// 关闭HTTP服务
		err := server.Shutdown(context.Background())
		if err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
		log.Println("HTTP server shut down.")
	}
}
