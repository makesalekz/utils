package jwt

import (
	"context"
	"fmt"

	"github.com/makesalekz/utils/v1/config"
)

type IJwtSecret interface {
	Get() []byte
}

type jwtSecret struct {
	data []byte
}

func (j *jwtSecret) Get() []byte {
	return j.data
}

// NewGlobalJwt .
func NewGlobalJwt(c *config.Config) (IJwtSecret, error) {
	data, err := c.ReadJwt(context.TODO(), "global")
	if err != nil {
		return nil, fmt.Errorf("jwt secret not found, error: %w", err)
	}

	return &jwtSecret{
		data: data,
	}, nil
}

// NewPrivateJwt .
func NewPrivateJwt(c *config.Config) (IJwtSecret, error) {
	data, err := c.ReadJwt(context.TODO(), c.GetAppName())
	if err != nil {
		return nil, fmt.Errorf("jwt secret not found, error: %w", err)
	}

	return &jwtSecret{
		data: data,
	}, nil
}
