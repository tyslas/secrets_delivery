package main

import (
	"fmt"
)

func main() {
	privateKey := createRSAKey()
	publicKey := privateKey.PublicKey

	encryptedBytes := encryptBytes256(publicKey, []byte("super secret message"))

	//msg := []byte("verifiable message")

	msgHashSum := makeHashSum(encryptedBytes)
	signature := signHash(privateKey, msgHashSum)

	verifySignature(publicKey, msgHashSum, signature)

	// We get back the original information in the form of bytes, which we
	// then cast to a string and print
	decryptedBytes := decryptBytes(privateKey, encryptedBytes)
	fmt.Println("decrypted message: ", string(decryptedBytes))
}
