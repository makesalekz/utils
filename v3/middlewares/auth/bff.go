package auth

import (
	"context"
	"strconv"
	"strings"

	v1 "github.com/makesalekz/utils/v1/jwt"
	u_auth "github.com/makesalekz/utils/v2/auth"
	v2 "github.com/makesalekz/utils/v2/jwt"
	u_jwt "github.com/makesalekz/utils/v3/jwt"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	kjwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/golang-jwt/jwt/v5"
)

// BffServer is a middleware for HTTP-requests on BFF that extracts the claims from the jwt token and adds them to the context.
func BffServer(j u_jwt.IJwtSecret) middleware.Middleware {
	return kjwt.Server(
		func(token *jwt.Token) (interface{}, error) {
			return j.Get(), nil
		}, kjwt.WithSigningMethod(jwt.SigningMethodHS256),
		kjwt.WithClaims(func() jwt.Claims { return &v2.TenantClaims{} }),
	)
}

// BffMetaServer is a middleware that extracts the actor id and tenant id from the jwt token and adds them to the metadata in global context.
// It allows to use the actor id and tenant id in the whole service calls chain.
// Requires the jwt & metadata middlewares to be used before this middleware.
// Requires user claims to be passed in the jwt token.
func BffMetaServer(appId string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			claims, ok := getClaimsFromContext(ctx)
			if !ok {
				return nil, errors.New(401, "UNAUTHORIZED", "claims are required")
			}

			if !claims.IsUserRequest() {
				return nil, errors.New(401, "UNAUTHORIZED", "user claims are required")
			}

			actorId := claims.GetUserId()
			ctx = u_auth.NewActorContext(ctx, actorId)
			ctx = metadata.AppendToClientContext(ctx, "x-md-global-actor-id", strconv.FormatInt(actorId, 10))

			if appId != "" {
				ctx = metadata.AppendToClientContext(ctx, "x-md-global-app-id", appId)
			}

			tenantId := claims.GetTenantId()
			if tenantId != 0 {
				ctx = u_auth.NewTenantContext(ctx, tenantId)
				ctx = metadata.AppendToClientContext(ctx, "x-md-global-tenant-id", strconv.FormatInt(tenantId, 10))

				identities := claims.GetIdentities()
				if len(identities) > 0 {
					ctx = metadata.AppendToClientContext(ctx, "x-md-global-identities", strings.Join(identities, ","))
				}
			}

			return handler(ctx, req)
		}
	}
}

func getClaimsFromContext(ctx context.Context) (v2.ITenantClaims, bool) {
	token, ok := kjwt.FromContext(ctx)
	if !ok {
		return nil, false
	}

	var claims v2.ITenantClaims

	claims, ok = token.(*v2.TenantClaims)
	if !ok {
		claims, ok = token.(*v1.TenantClaims)
		if !ok {
			return nil, false
		}
	}

	return claims, true
}
