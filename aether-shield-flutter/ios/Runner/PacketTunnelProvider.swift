import NetworkExtension

@_silgen_name("aether_tunnel_start")
private func aetherTunnelStart(_ tunFd: Int32) -> Int32

@_silgen_name("aether_tunnel_stop")
private func aetherTunnelStop() -> Int32

final class PacketTunnelProvider: NEPacketTunnelProvider {
    private var nativeStarted = false

    override func startTunnel(
        options: [String: NSObject]?,
        completionHandler: @escaping (Error?) -> Void
    ) {
        let settings = NEPacketTunnelNetworkSettings(tunnelRemoteAddress: "not-configured")
        let ipv4 = NEIPv4Settings(
            addresses: ["10.88.0.2"],
            subnetMasks: ["255.255.255.255"]
        )
        ipv4.includedRoutes = [NEIPv4Route.default()]
        settings.ipv4Settings = ipv4

        setTunnelNetworkSettings(settings) { error in
            if let error {
                completionHandler(error)
                return
            }

            // NetworkExtension does not expose packetFlow as a Unix TUN fd.
            // Do not call aetherTunnelStart until a real packetFlow/FFI adapter exists.
            self.nativeStarted = false
            completionHandler(nil)
        }
    }

    override func stopTunnel(
        with reason: NEProviderStopReason,
        completionHandler: @escaping () -> Void
    ) {
        if nativeStarted {
            _ = aetherTunnelStop()
            nativeStarted = false
        }
        completionHandler()
    }
}
