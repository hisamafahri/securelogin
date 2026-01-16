package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

var (
	ErrInvalidCodeVerifier        = errors.New("invalid code verifier")
	ErrCodeChallengeMismatch      = errors.New("code challenge mismatch")
	ErrUnsupportedChallengeMethod = errors.New(
		"unsupported code challenge method",
	)
)

func VerifyPKCE(codeChallenge, codeChallengeMethod, codeVerifier string) error {
	if codeVerifier == "" {
		return ErrInvalidCodeVerifier
	}

	var computedChallenge string

	switch codeChallengeMethod {
	case "S256":
		hash := sha256.Sum256([]byte(codeVerifier))
		computedChallenge = base64.RawURLEncoding.EncodeToString(hash[:])
	case "plain":
		computedChallenge = codeVerifier
	default:
		return ErrUnsupportedChallengeMethod
	}

	if computedChallenge != codeChallenge {
		return ErrCodeChallengeMismatch
	}

	return nil
}
