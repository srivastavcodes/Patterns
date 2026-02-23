package main

import (
	"context"
	"fmt"
	"time"
)

type WithContext func(context.Context, string) (string, error)

type SlowFunction func(string) (string, error)

func Timeout(fn SlowFunction) WithContext {
	return func(ctx context.Context, str string) (string, error) {
		ch := make(chan struct {
			result string
			err    error
		}, 1)
		go func() {
			res, err := fn(str)
			ch <- struct {
				result string
				err    error
			}{res, err}
		}()
		select {
		case res := <-ch:
			return res.result, res.err
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

func Slow(string) (string, error) { return "assumed function", nil }

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	timeout := Timeout(Slow)
	res, err := timeout(ctx, "some input")

	fmt.Println(res, err)
}
