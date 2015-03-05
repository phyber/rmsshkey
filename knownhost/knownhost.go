package knownhost

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"strings"
)

const (
	// HashMagic and HashDelim from hostfile.h in OpenSSH.
	HashDelim = "|"
	HashMagic = "|1|"

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

func (kh *KnownHost) extractB64Len() error {
	// Ensure that string contains the HashDelim.
	b64len := strings.Index(kh.entry, HashDelim)
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

func (kh *KnownHost) extractHost() error {
	host := kh.entry[kh.b64len+1:]
	kh.host = host
	return nil
}

func (kh *KnownHost) extractSalt() error {
	salt := kh.entry[0:kh.b64len]
	kh.salt = salt
	return nil
}

func (kh *KnownHost) setup() error {
	// Ensure that string is long enough.
	lenHashMagic := len(HashMagic)
	lenentry := len(kh.entry)
	if lenentry < lenHashMagic {
		return errors.New(errStringTooShort)
	}

	// Ensure that HashMagic is present at the beginning of the entry
	if !strings.HasPrefix(kh.entry, HashMagic) {
		return errors.New(errInvalidMagicIdentifier)
	}

	// Strip out the hashmagic.
	kh.entry = kh.entry[lenHashMagic:]

	kh.extractB64Len()
	kh.extractSalt()
	kh.extractHost()

	return nil
}

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

func (kh *KnownHost) Equals(host string) (bool, error) {
	hostHash, err := kh.b64Hmac(host)
	if err != nil {
		return false, err
	}
	return hostHash == kh.Host(), nil
}

func (kh *KnownHost) Host() string {
	return kh.host
}

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

func (kh *KnownHost) SaltB64() string {
	return kh.salt
}

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