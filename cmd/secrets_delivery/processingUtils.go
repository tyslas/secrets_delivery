package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func checkFlags(flags int) {
	if flag.NFlag() < 1 {
		log.Fatal("at least one flag is missing from: prKey, pbKey, or msg")
	}
}

func checkPath(path string) {
	_, err := os.Stat(path)
	if err != nil {
		fmt.Sprintf("Path does not exist for file: %s", path)
		log.Fatal(err)
	}
}

func readFile(path string) []byte {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	//fmt.Println(string(body))
	return body
}

func readKeyFiles(prKeypath string, pbKeyPath string) (rsa.PrivateKey, rsa.PublicKey, error) {
	var privateKey *rsa.PrivateKey
	var publicKey *rsa.PublicKey

	priv, err := ioutil.ReadFile(prKeypath)
	if err != nil {
		fmt.Println("No RSA private key found")
		return *privateKey, *publicKey, err
	}

	privPem, _ := pem.Decode(priv)
	var privPemBytes []byte
	if privPem.Type != "RSA PRIVATE KEY" {
		fmt.Println("RSA private key is of the wrong type")
		return *privateKey, *publicKey, err
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil { // note this returns type `interface{}`
			fmt.Println("Unable to parse RSA private key")
			return *privateKey, *publicKey, err
		}
	}

	var ok bool
	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		fmt.Println("Unable to parse RSA private key")
		return *privateKey, *publicKey, err
	}

	pub, err := ioutil.ReadFile(pbKeyPath)
	if err != nil {
		fmt.Println("No RSA public key found")
		return *privateKey, *publicKey, err
	}
	pubPem, _ := pem.Decode(pub)
	if pubPem == nil {
		fmt.Println("Use `ssh-keygen -f id_rsa.pub -e -m pem > id_rsa.pem` to generate the pem encoding of your RSA public key")
		return *privateKey, *publicKey, err
	}
	if pubPem.Type != "RSA PUBLIC KEY" {
		fmt.Println("RSA public key is of the wrong type")
	}

	if parsedKey, err = x509.ParsePKIXPublicKey(pubPem.Bytes); err != nil {
		fmt.Println("Unable to parse RSA public key")
		return *privateKey, *publicKey, err
	}

	if publicKey, ok = parsedKey.(*rsa.PublicKey); !ok {
		fmt.Println("Unable to parse RSA public key")
		return *privateKey, *publicKey, err
	}

	return *privateKey, *publicKey, nil
}

func writeFile(signedHash []byte, path string) {
	destination, err := os.Create(path)
	if err != nil {
		fmt.Println("os.Create:", err)
		return
	}
	defer destination.Close()

	fmt.Fprint(destination, signedHash)
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := asn1.Marshal(pubkey)
	checkError(err)

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(fileName)
	checkError(err)
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	checkError(err)
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

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pub, nil
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
