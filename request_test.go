package digestsig_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/MadAppGang/digestsig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifySignature(t *testing.T) {
	secret := []byte("the most secret secret")
	body := "my request"
	request, _ := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		"http://google.com/whatever",
		strings.NewReader(body),
	)
	request.Header.Add(digestsig.ContentTypeHeaderKey, "application/json")
	bodyMD5 := digestsig.GetMD5([]byte(body))
	err := digestsig.AddHeadersAndSignRequest(request, secret, bodyMD5)

	require.NoError(t, err)

	assert.Equal(t, bodyMD5, request.Header["Content-Md5"][0])
	assert.NotEmpty(t, request.Header["Expires"][0])
	assert.NotEmpty(t, request.Header["Date"][0])
	assert.NotEmpty(t, request.Header["Digest"][0])

	err = digestsig.VerifySignature(request, secret)
	assert.NoError(t, err)
}
