package permutation

import "fmt"

type BitOrder int

const (
	LSBToMSB BitOrder = iota
	MSBToLSB
)

type StartIndex int

const (
	StartFromZero StartIndex = 0
	StartFromOne  StartIndex = 1
)

func bitCount(nBytes int) int {
	return nBytes * 8
}

func localize(nBytes, index int, order BitOrder, start StartIndex) (byteIdx int, mask byte, err error) {
	if start != StartFromZero && start != StartFromOne {
		return 0, 0, fmt.Errorf("start index must be 0 or 1, got %d", start)
	}
	offset := index - int(start)
	nBits := bitCount(nBytes)
	if offset < 0 || offset >= nBits {
		return 0, 0, fmt.Errorf("bit index %d is out of range for %d bytes (start=%d)", index, nBytes, start)
	}
	switch order {
	case MSBToLSB:
		byteIdx = offset / 8
		bitInByte := 7 - (offset % 8)
		mask = 1 << bitInByte
	case LSBToMSB:
		byteIdx = nBytes - 1 - offset/8
		bitInByte := offset % 8
		mask = 1 << bitInByte
	default:
		return 0, 0, fmt.Errorf("unknown bit order %d", order)
	}
	return byteIdx, mask, nil
}

func Bit(value []byte, index int, order BitOrder, start StartIndex) (byte, error) {
	byteIdx, mask, err := localize(len(value), index, order, start)
	if err != nil {
		return 0, err
	}
	if value[byteIdx]&mask != 0 {
		return 1, nil
	}
	return 0, nil
}

func SetBit(value []byte, index int, bit byte, order BitOrder, start StartIndex) error {
	byteIdx, mask, err := localize(len(value), index, order, start)
	if err != nil {
		return err
	}
	if bit != 0 {
		value[byteIdx] |= mask
	} else {
		value[byteIdx] &^= mask
	}
	return nil
}

func Permute(value []byte, pblock []int, order BitOrder, start StartIndex) ([]byte, error) {
	if len(pblock) == 0 {
		return []byte{}, nil
	}
	outBits := len(pblock)
	out := make([]byte, (outBits+7)/8)
	for i, src := range pblock {
		bit, err := Bit(value, src, order, start)
		if err != nil {
			return nil, fmt.Errorf("p-block[%d]=%d: %w", i, src, err)
		}
		if err := SetBit(out, int(start)+i, bit, order, start); err != nil {
			return nil, err
		}
	}
	return out, nil
}
