# Android client status

The Android scaffold now has a control-plane provisioning flow:

1. The client supplies its **public** WireGuard key to the control plane.
2. It creates a tunnel session.
3. It polls the opaque session ID while the node applies the peer.
4. It accepts endpoint/server public key/client address only after state becomes `READY`.
5. The native VPN service must still receive a real client-side private key/config and a functioning Rust WireGuard packet engine before encrypted packet forwarding can start.

The control plane never needs the client's private key. Key generation/storage should use an appropriate client-side secure storage mechanism. The current repository does not yet implement that key lifecycle or the Rust packet engine.
