package nsqplus

import (
	"crypto/md5"
	"encoding/hex"
)

/**
	md5Ctx := md5.New()
	//io.WriteString(md5Ctx, key)
	md5Ctx.Write([]byte(source))
	md5Ctx.Write([]byte(key))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
 */
func GetMD5Digest(source string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(source))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
