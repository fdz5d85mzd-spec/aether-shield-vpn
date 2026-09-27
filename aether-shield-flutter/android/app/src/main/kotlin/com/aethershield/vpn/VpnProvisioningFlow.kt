package com.aethershield.vpn

import android.content.Context

class VpnProvisioningFlow(
    context: Context,
    private val coordinator: ConnectionCoordinator,
    private val nativeKeys: NativeKeys = NativeKeys()
) {
    private val secureKeys = SecureKeyStore(context.applicationContext)

    suspend fun prepare(region: String = ""): ReadySession {
        val publicKey = secureKeys.publicKey(nativeKeys)
        require(publicKey.isNotBlank())
        return coordinator.provision(
            clientPublicKey = publicKey,
            region = region
        )
    }
}
