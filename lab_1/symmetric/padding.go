package symmetric

import (
	"crypto/rand"
	"fmt"
)

func pad(data []byte, blockSize int, mode PaddingMode) ([]byte, error) {
	if blockSize <= 0 || blockSize > 255 {
		return nil, fmt.Errorf("unsupported block size %d for padding", blockSize)
	}
	n := len(data)
	padLen := blockSize - (n % blockSize)
	if padLen == 0 {
		if mode == Zeros {
			out := make([]byte, n)
			copy(out, data)
			return out, nil
		}
		padLen = blockSize
	}
	out := make([]byte, n+padLen)
	copy(out, data)
	switch mode {
	case Zeros:
		//
	case ANSIX923:
		out[len(out)-1] = byte(padLen)
	case PKCS7:
		for i := n; i < len(out); i++ {
			out[i] = byte(padLen)
		}
	case ISO10126:
		if padLen > 1 {
			if _, err := rand.Read(out[n : len(out)-1]); err != nil {
				return nil, err
			}
		}
		out[len(out)-1] = byte(padLen)
	default:
		return nil, fmt.Errorf("unknown padding mode %d", mode)
	}
	return out, nil
}

func unpad(data []byte, blockSize int, mode PaddingMode) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, fmt.Errorf("ciphertext length %d is not a multiple of block size %d", len(data), blockSize)
	}
	switch mode {
	case Zeros:
		i := len(data)
		for i > 0 && data[i-1] == 0 {
			i--
		}
		out := make([]byte, i)
		copy(out, data[:i])
		return out, nil
	case PKCS7:
		padLen := int(data[len(data)-1])
		if err := validatePadLen(padLen, blockSize, len(data)); err != nil {
			return nil, err
		}
		for i := 0; i < padLen; i++ {
			if data[len(data)-1-i] != byte(padLen) {
				return nil, fmt.Errorf("invalid PKCS7 padding")
			}
		}
		return copyPrefix(data, len(data)-padLen), nil
	case ANSIX923:
		padLen := int(data[len(data)-1])
		if err := validatePadLen(padLen, blockSize, len(data)); err != nil {
			return nil, err
		}
		for i := 0; i < padLen-1; i++ {
			if data[len(data)-padLen+i] != 0 {
				return nil, fmt.Errorf("invalid ANSI X.923 padding")
			}
		}
		return copyPrefix(data, len(data)-padLen), nil
	case ISO10126:
		padLen := int(data[len(data)-1])
		if err := validatePadLen(padLen, blockSize, len(data)); err != nil {
			return nil, err
		}
		return copyPrefix(data, len(data)-padLen), nil
	default:
		return nil, fmt.Errorf("unknown padding mode %d", mode)
	}
}

func validatePadLen(padLen, blockSize, dataLen int) error {
	if padLen < 1 || padLen > blockSize || padLen > dataLen {
		return fmt.Errorf("invalid padding length %d", padLen)
	}
	return nil
}

func copyPrefix(data []byte, n int) []byte {
	out := make([]byte, n)
	copy(out, data[:n])
	return out
}
