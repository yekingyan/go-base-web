package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"

	utils "gService/share/utils"
)

const privatekey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAle+EG1sX1mNxo1sl/ZJ9ejcOBUkRoe7OThD2g6UATls7GrKj
tZv/iBf49+s2aIz9ZlZVO2vckFkusvd9e/4Dk/0Kkt3pdzzMpkOEthHSTJhHUGDS
FUKz9fYPBbYgf/tUh6bAGNzoaKskAXTSHmEGm7tv78nfJOE1IFXgH91u44oPdY/R
GUzHzoLgr/oPxNqisi0vmePqThTAFEC7aYTKNV/R+VXJLesfc/qs9JQmMppyJdac
PJqoq2gLoOdaZPfXlxSjJ8IDSbTVOoMbwZOdM1MUDjQgclbR7qi7cmkRSTdMOWLP
Qq2Cqce3QUl9t8+hHc63DFHPkmMkhPhRVLkdEwIDAQABAoIBAQCUSmoMbFHtNjmE
aXfvjV0Q7Tauv4/y6q+cwmYt2Zmt5clVJG2KBxn6TkttWHNdgGR3PNfbh1iEY9Au
LV0srxR234P/sf9gcP+vSYQOVx35p3qQ4tGlOW2nzI8bBQCr07XM3H2NmZ87NljM
H/BZH3lbuUyCitRqM1tJ+eVLXsR3wdw458YeHXiuY/nOgIsB984o1k9nS8kmELsA
hUh0jT78pUaRSTPAiRHw5o+cmcYAZ16y/gDYO9ev7lsxWmMk71ljEKsaJhr1Ox5F
F6SRVqpCJmTVD66KdE6/WMY3Kr03a8mz9+4KYZ5NvMYmgna4JYwIylhy1UtQdWcw
o2UwvorxAoGBAPe9tGPoUH65xY7l1vxjiBTfmU5M00Riy8mUUEnmQLGsSzoySe4o
0zYsS0QbkHntc7RomxPuuEs6X9VkGVVSOtyIGVGEJ9AouLXjNHpdOxil809a86RZ
RwG+EGXGXb+KkfgXuLNMTVO0ZuLxvzY0x2cuDcwYO8tJhsW23Db2Cst1AoGBAJrv
HGZmC4L9SuavQ1pl/tjsklCr0TWMRJcrUNDUxPx+isYGMTCxpQqhQBwsz07M0eCe
szIV1HaHn4eDC3DXrNVzFgZfsISeMzxqpPPWmbvdsmV2sBuQLQr3y4GhZH11FgMM
fxtYVUOmDUQNqOQjtp7ghdH20fE1m6emYBMuYR1nAoGBAOAq+uBwJO4WYOnllTu5
QDhK/yh3oa8+ilGzb1b/Dsj2MvfM61KiBEP/nndZTcjWEJ1NDg3iM0Z24qJvOfEX
QAWu7OHy6CeVwVVr6l9Snxe1iczjapTq57Ju7d15ufiIhX5s0IcE4u97zKFLyA5f
gZMefAOSZgCTXmznAqqExc3FAoGAemlZgv41MjV2LHVlPdNAW5vurpZSIYDF4Lp0
i8rQKL7CXlseGl8BCzkEMj6lPJlPaa1536SnzU6ymJrNO0bsY2keicKo8N8dlCqe
UZnItUogXVI9KknrjLLjs3QUtZsA3T/OXYiZNW3JJW+1dStSaozyrkXY8j9s0DT/
y6flSxkCgYEA41qcXKsGet39u7DpZjkCeCkO0nGj8mXBh/KL1xueFAXdlyq+01CG
iQtJD/+/FdK7HUZl5HsD8BhZeEhXqG8ZlRLreZ8U+KOpLrJvqU3hBh27IYaFOVSE
7McmIiUuplKnThuWQQr+CNuwkyFRLS6XEFI4+sN5TqyNu/L4VD0Lb9o=
-----END RSA PRIVATE KEY-----`

func TestGenerateToken(t *testing.T) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privatekey))
	if err != nil {
		t.Fatalf("can not parse private key: %v", err)
	}
	gen := NewJWTToken(int64((time.Minute * 2).Seconds()), utils.FakeNow("2022-01-01 8:00:00"), key)
	tk, expire, err := gen.GenerateToken("user1IDasdf")
	if err != nil {
		t.Fatalf("can not generate token: %v", err)
	}
	want := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDA5OTUzMjAsImlhdCI6MTY0MDk5NTIwMCwiaXNzIjoiZ1NlcnZpY2UiLCJzdWIiOiJ1c2VyMUlEYXNkZiJ9.RVbvPcytgCOaP9YWdlpmtO-XupzK1T7QCjC4WGeWJMmSoRqincp-AO6N7ofjva8MvpWr-zL-TdqnxwP9YURmVtNak8i-awrL9QSgvATMhU9BKMmLb3srPXcTeda9d7mnaUpcr2Sp2tQjYmbVo1OukJv7rzxqHEKK3-GQ1YYMpEzzKjucQKqMxBtu_NJQXBGhnUlfhI1bwKdevHqgYfE_5pkEiCrz690mcGOcqZmnfovCfvjvbyZBE9FM4nEpyKBLeLpi0HjvGrmIsl6A71-Uu-qaJlKLkf1G5N6iS3mRvuNkP4uPPP2kdzDaxMgWcjmn1E9r-wsvN72RcL6_MDq6Kw"
	if tk != want {
		t.Errorf("want token: %s, got: %s", want, tk)
	}

	exTime := time.Unix(expire, 0)
	wantExTime := utils.FakeNow("2022-01-01 8:02:00")()
	fmt.Println(expire, exTime)
	if exTime != wantExTime {
		t.Errorf("want expiration time: %s, got: %s", wantExTime, exTime)
	}
}
