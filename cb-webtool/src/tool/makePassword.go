package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// 평문의 패스워드에서 단방향 해시를 생성
// 패스워드는 72글자 이내로 해야함. (초반 72글자까지 일치하면 일치한 것으로 간주 됨.)
func GeneratePassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// 평문 패스워드와 단방향 해시가 일치한지 판단.
func CheckPassword(password string, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, err
		}
		return false, err
	}
	return true, nil
}

func main() {
	//$2a$10$hw32uwvJVxeE4YQNk31DT.bEN2UOs7/8UrvpVahyGJboGHLU4uuES
	//$2a$10$yQ8jUHlwnm1r/cqsxLxq7ushC1MaT6CYY92BZyrghTM05bqiv5riW

	if len(os.Args) < 2 {
		fmt.Printf("go run %s password\n", "makePassword")
		return
	}

	inputPassword := os.Args[1:2][0]

	//inputPassword := "admin"
	hash, err := GeneratePassword(inputPassword)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hash)
	/*
		hash := "$2a$10$hw32uwvJVxeE4YQNk31DT.bEN2UOs7/8UrvpVahyGJboGHLU4uuES"
		isOk, err := CheckPassword(inputPassword, hash)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(isOk)
	*/
}
