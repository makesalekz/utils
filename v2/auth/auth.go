package auth

import (
	"context"
	"strconv"

	"github.com/go-kratos/kratos/v2/metadata"
)

func GetActorIdFromContext(ctx context.Context) int64 {
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
