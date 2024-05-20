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
    // Преобразуем ключ из конфигурации шифрования в байты
    key := []byte(l.encryptionConfig.Key) 
    // Создаем шифр на основе ключа
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    // Преобразуем строку для шифрования в байты
    plaintext := []byte(val)
    // Создаем массив для шифротекста с дополнительным местом для IV
    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    // Используем первый блок массива в качестве IV
    iv := ciphertext[:aes.BlockSize] 
    // Заполняем IV случайными данными
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    // Создаем потоковый шифр для шифрования
    stream := cipher.NewCFBEncrypter(block, iv)
    // Применяем XOR между ключом и текстом, начиная с позиции после IV
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

    // Возвращаем зашифрованное значение
    return val, nil
}

func (l *SecurityService) Decrypt(val string) (string, error) {
    // Преобразуем ключ из конфигурации шифрования в байты
    key := []byte(l.encryptionConfig.Key) 

    // Создаем шифр на основе ключа
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    // Проверяем, достаточно ли длинный шифротекст
    if len(val) < aes.BlockSize {
        return "", errors.New("ciphertext too short")
    }

    // Используем первый блок шифротекста в качестве IV
    iv := val[:aes.BlockSize]
    val = val[aes.BlockSize:]

    // Создаем потоковый шифр для дешифрования
    stream := cipher.NewCFBDecrypter(block, []byte(iv))
    // Применяем XOR между ключом и текстом
    stream.XORKeyStream([]byte(val), []byte(val))

    // Возвращаем дешифрованное значение
    return string(val), nil
}


