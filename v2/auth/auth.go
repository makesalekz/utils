package auth

import (
	"context"
	"strconv"

	"github.com/go-kratos/kratos/v2/metadata"
)

type actorKey struct{}
type tenantKey struct{}

func NewActorContext(ctx context.Context, actorId int64) context.Context {
	return context.WithValue(ctx, actorKey{}, actorId)
}

func NewTenantContext(ctx context.Context, tenantId int64) context.Context {
	return context.WithValue(ctx, tenantKey{}, tenantId)
}

func GetActorIdFromContext(ctx context.Context) int64 {
	actorId, ok := ctx.Value(actorKey{}).(int64)
	if ok {
		return actorId
	}

	if md, ok := metadata.FromServerContext(ctx); ok {
		idString := md.Get("x-md-global-actor-id")
		if idString != "" {
			id, err := strconv.ParseInt(idString, 10, 64)
			if err == nil {
				return id
			}
		}
	}
	return 0
}

func GetTenantIdFromContext(ctx context.Context) int64 {
	tenantId, ok := ctx.Value(tenantKey{}).(int64)
	if ok {
		return tenantId
	}

	if md, ok := metadata.FromServerContext(ctx); ok {
		idString := md.Get("x-md-global-tenant-id")
		if idString != "" {
			id, err := strconv.ParseInt(idString, 10, 64)
			if err == nil {
				return id
			}
		}
	}
	return 0
}

func AppendAuthIds(ctx context.Context, userId, tenantId int64) context.Context {
	ctx = metadata.AppendToClientContext(ctx, "x-md-global-actor-id", strconv.FormatInt(userId, 10))

	if tenantId != 0 {
		ctx = metadata.AppendToClientContext(ctx, "x-md-global-tenant-id", strconv.FormatInt(tenantId, 10))
	}

	return ctx
}
