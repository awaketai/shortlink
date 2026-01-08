package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"shortlink/internal/api/http/server"
	"shortlink/internal/config"
	"shortlink/internal/idgen"
	"shortlink/internal/shortener"
	"shortlink/internal/storage"
	"syscall"
	"time"
)

func main() {
	c, _ := config.LoadConfig()
	// 使用标准库
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting shortlink service", "version", "shortlink-demo1")
	// 初始化依赖
	storeImpl := storage.NewMemoryStore()
	defer func() {
		if err := storeImpl.Close(); err != nil {
			log.Println("Failed to close storage:", err)
		}
	}()
	idGenImpl := idgen.NewGenerator()

	shortenerSvc := shortener.NewService(shortener.Config{
		Store:         storeImpl,
		Generator:     idGenImpl,
		MaxGenAttemps: 3,
	})
	if shortenerSvc == nil {
		log.Fatal("Failed to create shortener service")
	}
	// 创建http服务器
	httpServer := server.NewServer(c.Server.Port, shortenerSvc)
	go func() {
		if err := httpServer.Start(); err != nil {
			log.Fatal("Failed to start http server:", err)
		}
	}()
	// 实现优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGILL, syscall.SIGTERM)
	sig := <-quit
	log.Printf("Receive signal %s,shutting down server...\n", sig.String())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Println("Failed to shutdown http server:", err)
	}

	log.Println("Server exiting")

}
