package tools

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
)

func NewCFBDecrypter(s string) (string, error) {

	key, _ := hex.DecodeString("346c6c796f75726234736534726562656c6f6e67746f7573") //6368616e676520746869732070617373

	ciphertext, _ := hex.DecodeString(s)

	block, err := aes.NewCipher(key)

	if err != nil {

		log.Println(err)
		return "", err

	}

	// The IV needs to be unique, but not secure. Therefore it's common to include it at the beginning of the ciphertext.

	if len(ciphertext) < aes.BlockSize {

		log.Println(err)
		return "", err

	}

	iv := ciphertext[:aes.BlockSize]

	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.

	stream.XORKeyStream(ciphertext, ciphertext)

	fmt.Printf("%s", ciphertext)

	// Output: some plaintext
	return string(ciphertext), nil
}

func NewCFBEncrypter(s string) (string, error) {

	key, _ := hex.DecodeString("346c6c796f75726234736534726562656c6f6e67746f7573")

	plaintext := []byte(s)

	block, err := aes.NewCipher(key)

	if err != nil {

		log.Println(err)
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {

		log.Println(err)
		return "", err

	}

	stream := cipher.NewCFBEncrypter(block, iv)

	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated(i.e. by using crypto/hmac) as well as being encrypted in order to be secure.

	var cipherstring string = fmt.Sprintf("%x\n", ciphertext)

	

	return cipherstring, nil

}
