package errors

import (
	"context"
	"os"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
)

type errorFunc func(format string, args ...interface{}) *errors.Error

type ErrorHelper struct {
	log              *log.Helper
	errorWrapperFunc errorFunc
	debug            bool
}

func NewErrorHelper(
	logger log.Logger,
) *ErrorHelper {
	debug := os.Getenv("DEBUG")

	return &ErrorHelper{
		log:   log.NewHelper(logger),
		debug: debug != "",
	}
}

func (h *ErrorHelper) wrapError(err error) error {
	if h.debug {
		return err
	}

	e := errors.FromError(err)
	if e.Code >= 500 {
		// log internal errors & don't return them to the client
		h.log.Errorf(err.Error())
		return h.errorWrapperFunc("internal error")
	}

	return err
}

func (h *ErrorHelper) Build(wrapperFunc errorFunc) middleware.Middleware {
	h.errorWrapperFunc = wrapperFunc
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			reply, err := handler(ctx, req)

			return reply, h.wrapError(err)
		}
	}
}
