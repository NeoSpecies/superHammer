package http

import (
	"bigHammer/internal/plugin/agilitymemdb"
	"bigHammer/internal/router"
	"bigHammer/internal/shared"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func StartHTTPServer(ctx context.Context, httpPort string) { // 添加一个上下文参数用于取消操作
	var loadedRouter router.Router // 使用 router 包中的 Router 类型
	var err error

	if err != nil {
		fmt.Println("Error loading router configuration:", err)
		return
	}
	db, err := shared.GlobalContainer.Resolve("database")
	if err != nil {
		log.Fatal("Error resolving db from container:", err)
	}

	// 断言 db 的类型，确保它实现了正确的接口
	dbInstance, ok := db.(*agilitymemdb.AgilityMemDB)
	if !ok {
		log.Fatal("The provided db service does not match the expected type.")
	}
	loadedRouter, _ = router.LoadRouterConfig() // 使用 router 包中的 LoadRouterConfig 函数
	// 设置数据库实例到路由器
	loadedRouter.DB = dbInstance
	log.Println("Router loaded successfully.")
	server := &http.Server{Addr: ":" + httpPort}  
	

	// 假设 Router 类型有一个 handleHTTP 方法
	http.HandleFunc("/", loadedRouter.HandleHTTP)
	log.Println("HTTP server listening on port:"+ httpPort)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("Error starting HTTP server:", err)
		}
	}()

	<-ctx.Done() // 等待上下文被取消

	// 创建一个新的超时上下文用于服务器关闭
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 尝试优雅地关闭服务器
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}
}
