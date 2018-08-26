package logger

import (
	"context"
	"log"
)

const (
	requestIDKey = iota
	ipIDKey
)

func Log(ctx context.Context, format string, v ...interface{}) {
	fmt := "[%+v] " + format
	uuid := ctx.Value(requestIDKey)
	var values []interface{}
	values = append(values, uuid)
	values = append(values, v...)
	log.Printf(fmt, values...)
}

func Error(ctx context.Context, err error) {
	fmt := "[%+v] Error: %+v"
	uuid := ctx.Value(requestIDKey)
	var values []interface{}
	values = append(values, uuid, err)
	log.Printf(fmt, values...)
}

func Println(ctx context.Context, message string) {
	Log(ctx, message)
}
