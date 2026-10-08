package symmetric

type KeyExpander interface {
	ExpandKey(key []byte) ([][]byte, error)
}

type RoundTransform interface {
	Transform(block, roundKey []byte) ([]byte, error)
}

type SymmetricAlgorithm interface {
	BlockSize() int
	SetKey(key []byte) error
	EncryptBlock(block []byte) ([]byte, error)
	DecryptBlock(block []byte) ([]byte, error)
}

type CipherMode int

const (
	ECB CipherMode = iota
	CBC
	PCBC
	CFB
	OFB
	CTR
	RandomDelta
)

func (m CipherMode) String() string {
	switch m {
	case ECB:
		return "ECB"
	case CBC:
		return "CBC"
	case PCBC:
		return "PCBC"
	case CFB:
		return "CFB"
	case OFB:
		return "OFB"
	case CTR:
		return "CTR"
	case RandomDelta:
		return "RandomDelta"
	default:
		return "unknown"
	}
}

type PaddingMode int

const (
	Zeros PaddingMode = iota
	ANSIX923
	PKCS7
	ISO10126
)

func (p PaddingMode) String() string {
	switch p {
	case Zeros:
		return "Zeros"
	case ANSIX923:
		return "ANSI X.923"
	case PKCS7:
		return "PKCS7"
	case ISO10126:
		return "ISO 10126"
	default:
		return "unknown"
	}
}

func (m CipherMode) needsIV() bool {
	return m != ECB
}
