package auth

import (
	"context"
	"strconv"
	"strings"
	"time"

	"gitlab.calendaria.team/services/utils/v1/jwt"
	"gitlab.calendaria.team/services/utils/v2/auth"

	errors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	kjwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwtv4 "github.com/golang-jwt/jwt/v4"
)

const S2S_TOKEN_DURATION = 60 * time.Minute

// Server is a middleware that extracts the claims from the jwt token and adds them to the context.
func Server(jwtp *jwt.JwtProcessor) middleware.Middleware {
	return kjwt.Server(func(token *jwtv4.Token) (interface{}, error) {
		return jwtp.GetSecret(), nil
	}, kjwt.WithSigningMethod(jwtv4.SigningMethodHS256), kjwt.WithClaims(func() jwtv4.Claims { return &jwt.TenantClaims{} }))
}

// BffMetaServer is a middleware that extracts the actor id and tenant id from the jwt token and adds them to the metadata in global context.
// It allows to use the actor id and tenant id in the whole service calls chain.
// Requires the jwt & metadata middlewares to be used before this middleware.
// Requires user claims to be passed in the jwt token.
func BffMetaServer(jwtp *jwt.JwtProcessor) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			claims, _ := jwtp.GetClaimsFromContext(ctx)

			if !claims.IsUserRequest() {
				return nil, errors.New(401, "UNAUTHORIZED", "user claims are required")
			}

			actorId := claims.GetUserId()
			ctx = auth.NewActorContext(ctx, actorId)
			ctx = metadata.AppendToClientContext(ctx, "x-md-global-actor-id", strconv.FormatInt(actorId, 10))

			tenantId := claims.GetTenantId()
			if tenantId != 0 {
				ctx = auth.NewTenantContext(ctx, tenantId)
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

// Client is a middleware that adds the jwt token to the client grpc request.
func Client(
	jwtp *jwt.JwtProcessor,
	issuer string,
	audience jwtv4.ClaimStrings,
) middleware.Middleware {
	return kjwt.Client(func(token *jwtv4.Token) (interface{}, error) {
		return jwtp.GetSecret(), nil
	}, kjwt.WithSigningMethod(jwtv4.SigningMethodHS256), kjwt.WithClaims(func() jwtv4.Claims {
		return &jwtv4.RegisteredClaims{
			Issuer:    issuer,
			Audience:  audience,
			Subject:   "s2s",
			IssuedAt:  jwtv4.NewNumericDate(time.Now()),
			ExpiresAt: jwtv4.NewNumericDate(time.Now().Add(S2S_TOKEN_DURATION)),
		}
	}))
}
