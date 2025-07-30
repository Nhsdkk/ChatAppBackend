package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

const KBytesInMBytes = 1024
const Separator = "$"

type HashPasswordConfig struct {
	SaltSize   uint32
	Iterations uint32
	Memory     uint32
	Threads    uint8
	KeyLength  uint32
}

var options = HashPasswordConfig{
	SaltSize:   16,
	Iterations: 3,
	Memory:     64 * KBytesInMBytes,
	Threads:    2,
	KeyLength:  32,
}

func generateRandomSalt(saltSize uint32) []byte {
	bytes := make([]byte, saltSize)
	_, _ = rand.Read(bytes)
	return bytes
}

func hashWithSalt(password string, salt []byte, opts *HashPasswordConfig) []byte {
	hashedPassword := argon2.IDKey(
		[]byte(password),
		salt,
		opts.Iterations,
		opts.Memory,
		opts.Threads,
		opts.KeyLength,
	)
	return hashedPassword
}

func hashPassword(password string, opts *HashPasswordConfig) (salt []byte, hashedPassword []byte) {
	salt = generateRandomSalt(opts.SaltSize)
	hashedPassword = hashWithSalt(password, salt, opts)
	return salt, hashedPassword
}

func HashPassword(password string) []byte {
	salt, hashedPassword := hashPassword(password, &options)
	return []byte(
		fmt.Sprintf(
			"%s%s%s",
			base64.RawStdEncoding.EncodeToString(salt),
			Separator,
			base64.RawStdEncoding.EncodeToString(hashedPassword),
		),
	)
}

func ComparePassword(rawPassword string, encodedPassword []byte) bool {
	salt, storedHashed := decodePasswordAndSalt(encodedPassword)
	rawHashed := hashWithSalt(rawPassword, salt, &options)

	return subtle.ConstantTimeCompare(rawHashed, storedHashed) == 1
}

func decodePasswordAndSalt(encodedPassword []byte) (salt []byte, password []byte) {
	str := string(encodedPassword)
	split := strings.Split(str, Separator)
	b64Salt, b64Password := split[0], split[1]
	salt, _ = base64.RawStdEncoding.DecodeString(b64Salt)
	password, _ = base64.RawStdEncoding.DecodeString(b64Password)

	return salt, password
}
