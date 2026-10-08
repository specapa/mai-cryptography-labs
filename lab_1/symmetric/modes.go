package symmetric

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

func (c *SymmetricContext) encryptBlocks(ctx context.Context, blocks [][]byte) ([][]byte, error) {
	switch c.mode {
	case ECB:
		return mapParallel(ctx, blocks, c.algorithm.EncryptBlock)
	case CBC:
		return c.encryptCBC(ctx, blocks)
	case PCBC:
		return c.encryptPCBC(ctx, blocks)
	case CFB:
		return c.encryptCFB(ctx, blocks)
	case OFB:
		return c.cryptOFB(ctx, blocks)
	case CTR:
		return c.cryptCTR(ctx, blocks)
	case RandomDelta:
		return c.encryptRandomDelta(ctx, blocks)
	default:
		return nil, fmt.Errorf("unsupported mode %s", c.mode)
	}
}

func (c *SymmetricContext) decryptBlocks(ctx context.Context, blocks [][]byte) ([][]byte, error) {
	switch c.mode {
	case ECB:
		return mapParallel(ctx, blocks, c.algorithm.DecryptBlock)
	case CBC:
		return c.decryptCBC(ctx, blocks)
	case PCBC:
		return c.decryptPCBC(ctx, blocks)
	case CFB:
		return c.decryptCFB(ctx, blocks)
	case OFB:
		return c.cryptOFB(ctx, blocks)
	case CTR:
		return c.cryptCTR(ctx, blocks)
	case RandomDelta:
		return c.decryptRandomDelta(ctx, blocks)
	default:
		return nil, fmt.Errorf("unsupported mode %s", c.mode)
	}
}

func (c *SymmetricContext) encryptCBC(ctx context.Context, blocks [][]byte) ([][]byte, error) {
	out := make([][]byte, len(blocks))
	prev := cloneBytes(c.iv)
	for i, p := range blocks {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		xored := xorBytes(p, prev)
		ciph, err := c.algorithm.EncryptBlock(xored)
		if err != nil {
			return nil, err
		}
		out[i] = ciph
		prev = ciph
	}
	return out, nil
}

func (c *SymmetricContext) decryptCBC(ctx context.Context, blocks [][]byte) ([][]byte, error) {
	plainEnc, err := mapParallel(ctx, blocks, c.algorithm.DecryptBlock)
	if err != nil {
		return nil, err
	}
	out := make([][]byte, len(blocks))
	prev := cloneBytes(c.iv)
	for i := range blocks {
		out[i] = xorBytes(plainEnc[i], prev)
		prev = blocks[i]
	}
	return out, nil
}

func (c *SymmetricContext) encryptPCBC(ctx context.Context, blocks [][]byte) ([][]byte, error) {
	out := make([][]byte, len(blocks))
	prevP := cloneBytes(c.iv)
	prevC := cloneBytes(c.iv)
	for i, p := range blocks {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		xored := xorBytes(xorBytes(p, prevP), prevC)
		if i == 0 {
			xored = xorBytes(p, c.iv)
		}
		ciph, err := c.algorithm.EncryptBlock(xored)
		if err != nil {
			return nil, err
		}
		out[i] = ciph
		prevP = p
		prevC = ciph
	}
	return out, nil
}

func (c *SymmetricContext) decryptPCBC(ctx context.Context, blocks [][]byte) ([][]byte, error) {
	out := make([][]byte, len(blocks))
	prevP := cloneBytes(c.iv)
	prevC := cloneBytes(c.iv)
	for i, ct := range blocks {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		dec, err := c.algorithm.DecryptBlock(ct)
		if err != nil {
			return nil, err
		}
		var p []byte
		if i == 0 {
			p = xorBytes(dec, c.iv)
		} else {
			p = xorBytes(xorBytes(dec, prevP), prevC)
		}
		out[i] = p
		prevP = p
		prevC = ct
	}
	return out, nil
}

func (c *SymmetricContext) encryptCFB(ctx context.Context, blocks [][]byte) ([][]byte, error) {
	out := make([][]byte, len(blocks))
	prev := cloneBytes(c.iv)
	for i, p := range blocks {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		gamma, err := c.algorithm.EncryptBlock(prev)
		if err != nil {
			return nil, err
		}
		ct := xorBytes(p, gamma)
		out[i] = ct
		prev = ct
	}
	return out, nil
}

func (c *SymmetricContext) decryptCFB(ctx context.Context, blocks [][]byte) ([][]byte, error) {
	inputs := make([][]byte, len(blocks))
	if len(blocks) > 0 {
		inputs[0] = cloneBytes(c.iv)
		for i := 1; i < len(blocks); i++ {
			inputs[i] = blocks[i-1]
		}
	}
	gammas, err := mapParallel(ctx, inputs, c.algorithm.EncryptBlock)
	if err != nil {
		return nil, err
	}
	out := make([][]byte, len(blocks))
	for i := range blocks {
		out[i] = xorBytes(blocks[i], gammas[i])
	}
	return out, nil
}

func (c *SymmetricContext) cryptOFB(ctx context.Context, blocks [][]byte) ([][]byte, error) {
	out := make([][]byte, len(blocks))
	state := cloneBytes(c.iv)
	for i, b := range blocks {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		gamma, err := c.algorithm.EncryptBlock(state)
		if err != nil {
			return nil, err
		}
		out[i] = xorBytes(b, gamma)
		state = gamma
	}
	return out, nil
}

func (c *SymmetricContext) cryptCTR(ctx context.Context, blocks [][]byte) ([][]byte, error) {
	counters := make([][]byte, len(blocks))
	ctr := cloneBytes(c.iv)
	for i := range blocks {
		counters[i] = cloneBytes(ctr)
		ctr = incrementBE(ctr)
	}
	gammas, err := mapParallel(ctx, counters, c.algorithm.EncryptBlock)
	if err != nil {
		return nil, err
	}
	out := make([][]byte, len(blocks))
	for i := range blocks {
		out[i] = xorBytes(blocks[i], gammas[i])
	}
	return out, nil
}

func (c *SymmetricContext) maskRD(i int) []byte {
	return addBigEndian(c.iv, mulBigEndian(c.delta, uint64(i)))
}

func (c *SymmetricContext) encryptRandomDelta(ctx context.Context, blocks [][]byte) ([][]byte, error) {
	masked := make([][]byte, len(blocks))
	for i, p := range blocks {
		masked[i] = xorBytes(p, c.maskRD(i))
	}
	return mapParallel(ctx, masked, c.algorithm.EncryptBlock)
}

func (c *SymmetricContext) decryptRandomDelta(ctx context.Context, blocks [][]byte) ([][]byte, error) {
	dec, err := mapParallel(ctx, blocks, c.algorithm.DecryptBlock)
	if err != nil {
		return nil, err
	}
	out := make([][]byte, len(dec))
	for i := range dec {
		out[i] = xorBytes(dec[i], c.maskRD(i))
	}
	return out, nil
}

type blockFn func([]byte) ([]byte, error)

func mapParallel(ctx context.Context, in [][]byte, fn blockFn) ([][]byte, error) {
	out := make([][]byte, len(in))
	if len(in) == 0 {
		return out, nil
	}
	workers := runtime.GOMAXPROCS(0)
	if workers < 1 {
		workers = 1
	}
	if workers > len(in) {
		workers = len(in)
	}
	idx := make(chan int)
	errCh := make(chan error, workers)
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range idx {
				if err := ctx.Err(); err != nil {
					errCh <- err
					return
				}
				res, err := fn(in[i])
				if err != nil {
					errCh <- err
					return
				}
				out[i] = res
			}
		}()
	}
	go func() {
		defer close(idx)
		for i := range in {
			select {
			case <-ctx.Done():
				return
			case idx <- i:
			}
		}
	}()
	wg.Wait()
	select {
	case err := <-errCh:
		return nil, err
	default:
		return out, nil
	}
}
