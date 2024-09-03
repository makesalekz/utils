package auth

import (
	"time"

	u_jwt "gitlab.calendaria.team/services/utils/v3/jwt"

	"github.com/go-kratos/kratos/v2/middleware"
	kjwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/golang-jwt/jwt/v5"
)

const S2S_TOKEN_DURATION = 60 * time.Minute

// Server is a middleware for gRPC-requests on microservices that extracts the claims from the s2s-jwt token and adds them to the context.
func Server(jwtp u_jwt.IJwtSecret) middleware.Middleware {
	return kjwt.Server(
		func(token *jwt.Token) (interface{}, error) {
			return jwtp.Get(), nil
		}, kjwt.WithSigningMethod(jwt.SigningMethodHS256),
		kjwt.WithClaims(func() jwt.Claims { return &jwt.RegisteredClaims{} }),
	)
}

// Client is a middleware that adds the jwt token to the client grpc request.
func Client(jwtSecret []byte, issuer, audience string) middleware.Middleware {
	return kjwt.Client(
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		}, kjwt.WithSigningMethod(jwt.SigningMethodHS256), kjwt.WithClaims(
			func() jwt.Claims {
				return &jwt.RegisteredClaims{
					Issuer:    issuer,
					Audience:  jwt.ClaimStrings{audience},
					Subject:   "s2s",
					IssuedAt:  jwt.NewNumericDate(time.Now()),
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(S2S_TOKEN_DURATION)),
				}
			},
		),
	)
}
