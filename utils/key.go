package utils

import (
	"encoding/base64"
	"fmt"
	"github.com/kevinburke/nacl"
)

func ToKey(src []byte) (nacl.Key, error) {
	if len(src) != nacl.KeySize {
		return nil, fmt.Errorf("nacl: incorrect key length: %d", len(src))
	}
	key := new([nacl.KeySize]byte)
	copy(key[:], src)
	return key, nil
}

func LoadBase64Key(key string) (nacl.Key, error) {
	src, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	return ToKey(src)
}
