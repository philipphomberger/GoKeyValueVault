package cryptostore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/thanhpk/randstr"
	"io"
	"io/ioutil"
	"log"
)

func readKey() []byte {
	key, err := ioutil.ReadFile("assets/key.txt")
	if err != nil {
		log.Fatalf("read file err: %v", err.Error())
	}
	return key
}

func createBlockOfAlgorithm(key []byte) cipher.Block {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("cipher err: %v", err.Error())
	}
	return block
}

func createGcmMode(block cipher.Block) cipher.AEAD {
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("cipher GCM err: %v", err.Error())
	}
	return gcm
}

func generatingRandomNonce(gcm cipher.AEAD) []byte {
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf("nonce  err: %v", err.Error())
	}
	return nonce
}

func GenerateKey() {
	err := ioutil.WriteFile("assets/key.txt", []byte(randstr.String(32)), 0777)
	if err != nil {
		log.Fatalf("write file err: %v", err.Error())
	}
}

func EncryptText(plainText []byte) []byte {
	// Reading key
	key := readKey()
	// Creating block of algorithm
	block := createBlockOfAlgorithm(key)
	// Creating GCM mode
	gcm := createGcmMode(block)
	// Generating random nonce
	nonce := generatingRandomNonce(gcm)
	// Decrypt file
	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return cipherText
}

func DecryptText(encryptText []byte) string {
	// Reading key
	key := readKey()
	// Creating block of algorithm
	block := createBlockOfAlgorithm(key)
	// Creating GCM mode
	gcm := createGcmMode(block)

	nonce := encryptText[:gcm.NonceSize()]
	encryptText = encryptText[gcm.NonceSize():]
	plainText, err := gcm.Open(nil, []byte(nonce), []byte(encryptText), nil)
	if err != nil {
		log.Fatalf("decrypt file err: %v", err.Error())
	}
	fmt.Println(string(plainText))
	return string(plainText)
}
