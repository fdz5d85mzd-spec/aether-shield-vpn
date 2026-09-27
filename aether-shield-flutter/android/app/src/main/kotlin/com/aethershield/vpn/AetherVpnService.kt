package com.aethershield.vpn

import android.content.Intent
import android.net.VpnService
import android.os.ParcelFileDescriptor

class AetherVpnService : VpnService() {
    private var tun: ParcelFileDescriptor? = null

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        if (tun != null) return START_STICKY
        tun = Builder()
            .setSession("Aether Shield")
            .addAddress("10.88.0.2", 32)
            .addRoute("0.0.0.0", 0)
            .setMtu(1420)
            .establish()
        // Deliberately no "connected" claim: the Rust packet-processing loop is not wired yet.
        return START_STICKY
    }

    override fun onDestroy() {
        tun?.close()
        tun = null
        super.onDestroy()
    }
}
