# WireGuard node provisioning

The provisioning script creates a local WireGuard server key only when one does not already exist. The private key remains under `/etc/wireguard` with mode 0600 and is never written to the repository or sent in node heartbeats.

The script creates the interface configuration and enables kernel IP forwarding. It deliberately does **not** add client peers and does not configure a NAT/firewall policy because those depend on the deployment network interface and security policy.

Before a node can be considered usable, operators must separately configure firewall/NAT routing, start `wg-quick@wg0`, install the node-agent binary and runtime environment, and verify real client traffic.
