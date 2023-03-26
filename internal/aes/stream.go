package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"io"

	"github.com/KirillMironov/beaver/internal/rand"
)

type Encrypter struct {
	src io.Reader
	dst io.Writer
}

func NewEncrypter(src io.Reader, dst io.Writer) *Encrypter {
	return &Encrypter{
		src: src,
		dst: dst,
	}
}

// Encrypt encrypts data from src and writes it to dst.
func (e Encrypter) Encrypt(key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv, err := rand.Key(block.BlockSize())
	if err != nil {
		return err
	}

	_, err = e.dst.Write(iv)
	if err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	writer := cipher.StreamWriter{
		S: stream,
		W: e.dst,
	}

	_, err = io.Copy(writer, e.src)

	return err
}

type Decrypter struct {
	src io.Reader
	dst io.Writer
}

func NewDecrypter(src io.Reader, dst io.Writer) *Decrypter {
	return &Decrypter{
		src: src,
		dst: dst,
	}
}

// Decrypt decrypts data from src and writes it to dst.
func (d Decrypter) Decrypt(key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, block.BlockSize())

	_, err = d.src.Read(iv)
	if err != nil {
		return err
	}

	stream := cipher.NewCFBDecrypter(block, iv)

	reader := cipher.StreamReader{
		S: stream,
		R: d.src,
	}

	_, err = io.Copy(d.dst, reader)

	return err
}
