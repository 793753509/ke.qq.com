package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"ke.qq.com/scrape"
	"ke.qq.com/storage"
	"ke.qq.com/web"
	"time"
)

const (
	user     = "root"
	pwd      = "Uv38ByGCZU8WP18PrLoHNrdJdiO8yvs22UPbKhVYnCodYJXupr2PufFU8dc="
	key      = "12345678901234587690123458769012"
	ip       = "localhost"
	port     = "3306"
	database = "ke_qq_com"
)

func main() {
	storage.InitDbService(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, AesGCMDecrypt(pwd, key), ip, port, database))
	go scrapeCoursesByDay()
	web.NewApiserver().Server()

}

func scrapeCoursesByDay() {
	ticker := time.NewTimer(24 * time.Hour)
	getAndSaveCourses()
	for {
		select {
		case <-ticker.C:
			getAndSaveCourses()
		}
	}
}

func getAndSaveCourses() {
	tableName := storage.CreateTable()
	scraper, curseChan := scrape.NewScrapeManager()
	go storage.SyncCourseInfo(curseChan, tableName)
	scraper.ScrapeCourseInfo(tableName)

}

func AesGCMDecrypt(cryted string, key string) string {
	data, err := base64.StdEncoding.DecodeString(cryted)
	if err != nil {
		fmt.Println("decode str err: ", err.Error())
		return ""
	}
	k := []byte(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		fmt.Println("aes.NewCipher err: ", err.Error())
		return ""
	}

	blockMode, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("cipher.NewGCM err: ", err.Error())
		return ""

	}

	nonceSize := blockMode.NonceSize()
	if len(data) < nonceSize {
		fmt.Println("len data < noncesize")
	}
	nonce, out := data[:nonceSize], data[nonceSize:]
	out, err = blockMode.Open(nil, nonce, out, nil)
	if err != nil {
		fmt.Println("blockMode.Open err: ", err.Error())
		return ""
	}
	out = PKCS7UnPadding(out)
	return string(out)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]

}
