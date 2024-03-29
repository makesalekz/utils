package auth

import (
	"context"

	u_jwt "gitlab.calendaria.team/services/utils/v1/jwt"

	"github.com/go-kratos/kratos/v2/middleware"
	kjwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/golang-jwt/jwt/v5"
)

func Server(jwtp *u_jwt.JwtProcessor) middleware.Middleware {
	return kjwt.Server(func(token *jwt.Token) (interface{}, error) {
		return jwtp.GetSecret(), nil
	}, kjwt.WithSigningMethod(jwt.SigningMethodHS256), kjwt.WithClaims(func() jwt.Claims { return &u_jwt.TenantClaims{} }))
}

func Client(
	ctx context.Context,
	jwtp *u_jwt.JwtProcessor,
	defaultClaims *u_jwt.TenantClaims,
) middleware.Middleware {
	return kjwt.Client(func(token *jwt.Token) (interface{}, error) {
		return jwtp.GetSecret(), nil
	}, kjwt.WithSigningMethod(jwt.SigningMethodHS256), kjwt.WithClaims(func() jwt.Claims {
		if defaultClaims != nil {
			return defaultClaims
		}

		claims, _ := jwtp.GetClaimsFromContext(ctx)

		return claims
	}))
}
