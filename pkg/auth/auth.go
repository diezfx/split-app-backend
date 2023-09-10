package auth

import (
	"fmt"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type Config struct {
	Key string
}

type Client struct {
	cfg Config
}

func New(cfg Config) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (c *Client) Validate(tokenString string) (jwt.Token, error) {
	token, err := jwt.Parse([]byte(tokenString), jwt.WithKey(jwa.HS256, []byte(c.cfg.Key)))
	if err != nil {
		return nil, fmt.Errorf("validate token: %s", err)
	}
	return token, nil
}
