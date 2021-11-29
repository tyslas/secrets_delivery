package main

import (
	"flag"
	"fmt"
)

func main() {
	privateKeyPtr := flag.String("prKey", "", "location of private key")
	publicKeyPtr := flag.String("pbKey", "", "location of public key")
	messagePtr := flag.String("msg", "", "location of message to encrypt/decrypt/verify")
	//encryptDecryptPtr := flag.Bool("enc", false, "including this flag will encrypt, omitting it will decrypt")

	flag.Parse()

	//checkFlags(flag.NFlag())
	checkPath(*privateKeyPtr)
	checkPath(*publicKeyPtr)
	checkPath(*messagePtr)

	privateKey := createRSAKey()
	publicKey := privateKey.PublicKey
	savePEMKey("/Users/titoyslas/.ssh/go.pem", &privateKey)
	savePublicPEMKey("/Users/titoyslas/.ssh/go.pub", publicKey)

	// once private and public key paths are verified
	// read the files and return the rsa.PrivateKey
	//privateKey, publicKey, err := readKeyFiles(*privateKeyPtr, *publicKeyPtr)
	//if err != nil {
	//	log.Fatal(err)
	//}

	// read the message file and return slice of bytes
	msgBytes := readFile(*messagePtr)

	// encrypt bytes with recipient's public key and sign with sender's private key
	var encryptedBytes []byte
	//var signedHash []byte
	//if *encryptDecryptPtr == true {
	encryptedBytes = encryptBytes256(publicKey, msgBytes)
	msgHashSum := makeHashSum(encryptedBytes)
	signature := signHash(privateKey, msgHashSum)
	//}

	writeFile(msgHashSum, "/Users/titoyslas/workspace/vandy/cyberSecurity-cs6387/final_project/hashSum")
	writeFile(signature, "/Users/titoyslas/workspace/vandy/cyberSecurity-cs6387/final_project/signature")

	//msgHashSum := readFile("/Users/titoyslas/workspace/vandy/cyberSecurity-cs6387/final_project/hashSum")
	//signature := readFile("/Users/titoyslas/workspace/vandy/cyberSecurity-cs6387/final_project/signature")

	//publicKeyFile := string(readFile(*publicKeyPtr))
	//publicKeyParsed, err := ParseRsaPublicKeyFromPemStr(publicKeyFile)
	//if err != nil {
	//	fmt.Println("cannot parse rsa public key")
	//	log.Fatal(err)
	//}

	//privateKeyFile := string(readFile(*privateKeyPtr))
	//privateKeyParsed, err := ParseRsaPrivateKeyFromPemStr(privateKeyFile)
	//if err != nil {
	//	fmt.Println("cannot parse rsa private key")
	//	log.Fatal(err)
	//}

	//fmt.Println("right before decrypt check")
	//if *encryptDecryptPtr != true {
	fmt.Println("decrypt")
	verifySignature(publicKey, msgHashSum, signature)
	//decryptedBytes := decryptBytes256(*privateKeyParsed, msgHashSum)
	//fmt.Println("decrypted message: ", string(decryptedBytes))
	//}

}
