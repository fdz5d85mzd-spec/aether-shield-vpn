# Android client status

The Android scaffold now has a control-plane provisioning flow:

1. The client creates or loads its WireGuard private key locally.
2. The private key is wrapped at rest with an AES-GCM key held by Android Keystore.
3. The Rust native core derives the client **public** key.
4. Only the public key is supplied to the control plane.
5. The client creates a tunnel session and polls the opaque session ID.
6. The node applies the peer and acknowledges the command.
7. Endpoint, server public key and client address are accepted only after state becomes `READY`.
8. The Rust core has a validated local configuration handoff for these values.

## Still not implemented

A production WireGuard packet engine is not yet connected to the Rust `PacketTransport` implementation. Therefore a verified READY control-plane session is **not** equivalent to encrypted packet forwarding.

The Android service integration that would feed the READY session and locally decrypted private key into the native configuration boundary remains incomplete. Private key material must never be sent to the control plane or written to logs.
