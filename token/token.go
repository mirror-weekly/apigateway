// Package token define the domain of token
package token

const (
	OK = "OK"
)

type Token interface {
	ExecuteTokenStateUpdate() error
	GetTokenString() (string, error)
	GetTokenState() string
}
