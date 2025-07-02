package services

import "context"

type IEncryptService interface {
	Encrypt(ctx context.Context, text string) (*string, error)
	Decrypt(ctx context.Context, encryptedText string) (string, error)
}