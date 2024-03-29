package session

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	mathRand "math/rand"
	"strings"
)

const SecureSessionStoreKey = "securedSessionDetails"

const SecureCookieName = "session_id"

const sessionIDLength = 20

const sessionIDSignatureDelimiter = "."

func GenerateSessionID(secret string) (string, error) {
	// Generate random string
	charSet := "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP123456789_"
	var output strings.Builder
	for i := 0; i < sessionIDLength; i++ {
		random := mathRand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	tok := output.String()
	// Sign it
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(tok))
	if err != nil {
		return "", err
	}
	// base64 encode signature
	signature := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	// Concat random string with its signature using "|" as delimiter
	signedToken := []byte(fmt.Sprintf("%s%s%s", tok, sessionIDSignatureDelimiter, signature))
	// base64 the whole signed token
	return base64.URLEncoding.EncodeToString(signedToken), nil
}

func VerifySessionID(sessionID, secret string) bool {
	// base64 decode the whole sessionID
	dec, _ := base64.URLEncoding.DecodeString(sessionID)
	// Split it for "|" delimiter
	tokParts := strings.Split(string(dec), sessionIDSignatureDelimiter)
	// It should contain two parts. Random string and signature
	if len(tokParts) != 2 {
		return false
	}
	decodedPayload := []byte(tokParts[0])
	// base64 decode the signature
	decodedSignature, err := base64.URLEncoding.DecodeString(tokParts[1])
	if err != nil {
		return false
	}
	// Sign the random string again and verify it's equal
	// to the signature from the sessionID
	mac := hmac.New(sha256.New, []byte(secret))
	_, err = mac.Write(decodedPayload)
	if err != nil {
		return false
	}
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(expectedMAC, decodedSignature)
}
