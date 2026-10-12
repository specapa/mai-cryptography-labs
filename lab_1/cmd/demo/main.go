package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"lab1/deal"
	"lab1/des"
	"lab1/symmetric"
)

func main() {
	root := findTestdata()
	outDir := filepath.Join(root, "out")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fatal(err)
	}

	desKey := []byte{0x13, 0x34, 0x57, 0x79, 0x9B, 0xBC, 0xDF, 0xF1}
	desIV := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}
	desDelta := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xA5}

	dealKey := append(append([]byte{}, desKey...), desIV...)
	dealIV := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF, 0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10}
	dealDelta := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xA5}

	if err := generateMedia(root); err != nil {
		fatal(err)
	}

	random := make([]byte, 4096)
	if _, err := rand.Read(random); err != nil {
		fatal(err)
	}

	files := []string{
		"sample.txt",
		"rbtree_allocator.go",
		"rbtree_allocator_test.go",
		"tone.wav",
		"image.ppm",
		"clip.y4m",
	}

	fmt.Println("=== Задание 5: DES, псевдослучайные байты ===")
	runModes("des-random-4096", func() symmetric.SymmetricAlgorithm { return des.New() }, random, desKey, desIV, desDelta, outDir)
	fmt.Println("\n=== Задание 5: DES, файлы ===")
	runFiles("des", func() symmetric.SymmetricAlgorithm { return des.New() }, files, root, desKey, desIV, desDelta, outDir)

	fmt.Println("\n=== Задание 7: DEAL-128, псевдослучайные байты ===")
	runModes("deal-random-4096", func() symmetric.SymmetricAlgorithm { return deal.New() }, random, dealKey, dealIV, dealDelta, outDir)
	fmt.Println("\n=== Задание 7: DEAL-128, файлы ===")
	runFiles("deal", func() symmetric.SymmetricAlgorithm { return deal.New() }, files, root, dealKey, dealIV, dealDelta, outDir)

	fmt.Println("\nГотово.")
}

func runModes(label string, newAlg func() symmetric.SymmetricAlgorithm, plain, key, iv, delta []byte, outDir string) {
	modes := []symmetric.CipherMode{
		symmetric.ECB, symmetric.CBC, symmetric.PCBC,
		symmetric.CFB, symmetric.OFB, symmetric.CTR, symmetric.RandomDelta,
	}
	for _, mode := range modes {
		ctx, err := symmetric.NewSymmetricContext(newAlg(), key, mode, symmetric.PKCS7, iv, delta)
		if err != nil {
			fatal(err)
		}
		var enc, dec []byte
		if err := ctx.Encrypt(plain, &enc); err != nil {
			fatal(err)
		}
		if err := ctx.Decrypt(enc, &dec); err != nil {
			fatal(err)
		}
		ok := bytes.Equal(plain, dec)
		mark := "OK"
		if !ok {
			mark = "FAIL"
		}
		_ = os.WriteFile(filepath.Join(outDir, label+"_"+mode.String()+".bin"), enc, 0o644)
		fmt.Printf("  %s PKCS7: %s  cipher=%d байт\n", mode, mark, len(enc))
		if !ok {
			os.Exit(1)
		}
	}
}

func runFiles(prefix string, newAlg func() symmetric.SymmetricAlgorithm, files []string, root string, key, iv, delta []byte, outDir string) {
	for _, name := range files {
		path := filepath.Join(root, name)
		data, err := os.ReadFile(path)
		if err != nil {
			fatal(err)
		}
		fmt.Printf("\n-- %s (%d байт) --\n", name, len(data))
		runFileModes(prefix, newAlg, path, key, iv, delta, outDir)
	}
}

func runFileModes(prefix string, newAlg func() symmetric.SymmetricAlgorithm, inPath string, key, iv, delta []byte, outDir string) {
	base := filepath.Base(inPath)
	modes := []symmetric.CipherMode{
		symmetric.ECB, symmetric.CBC, symmetric.PCBC,
		symmetric.CFB, symmetric.OFB, symmetric.CTR, symmetric.RandomDelta,
	}
	pads := []symmetric.PaddingMode{symmetric.PKCS7, symmetric.ANSIX923}
	for _, mode := range modes {
		for _, pad := range pads {
			ctx, err := symmetric.NewSymmetricContext(newAlg(), key, mode, pad, iv, delta)
			if err != nil {
				fatal(err)
			}
			encPath := filepath.Join(outDir, fmt.Sprintf("%s.%s.%s.%s.enc", prefix, base, mode, padName(pad)))
			decPath := filepath.Join(outDir, fmt.Sprintf("%s.%s.%s.%s.dec", prefix, base, mode, padName(pad)))
			if err := ctx.EncryptFile(inPath, encPath); err != nil {
				fatal(err)
			}
			if err := ctx.DecryptFile(encPath, decPath); err != nil {
				fatal(err)
			}
			want, _ := os.ReadFile(inPath)
			got, _ := os.ReadFile(decPath)
			mark := "OK"
			if !bytes.Equal(want, got) {
				mark = "FAIL"
				fmt.Printf("  %s/%s: %s\n", mode, pad, mark)
				os.Exit(1)
			}
			fmt.Printf("  %s/%s: %s\n", mode, pad, mark)
		}
	}
}

func padName(p symmetric.PaddingMode) string {
	switch p {
	case symmetric.PKCS7:
		return "pkcs7"
	case symmetric.ANSIX923:
		return "x923"
	case symmetric.ISO10126:
		return "iso10126"
	default:
		return "zeros"
	}
}

func generateMedia(root string) error {
	if err := os.WriteFile(filepath.Join(root, "rbtree_allocator_test.go"), rbtreeTestSource, 0o644); err != nil {
		return err
	}
	if err := writeWAV(filepath.Join(root, "tone.wav"), 440, 0.4); err != nil {
		return err
	}
	if err := writePPM(filepath.Join(root, "image.ppm"), 32, 24); err != nil {
		return err
	}
	return writeY4M(filepath.Join(root, "clip.y4m"), 16, 16, 4)
}

func writeWAV(path string, freq float64, seconds float64) error {
	const sampleRate = 8000
	n := int(seconds * sampleRate)
	data := make([]byte, n)
	for i := 0; i < n; i++ {
		s := math.Sin(2 * math.Pi * freq * float64(i) / sampleRate)
		data[i] = byte(128 + 80*s)
	}
	buf := &bytes.Buffer{}
	buf.WriteString("RIFF")
	binary.Write(buf, binary.LittleEndian, uint32(36+len(data)))
	buf.WriteString("WAVEfmt ")
	binary.Write(buf, binary.LittleEndian, uint32(16))
	binary.Write(buf, binary.LittleEndian, uint16(1))
	binary.Write(buf, binary.LittleEndian, uint16(1))
	binary.Write(buf, binary.LittleEndian, uint32(sampleRate))
	binary.Write(buf, binary.LittleEndian, uint32(sampleRate))
	binary.Write(buf, binary.LittleEndian, uint16(1))
	binary.Write(buf, binary.LittleEndian, uint16(8))
	buf.WriteString("data")
	binary.Write(buf, binary.LittleEndian, uint32(len(data)))
	buf.Write(data)
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

func writePPM(path string, w, h int) error {
	var b bytes.Buffer
	fmt.Fprintf(&b, "P6\n%d %d\n255\n", w, h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			b.WriteByte(byte(x * 255 / w))
			b.WriteByte(byte(y * 255 / h))
			b.WriteByte(80)
		}
	}
	return os.WriteFile(path, b.Bytes(), 0o644)
}

func writeY4M(path string, w, h, frames int) error {
	var b bytes.Buffer
	fmt.Fprintf(&b, "YUV4MPEG2 W%d H%d F10:1 Ip A1:1 Cmono\n", w, h)
	frame := make([]byte, w*h)
	for f := 0; f < frames; f++ {
		b.WriteString("FRAME\n")
		for i := range frame {
			frame[i] = byte((i + f*13) % 256)
		}
		b.Write(frame)
	}
	return os.WriteFile(path, b.Bytes(), 0o644)
}

func findTestdata() string {
	candidates := []string{
		"testdata",
		filepath.Join("lab_1", "testdata"),
		filepath.Join("..", "testdata"),
	}
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && st.IsDir() {
			abs, _ := filepath.Abs(c)
			return abs
		}
	}
	fatal(fmt.Errorf("testdata directory not found"))
	return ""
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

var rbtreeTestSource = []byte(`package rbtreealloc

import "testing"

func TestAllocAndLayout(t *testing.T) {
	a := NewAllocator(1024)
	off, ok := a.Alloc(64)
	if !ok {
		t.Fatalf("alloc 64 failed")
	}
	if off != 0 {
		t.Fatalf("first block offset %d, want 0", off)
	}
	off2, ok := a.Alloc(128)
	if !ok {
		t.Fatalf("alloc 128 failed")
	}
	if off2 == off {
		t.Fatalf("overlapping allocations")
	}
}
`)
