package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"itisadb/config"
	"itisadb/internal/constants"
	"itisadb/internal/models"

	"github.com/egorgasay/gost"
)

type SecurityService struct {
	cfg config.SecurityConfig
	encryptionConfig config.EncryptionConfig
}

func NewSecurityService(cfg config.SecurityConfig, encryptionConfig config.EncryptionConfig) *SecurityService {
	return &SecurityService{
		cfg: cfg,
		encryptionConfig: encryptionConfig,
	}
}

func (l *SecurityService) HasPermission(claimsOpt gost.Option[models.UserClaims], level models.Level) bool {
	// always ok when security is disabled
	if !l.cfg.MandatoryAuthorization {
		return true
	}

	// ok when security is not mandatory for Default level
	if level == constants.DefaultLevel {
		return true
	}

	if claimsOpt.IsNone() {
		return false
	}

	claims := claimsOpt.Unwrap()

	return claims.Level >= level
}

func (l *SecurityService) Encrypt(val string) (string, error) {
    key := []byte(l.encryptionConfig.Key) // Assuming the key is stored in the SecurityConfig and is of appropriate length
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    plaintext := []byte(val)
    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize] // using the first block as IV
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

    return val, nil
}

func (l *SecurityService) Decrypt(val string) (string, error) {
    key := []byte(l.encryptionConfig.Key) 

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    if len(val) < aes.BlockSize {
        return "", errors.New("ciphertext too short")
    }

    iv := val[:aes.BlockSize]
    val = val[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, []byte(iv))
    stream.XORKeyStream([]byte(val), []byte(val))

    return string(val), nil
}


