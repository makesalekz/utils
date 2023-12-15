package auth

import (
	"github.com/go-kratos/kratos/v2/middleware"
	kjwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"gitlab.calendaria.team/services/utils/v1/jwt"
)

func Authorization(jwtp *jwt.JwtProcessor) middleware.Middleware {
	return kjwt.Server(func(token *jwtv4.Token) (interface{}, error) {
		return jwtp.GetSecret(), nil
	}, kjwt.WithSigningMethod(jwtv4.SigningMethodHS256), kjwt.WithClaims(func() jwtv4.Claims { return &jwt.TenantClaims{} }))
}
