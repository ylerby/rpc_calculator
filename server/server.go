package main

import (
	"go.uber.org/zap"
	"log"
	"net"
	"net/rpc"
	"rpc_calculator_lab/logger"
)

const (
	serverLoggerKey   = "component"
	serverLoggerValue = "rpc_server"
	serverProtocol    = "tcp"
	serverAddress     = ":5050"
	serverLogDir      = "server_logs"
	serverLogFilePath = "server_logs/server.log"
)

func main() {
	zapLogger, err := logger.New(serverLogDir, serverLogFilePath)
	if err != nil {
		log.Fatalf("ошибка при инициализации логгера - %s", err)
	}

	defer func() {
		_ = zapLogger.Sync()
	}()

	l := zapLogger.With(zap.String(serverLoggerKey, serverLoggerValue))

	calculator := new(Calculator)

	err = rpc.Register(calculator)
	if err != nil {
		l.Fatalf("ошибка при инициализации rpc - %s", err)
	}

	listener, err := net.Listen(serverProtocol, serverAddress)
	if err != nil {
		l.Fatalf("ошибка при запуске сервера - %s", err)
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			l.Fatalf("ошибка при установлении соединения - %s", err)
		}
		go rpc.ServeConn(connection)
	}
}
