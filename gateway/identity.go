package gateway

import (
	"github.com/a3tai/openclaw-go/identity"
	"github.com/a3tai/openclaw-go/protocol"
)

// WithIdentity configures the client to authenticate using a device keypair.
// If deviceToken is non-empty, it overrides the bearer token for this connect
// (the device uses its issued token instead of the shared gateway token).
// The identity is used to sign the server's challenge nonce.
func WithIdentity(id *identity.Identity, deviceToken string) Option {
	return func(o *options) {
		if deviceToken != "" {
			o.token = deviceToken
		}
		o.deviceSigner = func(nonce string) *protocol.DeviceIdentity {
			proto := id.BuildDeviceIdentity(nonce)
			return &protocol.DeviceIdentity{
				ID:        proto.ID,
				PublicKey: proto.PublicKey,
				Signature: proto.Signature,
				SignedAt:  proto.SignedAt,
				Nonce:     proto.Nonce,
			}
		}
	}
}
