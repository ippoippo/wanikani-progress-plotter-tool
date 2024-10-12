package slog

import (
	"context"
	"errors"
	"log/slog"
	"os"
)

func ErrorContextFromMsgWithOSExit(ctx context.Context, msg string) {
	ErrorContextWithOSExit(ctx, "fatal error", errors.New(msg))
}

// ErrorContextWithOSExit replicates the behaviour of `log.Fatal`,
// and logs the underlying error via slogg.ErrorAttr, whilst also accepting
// a `context.Context`
func ErrorContextWithOSExit(ctx context.Context, msg string, err error) {
	slog.ErrorContext(ctx, msg, ErrorAttr(err))
	os.Exit(1)
}

// ErrorAttr returns an Attr for an error, with a key of `err`
func ErrorAttr(err error) slog.Attr {
	if err == nil {
		return slog.String("err", "<nil>")
	}
	return slog.String("err", err.Error())
}

// CtxErrorAttr returns an Attr for an error, with a key of `ctx-err`
func CtxErrorAttr(ctx context.Context) slog.Attr {
	if ctx.Err() == nil {
		return slog.String("ctx-err", "<nil>")
	}
	return slog.String("ctx-err", ctx.Err().Error())
}
