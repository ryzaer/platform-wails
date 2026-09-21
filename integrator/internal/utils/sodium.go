package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"app-platform/internal/config"

	"golang.org/x/crypto/nacl/secretbox"
)

const (
	sodiumKeySize   = 32
	sodiumNonceSize = 24
)

// func sodiumKeyCheck() (*[32]byte, error) {

// 	data, err := base64.StdEncoding.DecodeString(config.App.SodiumKey)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(data) != sodiumKeySize {
// 		return nil, errors.New("invalid sodium key")
// 	}

// 	var key [32]byte
// 	copy(key[:], data)

// 	return &key, nil
// }

func sodiumKey(textKeys string) *[32]byte {

	hash := sha256.Sum256([]byte(textKeys))
	return &hash
}

func Encrypt(text string, keys ...string) (string, error) {
	vkey := config.App.SodiumKey
	if len(keys) > 0 && keys[0] != "" {
		vkey = keys[0]
	}
	key := sodiumKey(string(vkey))

	var nonce [24]byte

	if _, err := rand.Read(nonce[:]); err != nil {
		return "", err
	}

	encrypted := secretbox.Seal(
		nonce[:],
		[]byte(text),
		&nonce,
		key,
	)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// untuk API React
func EnSodiumKeyU8Api(tkeys string) []int {

	key := sodiumKey(tkeys)

	out := make([]int, len(key))
	for i, b := range key {
		out[i] = int(b)
	}

	return out
}

func Decrypt(cipher string, keys ...string) (string, error) {
	vkey := config.App.SodiumKey
	if len(keys) > 0 && keys[0] != "" {
		vkey = keys[0]
	}
	key := sodiumKey(string(vkey))

	data, err := base64.StdEncoding.DecodeString(cipher)
	if err != nil {
		return "", err
	}

	if len(data) < sodiumNonceSize {
		return "", errors.New("invalid ciphertext")
	}

	var nonce [24]byte
	copy(nonce[:], data[:24])

	plain, ok := secretbox.Open(
		nil,
		data[24:],
		&nonce,
		key,
	)

	if !ok {
		return "", errors.New("decrypt failed")
	}

	return string(plain), nil
}

func GenerateSodiumKey() (string, error) {

	var key [32]byte

	if _, err := rand.Read(key[:]); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(key[:]), nil
}
