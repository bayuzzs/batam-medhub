package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidToken          = errors.New("token is invalid")
	ErrTokenExpired          = errors.New("token has expired")
	ErrInvalidTokenSignature = errors.New("token signature is invalid")
	ErrInvalidTokenIssuer    = errors.New("token issuer is invalid")
	ErrInvalidTokenAudience  = errors.New("token audience is invalid")
	ErrInvalidClaims         = errors.New("token claims are invalid")
)

// JWTClaims holds standard and custom claims required by the Batam MedHub contract.
type JWTClaims struct {
	Issuer            string `json:"iss"`
	Audience          string `json:"aud"`
	Subject           string `json:"sub"`
	SessionID         string `json:"sid"`
	PreferredCurrency string `json:"preferred_currency"`
	IssuedAt          int64  `json:"iat"`
	ExpiresAt         int64  `json:"exp"`
}

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// IssueAccessToken creates a signed HS256 JWT access token.
func IssueAccessToken(secret, issuer, audience, patientID, sessionID, preferredCurrency string, ttl time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := JWTClaims{
		Issuer:            issuer,
		Audience:          audience,
		Subject:           patientID,
		SessionID:         sessionID,
		PreferredCurrency: preferredCurrency,
		IssuedAt:          now.Unix(),
		ExpiresAt:         now.Add(ttl).Unix(),
	}

	header := jwtHeader{Alg: "HS256", Typ: "JWT"}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("marshal header: %w", err)
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal claims: %w", err)
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJSON)
	signingInput := encodedHeader + "." + encodedClaims

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	signature := mac.Sum(nil)
	encodedSignature := base64.RawURLEncoding.EncodeToString(signature)

	return signingInput + "." + encodedSignature, nil
}

// ValidateAccessToken parses and validates an HS256 JWT access token.
func ValidateAccessToken(tokenString, secret, expectedIssuer, expectedAudience string, expectedTTL time.Duration) (*JWTClaims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, ErrInvalidToken
	}

	var header jwtHeader
	if err := json.Unmarshal(headerBytes, &header); err != nil || header.Alg != "HS256" || header.Typ != "JWT" {
		return nil, ErrInvalidToken
	}

	signingInput := parts[0] + "." + parts[1]
	sigBytes, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, ErrInvalidToken
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	expectedSig := mac.Sum(nil)

	if !hmac.Equal(sigBytes, expectedSig) {
		return nil, ErrInvalidTokenSignature
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims JWTClaims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, ErrInvalidClaims
	}

	if claims.Issuer != expectedIssuer {
		return nil, ErrInvalidTokenIssuer
	}
	if claims.Audience != expectedAudience {
		return nil, ErrInvalidTokenAudience
	}

	const clockSkewSeconds = int64(60) // 60 seconds allowable clock skew for future iat

	if claims.IssuedAt <= 0 || claims.ExpiresAt <= 0 {
		return nil, ErrInvalidClaims
	}
	if claims.ExpiresAt <= claims.IssuedAt {
		return nil, ErrInvalidClaims
	}
	if expectedTTL > 0 && (claims.ExpiresAt-claims.IssuedAt) != int64(expectedTTL.Seconds()) {
		return nil, ErrInvalidClaims
	}

	now := time.Now().UTC().Unix()
	if claims.IssuedAt > (now + clockSkewSeconds) {
		return nil, ErrInvalidClaims // Token claims to be issued in the future beyond clock skew
	}
	if now >= claims.ExpiresAt {
		return nil, ErrTokenExpired // Strictly expired (no clock skew grace period on expiry)
	}

	if claims.Subject == "" || claims.SessionID == "" || (claims.PreferredCurrency != "SGD" && claims.PreferredCurrency != "IDR") {
		return nil, ErrInvalidClaims
	}

	return &claims, nil
}

// GenerateRefreshToken creates a high-entropy opaque token (43 characters).
func GenerateRefreshToken() (string, error) {
	var bytes [32]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes[:]), nil
}

// HashToken computes the lowercase SHA-256 hex string of a token or secret.
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// NewUUID generates an RFC 4122 v4 UUID string.
func NewUUID() string {
	var u [16]byte
	_, _ = rand.Read(u[:])
	u[6] = (u[6] & 0x0f) | 0x40 // version 4
	u[8] = (u[8] & 0x3f) | 0x80 // variant RFC 4122
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:16])
}
