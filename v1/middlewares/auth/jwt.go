package auth

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	kjwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"gitlab.calendaria.team/services/utils/v1/jwt"
)

func Server(jwtp *jwt.JwtProcessor) middleware.Middleware {
	return kjwt.Server(func(token *jwtv4.Token) (interface{}, error) {
		return jwtp.GetSecret(), nil
	}, kjwt.WithSigningMethod(jwtv4.SigningMethodHS256), kjwt.WithClaims(func() jwtv4.Claims { return &jwt.TenantClaims{} }))
}

func Client(
	ctx context.Context,
	jwtp *jwt.JwtProcessor,
	defaultClaims *jwt.TenantClaims,
) middleware.Middleware {
	return kjwt.Client(func(token *jwtv4.Token) (interface{}, error) {
		return jwtp.GetSecret(), nil
	}, kjwt.WithSigningMethod(jwtv4.SigningMethodHS256), kjwt.WithClaims(func() jwtv4.Claims {
		if defaultClaims != nil {
			return defaultClaims
		}

		claims, _ := jwtp.GetClaimsFromContext(ctx)

		return claims
	}))
}
