package whitelist

import (
	"crypto/x509"
	"strings"
	"unicode/utf8"

	"github.com/adamdecaf/cert-manage/pkg/certutil"
)

const (
	minimumFingerprintLength = 8
)

// fingerprint matches an incoming certificate's fingerprint (in hex)
type fingerprint string

func (f fingerprint) String() string {
	return string(f)
}

// Matches will check a given certificate against a hex encoded fingerprint
func (f fingerprint) Matches(c x509.Certificate) bool {
	fp := certutil.GetHexSHA256Fingerprint(c)

	// Check some constraints
	flen := utf8.RuneCountInString(f.String())
	if flen < minimumFingerprintLength {
		return false
	}

	// If the whitelist has a shortened fingerprint use it as a prefix
	// Otherwise, compare their full contents
	if flen < utf8.RuneCountInString(fp) {
		return strings.HasPrefix(fp, f.String())
	}
	return f.String() == fp
}