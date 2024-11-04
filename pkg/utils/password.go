package utils

import "golang.org/x/crypto/bcrypt"

// EncodePWD password encryption
func EncodePWD(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	return string(hash), err
}

// ComparePwd Compares the password ciphertext and plaintext for equality
func ComparePwd(pwd string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	if  err != nil {
		return false
	}  else {
		return true
	}
}
