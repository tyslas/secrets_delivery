package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

func createRSAKey() rsa.PrivateKey {
	// The GenerateKey method takes in a reader that returns random bits, and
	// the number of bits
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	return *privateKey
}

func encryptBytes256(pubKey rsa.PublicKey, bytes []byte) []byte {
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		&pubKey,
		bytes,
		nil)
	if err != nil {
		panic(err)
	}

	return encryptedBytes
}

func decryptBytes(privateKey rsa.PrivateKey, encryptedBytes []byte) []byte {
	// The first argument is an optional random data generator (the rand.Reader we used before)
	// we can set this value as nil
	// The OAEPOptions in the end signify that we encrypted the data using OAEP, and that we used
	// SHA256 to hash the input.
	decryptedBytes, err := privateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		panic(err)
	}

	return decryptedBytes
}

func makeHashSum(msg []byte) []byte {
	msgHash := sha256.New()
	_, err := msgHash.Write(msg)
	if err != nil {
		panic(err)
	}

	return msgHash.Sum(nil)
}

// In order to generate the signature, we provide a random number generator,
// our private key, the hashing algorithm that we used, and the hash sum
// of our message
func signHash(privateKey rsa.PrivateKey, msgHashSum []byte) []byte {
	signature, err := rsa.SignPSS(rand.Reader, &privateKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		panic(err)
	}

	return signature
}

func verifySignature(publicKey rsa.PublicKey, msgHashSum []byte, signature []byte) {
	err := rsa.VerifyPSS(&publicKey, crypto.SHA256, msgHashSum, signature, nil)
	if err != nil {
		fmt.Println("could not verify signature: ", err)
		return
	}

	// If we don't get any error from the `VerifyPSS` method, that means our
	// signature is valid
	fmt.Println("signature verified")
}