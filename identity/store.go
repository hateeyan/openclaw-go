// Package identity manages device keypairs for OpenClaw gateway authentication.
package identity

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	keypairFile     = "keypair.json"
	deviceTokenFile = "device-token"
)

// Store persists and manages device identity (Ed25519 keypair + device token).
type Store struct {
	dir string
}

// NewStore creates a new identity store backed by files in dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("identity store: mkdir %s: %w", dir, err)
	}
	return &Store{dir: dir}, nil
}

type keypairJSON struct {
	DeviceID   string `json:"deviceId"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

// Identity is the loaded device identity.
type Identity struct {
	DeviceID        string
	PublicKeyB64URL string
	privateKey      ed25519.PrivateKey
}

// DeviceIdentityProto mirrors protocol.DeviceIdentity without import cycles.
type DeviceIdentityProto struct {
	ID        string `json:"id"`
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
	SignedAt  int64  `json:"signedAt"`
	Nonce     string `json:"nonce,omitempty"`
}

// BuildDeviceIdentity creates a signed DeviceIdentity for the given nonce.
func (id *Identity) BuildDeviceIdentity(nonce string) *DeviceIdentityProto {
	now := time.Now().UnixMilli()
	msg := fmt.Sprintf("%s:%s:%d", nonce, id.DeviceID, now)
	h := sha256.Sum256([]byte(msg))
	sig := ed25519.Sign(id.privateKey, h[:])
	return &DeviceIdentityProto{
		ID:        id.DeviceID,
		PublicKey: id.PublicKeyB64URL,
		Signature: base64.RawURLEncoding.EncodeToString(sig),
		SignedAt:  now,
		Nonce:     nonce,
	}
}

// LoadOrGenerate loads the device keypair from disk, generating one if absent.
func (s *Store) LoadOrGenerate() (*Identity, error) {
	id, err := s.load()
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, errNotFound) {
		return nil, err
	}
	return s.generate()
}

var errNotFound = errors.New("keypair not found")

func (s *Store) load() (*Identity, error) {
	data, err := os.ReadFile(filepath.Join(s.dir, keypairFile))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errNotFound
		}
		return nil, fmt.Errorf("read keypair: %w", err)
	}
	var kp keypairJSON
	if err := json.Unmarshal(data, &kp); err != nil {
		return nil, fmt.Errorf("parse keypair: %w", err)
	}
	seed, err := base64.RawURLEncoding.DecodeString(kp.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("decode private key: %w", err)
	}
	return &Identity{
		DeviceID:        kp.DeviceID,
		PublicKeyB64URL: kp.PublicKey,
		privateKey:      ed25519.NewKeyFromSeed(seed),
	}, nil
}

func (s *Store) generate() (*Identity, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate keypair: %w", err)
	}
	pubB64 := base64.RawURLEncoding.EncodeToString(pub)
	h := sha256.Sum256(pub)
	deviceID := fmt.Sprintf("%x", h[:16])
	kp := keypairJSON{
		DeviceID:   deviceID,
		PublicKey:  pubB64,
		PrivateKey: base64.RawURLEncoding.EncodeToString(priv.Seed()),
	}
	data, err := json.MarshalIndent(kp, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal keypair: %w", err)
	}
	if err := os.WriteFile(filepath.Join(s.dir, keypairFile), data, 0600); err != nil {
		return nil, fmt.Errorf("write keypair: %w", err)
	}
	return &Identity{
		DeviceID:        deviceID,
		PublicKeyB64URL: pubB64,
		privateKey:      priv,
	}, nil
}

// LoadDeviceToken returns the stored device token, or "" if none.
func (s *Store) LoadDeviceToken() string {
	data, err := os.ReadFile(filepath.Join(s.dir, deviceTokenFile))
	if err != nil {
		return ""
	}
	return string(data)
}

// SaveDeviceToken persists the device token.
func (s *Store) SaveDeviceToken(token string) error {
	return os.WriteFile(filepath.Join(s.dir, deviceTokenFile), []byte(token), 0600)
}

// ClearDeviceToken removes the stored device token.
func (s *Store) ClearDeviceToken() error {
	err := os.Remove(filepath.Join(s.dir, deviceTokenFile))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Reset deletes all identity data (keypair + device token).
func (s *Store) Reset() error {
	for _, name := range []string{keypairFile, deviceTokenFile} {
		if err := os.Remove(filepath.Join(s.dir, name)); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("reset identity: remove %s: %w", name, err)
		}
	}
	return nil
}
