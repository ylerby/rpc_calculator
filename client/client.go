package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"log"
	"math"
	"net/rpc"
	"rpc_calculator_lab/logger"
)

const (
	clientLoggerKey   = "component"
	clientLoggerValue = "rpc_server"
	serverProtocol    = "tcp"
	serverAddress     = ":1234"
	clientLogDir      = "client_logs"
	clientLogFilePath = "client_logs/client.log"
)

type Arguments struct {
	A, B float64
}

func main() {
	operations := map[string]struct{}{
		"Multiply": {},
		"Divide":   {},
		"Add":      {},
		"Subtract": {},
		"Sqrt":     {},
		"Percent":  {},
		"Round":    {},
		"Pow":      {},
	}

	zapLogger, err := logger.New(clientLogDir, clientLogFilePath)
	if err != nil {
		log.Fatalf("ошибка при инициализации логгера - %s", err)
	}

	defer func() {
		_ = zapLogger.Sync()
	}()

	l := zapLogger.With(zap.String(clientLoggerKey, clientLoggerValue))

	l.Info(operations)

	client, err := rpc.Dial(serverProtocol, serverAddress)
	if err != nil {
		l.Fatalf("ошибка при подключении - %s", err)
	}

	var (
		isOperationSequence bool
		operation           string
		firstArgument,
		secondArgument,
		result float64
	)

	flag.BoolVar(&isOperationSequence, "is_sequence", false, "определение последовательности операции")
	flag.StringVar(&operation, "operation", "Add", "название операции")
	flag.Float64Var(&firstArgument, "A", -math.MaxFloat64, "первый аргумент")
	flag.Float64Var(&secondArgument, "B", -math.MaxFloat64, "второй аргумент")
	flag.Parse()

	if _, ok := operations[operation]; !ok {
		l.Fatalf("операция не найдена")
	}

	if firstArgument == -math.MaxFloat64 || secondArgument == -math.MaxFloat64 {
		l.Fatalf("некорректное значение аргументов")
	}

	args := Arguments{
		A: firstArgument,
		B: secondArgument,
	}

	if isOperationSequence {
		err = parseOperationSequence(client, operations, args, l)
		if err != nil {
			fmt.Printf("ошибка при вычислении результата операции - %s\n", err)
			l.Fatalf("ошибка при вычислении результата - %s", err)
		}

		return
	}

	err = client.Call(fmt.Sprintf("Calculator.%s", operation), args, &result)
	if err != nil {
		fmt.Printf("ошибка при вычислении результата операции - %s - %s\n", operation, err)
		l.Fatalf("ошибка при вычислении результата - %s", err)
	}

	fmt.Printf("Результат: %f\n", result)
	l.Infof("Результат: %f\n", result)
}

func parseOperationSequence(
	client *rpc.Client,
	operations map[string]struct{},
	args Arguments,
	l logger.Logger,
) error {
	for operation := range operations {
		var result float64
		err := client.Call(fmt.Sprintf("Calculator.%s", operation), args, &result)
		if err != nil {
			return fmt.Errorf("%s - %w", operation, err)
		}

		fmt.Printf("Результат: %f\n", result)
		l.Infof("Результат: %f\n", result)
	}

	return nil
}
