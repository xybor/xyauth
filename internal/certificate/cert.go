package certificate

import (
	"io/ioutil"
	"strings"

	"github.com/xybor/xyauth/internal/cloud"
)

// GetCertificateContent returns the certificate content based on the input. The
// input may be the content of certificate, local file name, or a s3 key.
func GetCertificateContent(cert string) ([]byte, error) {
	var certs []byte
	var err error
	if strings.HasPrefix(cert, "s3://") {
		certs, err = cloud.ReadS3(cert)
	} else if strings.Contains(cert, "BEGIN ") {
		certs = []byte(strings.ReplaceAll(cert, `\n`, "\n"))
	} else {
		certs, err = ioutil.ReadFile(cert)
	}

	if err != nil {
		return nil, err
	}

	return certs, nil
}
