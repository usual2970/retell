package encode

import "encoding/base64"

func Base64(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}
