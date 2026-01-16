package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/uuid"
	"github.com/hisamafahri/securelogin/pkg/utils"
)

type JWTKeyPair struct {
	KID        string
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

type persistedKeyPair struct {
	KID        string `json:"kid"`
	PrivateKey string `json:"private_key"`
}

type AccessTokenClaims struct {
	jwt.Claims
	Scope string `json:"scope,omitempty"`
}

type IDTokenClaims struct {
	jwt.Claims
	Name          *string `json:"name,omitempty"`
	Email         string  `json:"email"`
	EmailVerified bool    `json:"email_verified"`
	Picture       *string `json:"picture,omitempty"`
}

type JWTService struct {
	keypairs     []*JWTKeyPair
	currentKeyID string
	issuer       string
	keyStorePath string
}

func NewJWTService(issuer string, keyStorePath string) (*JWTService, error) {
	service := &JWTService{
		issuer:       issuer,
		keyStorePath: keyStorePath,
	}

	if err := os.MkdirAll(keyStorePath, 0700); err != nil {
		return nil, fmt.Errorf("failed to create key store directory: %w", err)
	}

	keypairs, err := service.loadOrGenerateKeypairs()
	if err != nil {
		return nil, err
	}

	service.keypairs = keypairs
	service.currentKeyID = keypairs[0].KID

	return service, nil
}

func (s *JWTService) loadOrGenerateKeypairs() ([]*JWTKeyPair, error) {
	keysFile := filepath.Join(s.keyStorePath, "keypairs.json")

	if _, err := os.Stat(keysFile); err == nil {
		return s.loadKeypairsFromFile(keysFile)
	}

	keypairs, err := s.generateKeypairs()
	if err != nil {
		return nil, err
	}

	if err := s.saveKeypairsToFile(keysFile, keypairs); err != nil {
		return nil, err
	}

	return keypairs, nil
}

func (s *JWTService) generateKeypairs() ([]*JWTKeyPair, error) {
	keypairs := make([]*JWTKeyPair, 3)

	for i := range 3 {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, fmt.Errorf("failed to generate RSA key: %w", err)
		}
		kid, _ := utils.GenerateRandomID("kid")

		keypairs[i] = &JWTKeyPair{
			KID:        kid,
			PrivateKey: privateKey,
			PublicKey:  &privateKey.PublicKey,
		}
	}

	return keypairs, nil
}

func (s *JWTService) saveKeypairsToFile(
	filename string,
	keypairs []*JWTKeyPair,
) error {
	persistedKeys := make([]persistedKeyPair, len(keypairs))

	for i, kp := range keypairs {
		privateKeyBytes := x509.MarshalPKCS1PrivateKey(kp.PrivateKey)
		privateKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		})

		persistedKeys[i] = persistedKeyPair{
			KID:        kp.KID,
			PrivateKey: string(privateKeyPEM),
		}
	}

	data, err := json.MarshalIndent(persistedKeys, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal keypairs: %w", err)
	}

	if err := os.WriteFile(filename, data, 0600); err != nil {
		return fmt.Errorf("failed to write keypairs file: %w", err)
	}

	return nil
}

func (s *JWTService) loadKeypairsFromFile(
	filename string,
) ([]*JWTKeyPair, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read keypairs file: %w", err)
	}

	var persistedKeys []persistedKeyPair
	if err := json.Unmarshal(data, &persistedKeys); err != nil {
		return nil, fmt.Errorf("failed to unmarshal keypairs: %w", err)
	}

	keypairs := make([]*JWTKeyPair, len(persistedKeys))

	for i, pk := range persistedKeys {
		block, _ := pem.Decode([]byte(pk.PrivateKey))
		if block == nil {
			return nil, fmt.Errorf("failed to decode PEM block for key %s", pk.KID)
		}

		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse private key for %s: %w",
				pk.KID,
				err,
			)
		}

		keypairs[i] = &JWTKeyPair{
			KID:        pk.KID,
			PrivateKey: privateKey,
			PublicKey:  &privateKey.PublicKey,
		}
	}

	return keypairs, nil
}

func (s *JWTService) getKeyPair(kid string) *JWTKeyPair {
	for _, kp := range s.keypairs {
		if kp.KID == kid {
			return kp
		}
	}
	return nil
}

func (s *JWTService) getCurrentKeyPair() *JWTKeyPair {
	return s.getKeyPair(s.currentKeyID)
}

func (s *JWTService) GenerateAccessToken(
	userID uuid.UUID,
	clientID string,
	scope string,
	expiresIn time.Duration,
) (string, error) {
	now := time.Now()

	claims := AccessTokenClaims{
		Claims: jwt.Claims{
			Issuer:   s.issuer,
			Subject:  userID.String(),
			Audience: jwt.Audience{clientID},
			IssuedAt: jwt.NewNumericDate(now),
			Expiry:   jwt.NewNumericDate(now.Add(expiresIn)),
		},
		Scope: scope,
	}

	return s.signToken(claims)
}

func (s *JWTService) GenerateIDToken(
	userID uuid.UUID,
	clientID string,
	email string,
	name *string,
	picture *string,
	expiresIn time.Duration,
) (string, error) {
	now := time.Now()

	claims := IDTokenClaims{
		Claims: jwt.Claims{
			Issuer:   s.issuer,
			Subject:  userID.String(),
			Audience: jwt.Audience{clientID},
			IssuedAt: jwt.NewNumericDate(now),
			Expiry:   jwt.NewNumericDate(now.Add(expiresIn)),
		},
		Name:          name,
		Email:         email,
		EmailVerified: true,
		Picture:       picture,
	}

	return s.signToken(claims)
}

func (s *JWTService) signToken(claims interface{}) (string, error) {
	keyPair := s.getCurrentKeyPair()

	signer, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.RS256,
			Key:       keyPair.PrivateKey,
		},
		(&jose.SignerOptions{}).
			WithType("JWT").
			WithHeader("kid", keyPair.KID),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create signer: %w", err)
	}

	token, err := jwt.Signed(signer).Claims(claims).Serialize()
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return token, nil
}

func (s *JWTService) GetPublicKeys() []jose.JSONWebKey {
	keys := make([]jose.JSONWebKey, len(s.keypairs))

	for i, kp := range s.keypairs {
		keys[i] = jose.JSONWebKey{
			Key:       kp.PublicKey,
			KeyID:     kp.KID,
			Algorithm: string(jose.RS256),
			Use:       "sig",
		}
	}

	return keys
}

func (s *JWTService) ValidateAccessToken(
	tokenString string,
) (*AccessTokenClaims, error) {
	token, err := jwt.ParseSigned(
		tokenString,
		[]jose.SignatureAlgorithm{jose.RS256},
	)
	if err != nil {
		return nil, errors.New("invalid token format")
	}

	var claims AccessTokenClaims

	kidFound := false
	var lastErr error

	for _, header := range token.Headers {
		kid, ok := header.KeyID, header.KeyID != ""
		if !ok {
			continue
		}

		keyPair := s.getKeyPair(kid)
		if keyPair == nil {
			continue
		}

		kidFound = true

		if err := token.Claims(keyPair.PublicKey, &claims); err != nil {
			lastErr = err
			continue
		}

		if err := claims.Validate(jwt.Expected{
			Issuer: s.issuer,
			Time:   time.Now(),
		}); err != nil {
			return nil, errors.New("token validation failed")
		}

		return &claims, nil
	}

	if !kidFound {
		return nil, errors.New("token signed with unknown key")
	}

	if lastErr != nil {
		return nil, errors.New("token verification failed")
	}

	return nil, errors.New("token validation failed")
}
