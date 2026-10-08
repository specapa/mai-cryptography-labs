package symmetric

import (
	"context"
	"fmt"
	"io"
	"os"
)

type SymmetricContext struct {
	algorithm SymmetricAlgorithm
	mode      CipherMode
	padding   PaddingMode
	iv        []byte
	delta     []byte
}

func NewSymmetricContext(algorithm SymmetricAlgorithm, key []byte, mode CipherMode, padding PaddingMode, iv []byte, extra ...[]byte) (*SymmetricContext, error) {
	if algorithm == nil {
		return nil, fmt.Errorf("algorithm is nil")
	}
	if err := algorithm.SetKey(key); err != nil {
		return nil, err
	}
	bs := algorithm.BlockSize()
	if mode.needsIV() {
		if len(iv) != bs {
			return nil, fmt.Errorf("%s requires an IV of %d bytes, got %d", mode, bs, len(iv))
		}
	}
	c := &SymmetricContext{
		algorithm: algorithm,
		mode:      mode,
		padding:   padding,
		iv:        cloneBytes(iv),
	}
	if mode == RandomDelta {
		if len(extra) > 0 && len(extra[0]) > 0 {
			if len(extra[0]) > bs {
				return nil, fmt.Errorf("random delta must be at most %d bytes", bs)
			}
			c.delta = make([]byte, bs)
			copy(c.delta[bs-len(extra[0]):], extra[0])
		} else {
			c.delta = make([]byte, bs)
			half := bs / 2
			if half == 0 {
				half = 1
			}
			copy(c.delta[bs-half:], iv[bs-half:])
		}
		if isZero(c.delta) {
			c.delta[bs-1] = 1
		}
	}
	return c, nil
}

func (c *SymmetricContext) Encrypt(data []byte, out *[]byte) error {
	return c.EncryptContext(context.Background(), data, out)
}

func (c *SymmetricContext) Decrypt(data []byte, out *[]byte) error {
	return c.DecryptContext(context.Background(), data, out)
}

func (c *SymmetricContext) EncryptContext(ctx context.Context, data []byte, out *[]byte) error {
	res := <-c.EncryptAsync(ctx, data)
	if res.Err != nil {
		return res.Err
	}
	*out = res.Data
	return nil
}

func (c *SymmetricContext) DecryptContext(ctx context.Context, data []byte, out *[]byte) error {
	res := <-c.DecryptAsync(ctx, data)
	if res.Err != nil {
		return res.Err
	}
	*out = res.Data
	return nil
}

type AsyncResult struct {
	Data []byte
	Err  error
}

func (c *SymmetricContext) EncryptAsync(ctx context.Context, data []byte) <-chan AsyncResult {
	ch := make(chan AsyncResult, 1)
	go func() {
		defer close(ch)
		out, err := c.encryptSync(ctx, data)
		ch <- AsyncResult{Data: out, Err: err}
	}()
	return ch
}

func (c *SymmetricContext) DecryptAsync(ctx context.Context, data []byte) <-chan AsyncResult {
	ch := make(chan AsyncResult, 1)
	go func() {
		defer close(ch)
		out, err := c.decryptSync(ctx, data)
		ch <- AsyncResult{Data: out, Err: err}
	}()
	return ch
}

func (c *SymmetricContext) EncryptFile(inPath, outPath string) error {
	return c.EncryptFileContext(context.Background(), inPath, outPath)
}

func (c *SymmetricContext) DecryptFile(inPath, outPath string) error {
	return c.DecryptFileContext(context.Background(), inPath, outPath)
}

func (c *SymmetricContext) EncryptFileContext(ctx context.Context, inPath, outPath string) error {
	return <-c.EncryptFileAsync(ctx, inPath, outPath)
}

func (c *SymmetricContext) DecryptFileContext(ctx context.Context, inPath, outPath string) error {
	return <-c.DecryptFileAsync(ctx, inPath, outPath)
}

func (c *SymmetricContext) EncryptFileAsync(ctx context.Context, inPath, outPath string) <-chan error {
	ch := make(chan error, 1)
	go func() {
		defer close(ch)
		ch <- c.cryptFile(ctx, inPath, outPath, true)
	}()
	return ch
}

func (c *SymmetricContext) DecryptFileAsync(ctx context.Context, inPath, outPath string) <-chan error {
	ch := make(chan error, 1)
	go func() {
		defer close(ch)
		ch <- c.cryptFile(ctx, inPath, outPath, false)
	}()
	return ch
}

func (c *SymmetricContext) encryptSync(ctx context.Context, data []byte) ([]byte, error) {
	padded, err := pad(data, c.algorithm.BlockSize(), c.padding)
	if err != nil {
		return nil, err
	}
	blocks := splitBlocks(padded, c.algorithm.BlockSize())
	outBlocks, err := c.encryptBlocks(ctx, blocks)
	if err != nil {
		return nil, err
	}
	return joinBlocks(outBlocks), nil
}

func (c *SymmetricContext) decryptSync(ctx context.Context, data []byte) ([]byte, error) {
	bs := c.algorithm.BlockSize()
	if len(data)%bs != 0 {
		return nil, fmt.Errorf("ciphertext length %d is not a multiple of %d", len(data), bs)
	}
	blocks := splitBlocks(data, bs)
	outBlocks, err := c.decryptBlocks(ctx, blocks)
	if err != nil {
		return nil, err
	}
	joined := joinBlocks(outBlocks)
	return unpad(joined, bs, c.padding)
}

func (c *SymmetricContext) cryptFile(ctx context.Context, inPath, outPath string, encrypt bool) error {
	in, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer in.Close()
	data, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	var out []byte
	if encrypt {
		out, err = c.encryptSync(ctx, data)
	} else {
		out, err = c.decryptSync(ctx, data)
	}
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, out, 0o644)
}

func isZero(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}
