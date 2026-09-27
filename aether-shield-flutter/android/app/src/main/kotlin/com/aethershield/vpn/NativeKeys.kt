package com.aethershield.vpn

class NativeKeys {
    external fun nativeGeneratePrivateKey(): ByteArray?
    external fun nativePublicKey(privateKey: ByteArray): String?

    companion object {
        init {
            System.loadLibrary("aether_shield_core")
        }
    }
}
