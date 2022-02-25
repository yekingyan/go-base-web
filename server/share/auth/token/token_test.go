package sharetoken

import (
	"testing"

	"github.com/dgrijalva/jwt-go"

	utils "gService/share/utils"
)

const publickey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAle+EG1sX1mNxo1sl/ZJ9
ejcOBUkRoe7OThD2g6UATls7GrKjtZv/iBf49+s2aIz9ZlZVO2vckFkusvd9e/4D
k/0Kkt3pdzzMpkOEthHSTJhHUGDSFUKz9fYPBbYgf/tUh6bAGNzoaKskAXTSHmEG
m7tv78nfJOE1IFXgH91u44oPdY/RGUzHzoLgr/oPxNqisi0vmePqThTAFEC7aYTK
NV/R+VXJLesfc/qs9JQmMppyJdacPJqoq2gLoOdaZPfXlxSjJ8IDSbTVOoMbwZOd
M1MUDjQgclbR7qi7cmkRSTdMOWLPQq2Cqce3QUl9t8+hHc63DFHPkmMkhPhRVLkd
EwIDAQAB
-----END PUBLIC KEY-----`

func TestVerify(t *testing.T) {
	jwt.TimeFunc = utils.FakeNow("2022-01-01 8:02:00")
	pk, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publickey))
	if err != nil {
		panic(err)
	}
	verifier := &JWTTokenVerifier{
		PublicKey: pk,
	}
	token := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDA5OTUzMjAsImlhdCI6MTY0MDk5NTIwMCwiaXNzIjoiZ1NlcnZpY2UiLCJzdWIiOiJ1c2VyMUlEYXNkZiJ9.RVbvPcytgCOaP9YWdlpmtO-XupzK1T7QCjC4WGeWJMmSoRqincp-AO6N7ofjva8MvpWr-zL-TdqnxwP9YURmVtNak8i-awrL9QSgvATMhU9BKMmLb3srPXcTeda9d7mnaUpcr2Sp2tQjYmbVo1OukJv7rzxqHEKK3-GQ1YYMpEzzKjucQKqMxBtu_NJQXBGhnUlfhI1bwKdevHqgYfE_5pkEiCrz690mcGOcqZmnfovCfvjvbyZBE9FM4nEpyKBLeLpi0HjvGrmIsl6A71-Uu-qaJlKLkf1G5N6iS3mRvuNkP4uPPP2kdzDaxMgWcjmn1E9r-wsvN72RcL6_MDq6Kw"

	uerid, err := verifier.Verify(token)
	if err != nil {
		t.Errorf("failed to verify token: %v", err)
	}
	if uerid != "user1IDasdf" {
		t.Errorf("token userid dont match: %v", uerid)
	}
}
