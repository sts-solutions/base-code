package cccorrelation

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sts-solutions/base-code/ccmiddlewares"
)

type echoMiddlewareBuilder struct {
	shouldSkip func(c ccmiddlewares.Context) bool
}

func NewEchoMiddlewareBuilder() ccmiddlewares.Builder[Correlation[echo.MiddlewareFunc]] {
	mb := echoMiddlewareBuilder{
		shouldSkip: func(c ccmiddlewares.Context) bool {
			return false
		},
	}

	return &mb
}

func (b *echoMiddlewareBuilder) WithSkipper(
	skipper func(c ccmiddlewares.Context) bool) ccmiddlewares.Builder[Correlation[echo.MiddlewareFunc]] {
	b.shouldSkip = skipper
	return b
}

func (b echoMiddlewareBuilder) Build() Correlation[echo.MiddlewareFunc] {
	return echoCorrelation{
		shouldSkip: func(e echo.Context) bool {
			return b.shouldSkip(e)
		},
	}
}

type echoCorrelation struct {
	shouldSkip middleware.Skipper
}

func (co echoCorrelation) Setup() echo.MiddlewareFunc {
	return func() echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(ctx echo.Context) error {
				if co.shouldSkip(ctx) {
					return next(ctx)
				}

				req := ctx.Request()
				cid := req.Header.Get(Key.String())
				if cid == "" {
					cid = uuid.New().String()
				}

				ctxt := context.WithValue(ctx.Request().Context(), Key, cid)
				request := ctx.Request().WithContext(ctxt)
				ctx.SetRequest(request)

				return next(ctx)
			}
		}
	}()
}
