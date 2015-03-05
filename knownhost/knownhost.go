// Package knownhost provides an interface for working with individual
// known_hosts entries.

package knownhost

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"strings"
)

const (
	// hashMagic and hashDelim from hostfile.h in OpenSSH.
	hashDelim = "|"
	hashMagic = "|1|"

	// Error texts are mostly copied from OpenSSH hostfile.c.
	errBadSaltLength          = "extract_salt: bad encoded salt length"
	errInvalidMagicIdentifier = "extract_salt: invalid magic identifier"
	errInvalidSaltLen         = "extract_salt: invalid salt length after decode"
	errMissingHashDelim       = "extract_salt: missing hash delimiter"
	errSaltDecodeError        = "extract_salt: salt decode error"
	errStringTooShort         = "extract_salt: string too short"
)

type KnownHost struct {
	entry  string
	salt   string
	host   string
	b64len int
}

// extractB64Len extracts the base64 length from a KnownHost entry.
func (kh *KnownHost) extractB64Len() error {
	// Ensure that string contains the hashDelim.
	b64len := strings.Index(kh.entry, hashDelim)
	if b64len == -1 {
		return errors.New(errMissingHashDelim)
	}

	// Ensure that the b64len is a sensible size.
	if b64len == 0 || b64len > 1024 {
		return errors.New(errBadSaltLength)
	}

	// Store the b64len
	kh.b64len = b64len

	return nil
}

// extractHost extracts the host portion of a KnownHosts entry.
func (kh *KnownHost) extractHost() error {
	host := kh.entry[kh.b64len+1:]
	kh.host = host
	return nil
}

// extractSalt extracts the salt portion of a KnownHosts entry.
func (kh *KnownHost) extractSalt() error {
	salt := kh.entry[0:kh.b64len]
	kh.salt = salt
	return nil
}

// setup prepares the various entries in the KnownHost struct.
func (kh *KnownHost) setup() error {
	// Ensure that string is long enough.
	lenHashMagic := len(hashMagic)
	lenentry := len(kh.entry)
	if lenentry < lenHashMagic {
		return errors.New(errStringTooShort)
	}

	// Ensure that hashMagic is present at the beginning of the entry
	if !strings.HasPrefix(kh.entry, hashMagic) {
		return errors.New(errInvalidMagicIdentifier)
	}

	// Strip out the hashmagic.
	kh.entry = kh.entry[lenHashMagic:]

	kh.extractB64Len()
	kh.extractSalt()
	kh.extractHost()

	return nil
}

// b64Hmac computes the base64 HMAC for a given host using the salt of the
// current KnownHost entry.
func (kh *KnownHost) b64Hmac(host string) (string, error) {
	salt, err := kh.Salt()
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha1.New, salt)
	mac.Write([]byte(host))
	newMac := mac.Sum(nil)

	return base64.StdEncoding.EncodeToString(newMac), nil
}

// Equals checks if the provided host matches the current KnownHost.
func (kh *KnownHost) Equals(host string) (bool, error) {
	hostHash, err := kh.b64Hmac(host)
	if err != nil {
		return false, err
	}
	return hostHash == kh.Host(), nil
}

// Host returns the Host portion of the KnownHost.
func (kh *KnownHost) Host() string {
	return kh.host
}

// Salt returns the decoded Salt portion of the KnownHost.
func (kh *KnownHost) Salt() ([]byte, error) {
	// Snip the b64 encoded salt out of the known_host line.
	b64salt := kh.SaltB64()

	// Decode it.
	data, err := base64.StdEncoding.DecodeString(b64salt)
	if err != nil {
		return []byte(""), err //errors.New(errSaltDecodeError)
	}

	// Check size.
	if len(data) != sha1.Size {
		return []byte(""), errors.New(errInvalidSaltLen)
	}

	// Yay
	return data, nil
}

// SaltB64 returns the base64 encoded salt of the KnownHost
func (kh *KnownHost) SaltB64() string {
	return kh.salt
}

// New takes a raw known_hosts line and prepares a *KnownHost.
func New(knownHost string) (*KnownHost, error) {
	components := strings.Split(knownHost, " ")
	if len(components) < 2 {
		return nil, errors.New("Invalid known host entry.")
	}

	kh := KnownHost{
		entry: components[0],
	}

	err := kh.setup()
	if err != nil {
		return nil, err
	}

	return &kh, nil
}
