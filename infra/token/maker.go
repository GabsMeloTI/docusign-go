package token

type Maker interface {
	VerifyToken(token string) (*Payload, error)
}
