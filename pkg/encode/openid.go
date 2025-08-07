package encode

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

var blockSize = 20

const encodeSq = "A49Gud_gY8BOk0xfqMTt1Cj3lezZIKsJLyniRvHUoE7ra-hSwWmbcFD5pQ2NP6VX"

func padding(platForm, data []byte) []byte {
	padding := blockSize - len(data)%blockSize
	padText := strings.Repeat(strconv.Itoa(padding), padding)
	return append(append(platForm, data...), padText...)
}

func unPadding(origData []byte) []byte {
	length := len(origData)
	str := string(origData[length-1])
	unPadding, _ := strconv.Atoi(str)
	return origData[:(length - unPadding)]
}

func Openid(platformId int8, identifier string) string {
	return base64.NewEncoding(encodeSq).EncodeToString(padding([]byte(fmt.Sprintf("%04d", platformId)), []byte(identifier)))
}

func FromOpenid(openid string) (platformId int, identifier string) {
	bts, err := base64.NewEncoding(encodeSq).DecodeString(openid)
	if err != nil {
		return 0, ""
	}
	bts = unPadding(bts)

	platformId, _ = strconv.Atoi(string(bts[:4]))
	identifier = string(bts[4:])
	return
}
