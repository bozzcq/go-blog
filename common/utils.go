package common

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

var Log *Logger

func SetLog() {
	var out io.Writer
	if Cfg.MustBool("", "log_file", false) {
		f, _ := os.OpenFile(Cfg.MustValue("", "log_path", "./log.txt"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		out = io.MultiWriter(f)
	} else {
		out = os.Stdout
	}
	Log = New(out, "[go-blog]", Lshortfile|Ldate|Lmicroseconds)
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func ParseInt(value string) int {
	if value == "" {
		return 0
	}
	val, _ := strconv.Atoi(value)
	return val
}

func IntString(value int) string {
	return strconv.Itoa(value)
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Atoa(str string) string {
	var result string
	for i := 0; i < len(str); i++ {
		c := rune(str[i])
		if 'A' <= c && c <= 'Z' && i > 0 {
			result = result + "_" + strings.ToLower(string(str[i]))
		} else {
			result = result + string(str[i])
		}
	}
	return result
}

func GetRemoteIp(r *http.Request) (ip string) {
	ip = r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.RemoteAddr
	}
	ip = strings.Split(ip, ":")[0]
	if len(ip) < 7 || ip == "127.0.0.1" {
		ip = "localhost"
	}
	return
}

/* Test Helpers */
func Expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func GenToken() string {
	nano := time.Now().UnixNano()
	rand.Seed(nano)
	rndNum := rand.Int63()
	uuid := Md5(Md5(strconv.FormatInt(nano, 10)) + Md5(strconv.FormatInt(rndNum, 10)))
	return uuid
}
