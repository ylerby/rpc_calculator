package main

import (
	"fmt"
	"math"
)

type Arguments struct {
	A, B float64
}

type Calculator struct{}

func (c *Calculator) Multiply(args *Arguments, reply *float64) error {
	*reply = args.A * args.B
	return nil
}

func (c *Calculator) Divide(args *Arguments, quo *float64) error {
	if args.B == 0 {
		return fmt.Errorf("деление на ноль")
	}

	*quo = args.A / args.B
	return nil
}

func (c *Calculator) Add(args *Arguments, sum *float64) error {
	*sum = args.A + args.B
	return nil
}

func (c *Calculator) Subtract(args *Arguments, diff *float64) error {
	*diff = args.A - args.B
	return nil
}

func (c *Calculator) Sqrt(args *Arguments, sqrt *float64) error {
	if args.A < 0 {
		return fmt.Errorf("отрицательное число")
	}

	*sqrt = math.Sqrt(args.A)
	return nil
}

func (c *Calculator) Percent(args *Arguments, percent *float64) error {
	if args.A < 0 {
		return fmt.Errorf("отрицательное число")
	}

	*percent = args.A / 100 * args.B
	return nil
}

func (c *Calculator) Round(args *Arguments, rand *float64) error {
	*rand = math.Round(args.A)
	return nil
}

func (c *Calculator) Pow(args *Arguments, pow *float64) error {
	*pow = math.Pow(args.A, args.B)
	return nil
}
