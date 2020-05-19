package certutil

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadOriginCert(t *testing.T) {
	cert, err := DecodeOriginCert([]byte{})
	assert.Equal(t, fmt.Errorf("Cannot decode empty certificate"), err)
	assert.Nil(t, cert)

	blocks, err := ioutil.ReadFile("test-cert-no-key.pem")
	assert.Nil(t, err)
	cert, err = DecodeOriginCert(blocks)
	assert.Equal(t, fmt.Errorf("Missing private key in the certificate"), err)
	assert.Nil(t, cert)

	blocks, err = ioutil.ReadFile("test-cert-two-certificates.pem")
	assert.Nil(t, err)
	cert, err = DecodeOriginCert(blocks)
	assert.Equal(t, fmt.Errorf("Found multiple certificates in the certificate"), err)
	assert.Nil(t, cert)

	blocks, err = ioutil.ReadFile("test-cert-unknown-block.pem")
	assert.Nil(t, err)
	cert, err = DecodeOriginCert(blocks)
	assert.Equal(t, fmt.Errorf("Unknown block RSA PRIVATE KEY in the certificate"), err)
	assert.Nil(t, cert)

	blocks, err = ioutil.ReadFile("test-cert.pem")
	assert.Nil(t, err)
	cert, err = DecodeOriginCert(blocks)
	assert.Nil(t, err)
	assert.NotNil(t, cert)
	assert.Equal(t, "7b0a4d77dfb881c1a3b7d61ea9443e19", cert.ZoneID)
	key := "v1.0-58bd4f9e28f7b3c28e05a35ff3e80ab4fd9644ef3fece537eb0d12e2e9258217-183442fbb0bbdb3e571558fec9b5589ebd77aafc87498ee3f09f64a4ad79ffe8791edbae08b36c1d8f1d70a8670de56922dff92b15d214a524f4ebfa1958859e-7ce80f79921312a6022c5d25e2d380f82ceaefe3fbdc43dd13b080e3ef1e26f7"
	assert.Equal(t, key, cert.ServiceKey)
}

func TestNewlineArgoTunnelToken(t *testing.T) {
	ArgoTunnelTokenTest(t, "test-argo-tunnel-cert.pem")
}

func TestJSONArgoTunnelToken(t *testing.T) {
	// The given cert's Argo Tunnel Token was generated by base64 encoding this JSON:
	// {
	// "zoneID": "7b0a4d77dfb881c1a3b7d61ea9443e19",
	// "serviceKey": "test-service-key",
	// "accountID": "abcdabcdabcdabcd1234567890abcdef"
	// }
	ArgoTunnelTokenTest(t, "test-argo-tunnel-cert-json.pem")
}

func ArgoTunnelTokenTest(t *testing.T, path string) {
	blocks, err := ioutil.ReadFile(path)
	assert.Nil(t, err)
	cert, err := DecodeOriginCert(blocks)
	assert.Nil(t, err)
	assert.NotNil(t, cert)
	assert.Equal(t, "7b0a4d77dfb881c1a3b7d61ea9443e19", cert.ZoneID)
	key := "test-service-key"
	assert.Equal(t, key, cert.ServiceKey)
}
