// app.go - Основная точка входа
package main

import (
    "context"
    "os"
    "os/signal"
    "syscall"

    "github.com/your-org/dkvs/node"
    "go.uber.org/zap"
)

func main() {
    // Инициализация structured logging
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    // Создание корневого контекста с graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Загрузка конфигурации (env vars, config file, flags)
    cfg := node.LoadConfig()

    // Инициализация и запуск узла
    n := node.NewNode(cfg, logger)
    
    if err := n.Start(ctx); err != nil {
        logger.Fatal("Failed to start node", zap.Error(err))
    }

    // Ожидание сигналов завершения
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    
    select {
    case sig := <-sigCh:
        logger.Info("Received signal, shutting down", zap.String("signal", sig.String()))
        n.Stop()
    case <-ctx.Done():
    }
}
