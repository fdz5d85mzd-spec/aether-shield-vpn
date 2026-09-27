package com.aethershield.vpn

import android.content.Intent
import android.net.VpnService
import android.os.ParcelFileDescriptor

class AetherVpnService : VpnService() {
    private var tun: ParcelFileDescriptor? = null
    private var nativeStarted = false

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        if (tun != null) return START_STICKY

        val established = Builder()
            .setSession("Aether Shield")
            .addAddress("10.88.0.2", 32)
            .addRoute("0.0.0.0", 0)
            .setMtu(1420)
            .establish() ?: return START_NOT_STICKY

        tun = established
        val result = nativeTunnelStart(established.fd)
        nativeStarted = result == 0

        if (!nativeStarted) {
            established.close()
            tun = null
            stopSelf()
            return START_NOT_STICKY
        }

        // Native lifecycle is running, but packet forwarding remains unimplemented in Rust.
        return START_STICKY
    }

    override fun onDestroy() {
        if (nativeStarted) {
            nativeTunnelStop()
            nativeStarted = false
        }
        tun?.close()
        tun = null
        super.onDestroy()
    }

    private external fun nativeTunnelStart(tunFd: Int): Int
    private external fun nativeTunnelStop(): Int

    companion object {
        init {
            System.loadLibrary("aether_shield_core")
        }
    }
}
