package symmetric

func addBigEndian(a, b []byte) []byte {
	n := len(a)
	out := make([]byte, n)
	var carry uint16
	bOff := n - len(b)
	for i := n - 1; i >= 0; i-- {
		sum := uint16(a[i]) + carry
		j := i - bOff
		if j >= 0 && j < len(b) {
			sum += uint16(b[j])
		}
		out[i] = byte(sum)
		carry = sum >> 8
	}
	return out
}

func mulBigEndian(delta []byte, k uint64) []byte {
	n := len(delta)
	out := make([]byte, n)
	if k == 0 {
		return out
	}
	var carry uint64
	for i := n - 1; i >= 0; i-- {
		v := uint64(delta[i])*k + carry
		out[i] = byte(v)
		carry = v >> 8
	}
	return out
}

func incrementBE(block []byte) []byte {
	out := append([]byte(nil), block...)
	for i := len(out) - 1; i >= 0; i-- {
		out[i]++
		if out[i] != 0 {
			break
		}
	}
	return out
}

func splitBlocks(data []byte, blockSize int) [][]byte {
	n := len(data) / blockSize
	blocks := make([][]byte, n)
	for i := 0; i < n; i++ {
		b := make([]byte, blockSize)
		copy(b, data[i*blockSize:(i+1)*blockSize])
		blocks[i] = b
	}
	return blocks
}

func joinBlocks(blocks [][]byte) []byte {
	if len(blocks) == 0 {
		return []byte{}
	}
	out := make([]byte, 0, len(blocks)*len(blocks[0]))
	for _, b := range blocks {
		out = append(out, b...)
	}
	return out
}

func cloneBytes(in []byte) []byte {
	if in == nil {
		return nil
	}
	out := make([]byte, len(in))
	copy(out, in)
	return out
}
