package com.aethershield.vpn

import android.content.Intent
import android.net.VpnService
import android.os.ParcelFileDescriptor

class AetherVpnService : VpnService() {
    private var tun: ParcelFileDescriptor? = null
    private var nativeTunAttached = false
    private var transportConnected = false

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        if (tun != null) return START_STICKY

        val established = Builder()
            .setSession("Aether Shield")
            .addAddress("10.88.0.2", 32)
            .addRoute("0.0.0.0", 0)
            .setMtu(1420)
            .establish() ?: return START_NOT_STICKY

        tun = established
        nativeTunAttached = nativeTunnelStart(established.fd) == 0
        if (!nativeTunAttached) {
            cleanup()
            return START_NOT_STICKY
        }

        transportConnected = nativeTransportConnect() == 0
        if (!transportConnected) {
            // Correct behavior for the current scaffold: a TUN interface alone is not a VPN.
            cleanup()
            return START_NOT_STICKY
        }

        return START_STICKY
    }

    override fun onDestroy() {
        cleanup()
        super.onDestroy()
    }

    private fun cleanup() {
        if (nativeTunAttached) nativeTunnelStop()
        nativeTunAttached = false
        transportConnected = false
        tun?.close()
        tun = null
        stopSelf()
    }

    private external fun nativeTunnelStart(tunFd: Int): Int
    private external fun nativeTransportConnect(): Int
    private external fun nativeTunnelStop(): Int

    companion object {
        init {
            System.loadLibrary("aether_shield_core")
        }
    }
}
