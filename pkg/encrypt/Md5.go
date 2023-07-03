package encrypt

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
)

func main() {
	password := "abcde"
	md5Hash := GetMD5Hash(password)

	fmt.Println("password:", password)
	fmt.Println("password result:", md5Hash)
}

func GetMD5Hash(text string) string {
	hash := md5.New()
	_, err := hash.Write([]byte(text))
	if err != nil {
		log.Fatal(err)
	}

	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
