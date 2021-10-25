package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var privateKeyFilePath = filepath.Join("keys", "key.pem")
var publicKeyFilePath = filepath.Join("keys", "key.pub")

func getRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, privErr := readRsaPrivateKeyFromFile()
	publicKey, pubErr := readRsaPublicKeyFromFile()
	if privErr != nil || pubErr != nil {
		fmt.Println("Generating key")
		generatedKey, _ := rsa.GenerateKey(rand.Reader, 4096)

		createKeysFolder()
		privWriteErr := writeRsaPrivateKeyToFile(generatedKey)
		pubWriteErr := writeRsaPublicKeyToFile(&generatedKey.PublicKey)

		if privWriteErr != nil || pubWriteErr != nil {
			panic("Unable to write keys to file")
		}

		return generatedKey, &generatedKey.PublicKey
	} else {
		return privateKey, publicKey
	}
}

func createKeysFolder() {
	_, err := os.Stat("keys")

	if os.IsNotExist(err) {
		errDir := os.MkdirAll("keys", 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}
}

func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
	privkeyBytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkeyBytes,
		},
	)
	return string(privkeyPem)
}

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func writeRsaPrivateKeyToFile(privateKey *rsa.PrivateKey) error {
	err := ioutil.WriteFile(privateKeyFilePath, []byte(ExportRsaPrivateKeyAsPemStr(privateKey)), 0644)
	return err
}

func readRsaPrivateKeyFromFile() (*rsa.PrivateKey, error) {
	bytes, err := ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		return nil, err
	}

	return ParseRsaPrivateKeyFromPemStr(string(bytes))
}

func exportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) string {
	pubkeyBytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return ""
	}
	pubkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkeyBytes,
		},
	)

	return string(pubkeyPem)
}

func getRsaPublicKeyAsBase64Str(pubkey *rsa.PublicKey) string {
	return base64.StdEncoding.EncodeToString([]byte(exportRsaPublicKeyAsPemStr(pubkey)))
}

func parseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("key type is not RSA")
}

func writeRsaPublicKeyToFile(privateKey *rsa.PublicKey) error {
	err := ioutil.WriteFile(publicKeyFilePath, []byte(exportRsaPublicKeyAsPemStr(privateKey)), 0644)
	return err
}

func readRsaPublicKeyFromFile() (*rsa.PublicKey, error) {
	bytes, err := ioutil.ReadFile(publicKeyFilePath)
	if err != nil {
		return nil, err
	}

	return parseRsaPublicKeyFromPemStr(string(bytes))
}
