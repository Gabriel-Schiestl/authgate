package adapters

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"

	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/services"
)

type encryptService struct {
    publicKey  *rsa.PublicKey
    privateKey *rsa.PrivateKey
}

func NewEncryptService() services.IEncryptService {
    publicKeyPEM := os.Getenv("RSA_PUBLIC_KEY")
    privateKeyPEM := os.Getenv("RSA_PRIVATE_KEY")

    if publicKeyPEM == "" || privateKeyPEM == "" {
        panic("RSA_PUBLIC_KEY and RSA_PRIVATE_KEY environment variables must be set")
    }

    publicKey, err := parsePublicKey(publicKeyPEM)
    if err != nil {
        panic(fmt.Sprintf("Failed to parse RSA public key: %v", err))
    }

    privateKey, err := parsePrivateKey(privateKeyPEM)
    if err != nil {
        panic(fmt.Sprintf("Failed to parse RSA private key: %v", err))
    }

    return &encryptService{
        publicKey:  publicKey,
        privateKey: privateKey,
    }
}

func (s *encryptService) Encrypt(ctx context.Context, token string) (*string, error) {
    aesKey := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
        return nil, fmt.Errorf("failed to generate AES key: %w", err)
    }

    block, err := aes.NewCipher(aesKey)
    if err != nil {
        return nil, fmt.Errorf("failed to create AES cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM: %w", err)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, fmt.Errorf("failed to generate nonce: %w", err)
    }

    ciphertext := gcm.Seal(nonce, nonce, []byte(token), nil)

    encryptedAESKey, err := rsa.EncryptOAEP(
        sha256.New(),
        rand.Reader,
        s.publicKey,
        aesKey,
        nil,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt AES key: %w", err)
    }

    result := append(encryptedAESKey, ciphertext...)
    encryptedToken := base64.StdEncoding.EncodeToString(result)
    
    return &encryptedToken, nil
}

func (s *encryptService) Decrypt(ctx context.Context, encryptedToken string) (string, error) {
    data, err := base64.StdEncoding.DecodeString(encryptedToken)
    if err != nil {
        return "", fmt.Errorf("failed to decode base64: %w", err)
    }

    rsaKeySize := s.privateKey.Size()
    if len(data) < rsaKeySize {
        return "", fmt.Errorf("encrypted data too short")
    }

    encryptedAESKey := data[:rsaKeySize]
    encryptedData := data[rsaKeySize:]

    aesKey, err := rsa.DecryptOAEP(
        sha256.New(),
        rand.Reader,
        s.privateKey,
        encryptedAESKey,
        nil,
    )
    if err != nil {
        return "", fmt.Errorf("failed to decrypt AES key: %w", err)
    }

    block, err := aes.NewCipher(aesKey)
    if err != nil {
        return "", fmt.Errorf("failed to create AES cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("failed to create GCM: %w", err)
    }

    nonceSize := gcm.NonceSize()
    if len(encryptedData) < nonceSize {
        return "", fmt.Errorf("encrypted data too short for nonce")
    }

    nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", fmt.Errorf("failed to decrypt data: %w", err)
    }

    return string(plaintext), nil
}

func parsePublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
    block, _ := pem.Decode([]byte(publicKeyPEM))
    if block == nil {
        return nil, fmt.Errorf("failed to parse PEM block containing public key")
    }

    pub, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, fmt.Errorf("failed to parse public key: %w", err)
    }

    rsaPublicKey, ok := pub.(*rsa.PublicKey)
    if !ok {
        return nil, fmt.Errorf("key is not RSA public key")
    }

    return rsaPublicKey, nil
}

func parsePrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
    block, _ := pem.Decode([]byte(privateKeyPEM))
    if block == nil {
        return nil, fmt.Errorf("failed to parse PEM block containing private key")
    }

    if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
        if rsaKey, ok := key.(*rsa.PrivateKey); ok {
            return rsaKey, nil
        }
        return nil, fmt.Errorf("PKCS#8 key is not RSA private key")
    }

    if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
        return key, nil
    }

    return nil, fmt.Errorf("failed to parse private key: unsupported format")
}