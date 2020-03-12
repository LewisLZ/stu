package ut

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/consts"
)

const (
	ctxUser = "user"
	ctxSrc  = "src"
	CtxDS   = "CTX_DS"
)

type SessionKey struct {
	Key   interface{}
	Value interface{}
}

type SessionOpt struct {
	Keys []SessionKey
}

func SetWebSession(c *gin.Context, opt SessionOpt) error {
	// set session
	session := sessions.Default(c)
	for _, key := range opt.Keys {
		session.Set(key.Key, key.Value)
	}

	if err := session.Save(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func GetWebSession(c *gin.Context, opt SessionOpt) (map[interface{}]interface{}, error) {
	session := sessions.Default(c)

	m := make(map[interface{}]interface{})
	for _, k := range opt.Keys {
		m[k] = session.Get(k.Key)
	}

	return m, nil
}

func DelWebSession(c *gin.Context, opt SessionOpt) error {
	session := sessions.Default(c)
	for _, k := range opt.Keys {
		session.Delete(k.Key)
	}

	if err := session.Save(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

var Passwd = utee.Md5Str("s_pe?ltZ8")

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 3DES加密
func TripleDesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	origData = pkcs5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// 3DES解密
func TripleDesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = pkcs5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

type TokenInfo struct {
	Uid    uint  `json:"uid"`
	Expire int64 `json:"expire"` // 到期tick
}

var desKey = []byte("lots6??Mojavlots6??Mojav") // 3des要求， 必须24位

func Info2Token(info *TokenInfo) (string, error) {
	b, err := json.Marshal(info)
	if err != nil {
		return "", err
	}

	result, err := TripleDesEncrypt(b, desKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(result), nil
}

func Token2Info(token string) (*TokenInfo, error) {
	b, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	result, err := TripleDesDecrypt(b, desKey)
	if err != nil {
		return nil, err
	}

	var info *TokenInfo
	err = json.Unmarshal(result, &info)

	return info, err
}

func TickToTime(tick int64) time.Time {
	return time.Unix(0, tick*int64(time.Millisecond))
}

func UtickToTime(tick uint64) time.Time {
	return TickToTime(int64(tick))
}

func IsDup(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Error 1062")
}

func Int32ToInf(src []int32) []interface{} {
	result := []interface{}{}
	for _, v := range src {
		result = append(result, v)
	}
	return result
}

const RMB_UNIT = 100 // 以分为单位

func FmtPrice(price int) float64 {
	return float64(price) / float64(RMB_UNIT)
}

func FmtInt64Price(price int64) float64 {
	return float64(price) / float64(RMB_UNIT)
}

func ParsePriceToInt(s string, times ...int) (int, error) {
	if s == "" {
		return 0, nil
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	t := float64(100)
	if len(times) > 0 {
		t = float64(times[0])
	}

	return int(math.Floor(f*t + 0.5)), nil
}

func ParsePriceToInt64(s string, times ...int) (int64, error) {
	if s == "" {
		return 0, nil
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	t := float64(100)
	if len(times) > 0 {
		t = float64(times[0])
	}

	return int64(math.Floor(f*t + 0.5)), nil
}

func truncatePrice(f float64) string {
	price := fmt.Sprintf("%.2f", f)
	idx := strings.LastIndex(price, ".")
	if idx == -1 {
		return price
	}

	if idx == len(price)-1 {
		pr1 := price[:idx]
		if pr1 == "" {
			return "0"
		}
		return pr1
	}

	pr1 := price[0:idx]
	pr2 := price[idx+1:]

	pr3, err := strconv.Atoi(pr2)
	if err != nil {
		return price
	}

	if pr3 != 0 {
		return price
	}

	if pr1 == "" {
		return "0"
	}

	return pr1
}

func TruncatePriceFloat(f float64) string {
	return truncatePrice(f)
}

func TruncatePriceInt(p int) string {
	f := FmtPrice(p)
	return truncatePrice(f)
}

func NowTimeFormat() string {
	return time.Now().Format(consts.DateFmt)
}

func FmtTime(time time.Time) string {
	return time.Format(consts.DateFmt)
}

func GetTime(t string) (time.Time, error) {
	return time.ParseInLocation(consts.DateFmt, t, time.Local)
}

func JsonString(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(b)
}

// 打印当前函数名
func GetCurrentFuncName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "???"
	}
	return runtime.FuncForPC(pc).Name()
}

func UintSlice2String(s []uint, split string) string {
	if len(s) == 0 {
		return ""
	}
	res := make([]string, len(s))
	for k, v := range s {
		res[k] = strconv.Itoa(int(v))
	}

	return strings.Join(res, split)
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func MustInterfaceToString(itf interface{}) string {
	if itf == nil {
		return ""
	}
	return itf.(string)
}
func MustInterfaceToFloat64(itf interface{}) float64 {
	if itf == nil {
		return 0.0
	}
	return itf.(float64)
}

// 判断是上半月还是下半月 true :上半月 false :下半月
func IsFirstHalfOfMonth(t time.Time) bool {
	return t.Day() < 16
}

// 获取传入的时间所在月份的月中一天，即某月16号的0点。如传入time.Now(), 返回当前月份16号0点时间。
func GetMidDateOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 16, 0, 0, 0, 0, t.Location())
}

func I64toa(price int64) string {
	return strconv.FormatFloat(float64(price)/100, 'f', 2, 64)
}

// 获取分页请求时的当前页码，每页条数，偏移量
func MakePager(page, limit, defaultLimit int) (int, int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = defaultLimit
	}
	offset := (page - 1) * limit

	return page, limit, offset
}
