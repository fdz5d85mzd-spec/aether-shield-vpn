package com.aethershield.vpn

import android.content.Context
import android.security.keystore.KeyGenParameterSpec
import android.security.keystore.KeyProperties
import android.util.Base64
import java.security.KeyStore
import javax.crypto.Cipher
import javax.crypto.KeyGenerator
import javax.crypto.SecretKey
import javax.crypto.spec.GCMParameterSpec

class SecureKeyStore(private val context: Context) {
    private val alias = "aether-shield-wrap-v1"
    private val prefs = context.getSharedPreferences("aether-secure", Context.MODE_PRIVATE)

    fun loadOrCreatePrivateKey(nativeKeys: NativeKeys): ByteArray {
        loadPrivateKey()?.let { return it }
        val generated = nativeKeys.nativeGeneratePrivateKey()
            ?: error("native private key generation failed")
        require(generated.size == 32)
        storePrivateKey(generated)
        return generated
    }

    fun publicKey(nativeKeys: NativeKeys): String {
        val privateKey = loadOrCreatePrivateKey(nativeKeys)
        return try {
            nativeKeys.nativePublicKey(privateKey)
                ?: error("native public key derivation failed")
        } finally {
            privateKey.fill(0)
        }
    }

    private fun storePrivateKey(privateKey: ByteArray) {
        val cipher = Cipher.getInstance("AES/GCM/NoPadding")
        cipher.init(Cipher.ENCRYPT_MODE, wrappingKey())
        val encrypted = cipher.doFinal(privateKey)
        prefs.edit()
            .putString("private_ciphertext", Base64.encodeToString(encrypted, Base64.NO_WRAP))
            .putString("private_iv", Base64.encodeToString(cipher.iv, Base64.NO_WRAP))
            .apply()
    }

    private fun loadPrivateKey(): ByteArray? {
        val ciphertext = prefs.getString("private_ciphertext", null) ?: return null
        val iv = prefs.getString("private_iv", null) ?: return null
        val cipher = Cipher.getInstance("AES/GCM/NoPadding")
        cipher.init(
            Cipher.DECRYPT_MODE,
            wrappingKey(),
            GCMParameterSpec(128, Base64.decode(iv, Base64.NO_WRAP))
        )
        return cipher.doFinal(Base64.decode(ciphertext, Base64.NO_WRAP))
    }

    private fun wrappingKey(): SecretKey {
        val keyStore = KeyStore.getInstance("AndroidKeyStore").apply { load(null) }
        (keyStore.getKey(alias, null) as? SecretKey)?.let { return it }
        val generator = KeyGenerator.getInstance(KeyProperties.KEY_ALGORITHM_AES, "AndroidKeyStore")
        generator.init(
            KeyGenParameterSpec.Builder(
                alias,
                KeyProperties.PURPOSE_ENCRYPT or KeyProperties.PURPOSE_DECRYPT
            )
                .setBlockModes(KeyProperties.BLOCK_MODE_GCM)
                .setEncryptionPaddings(KeyProperties.ENCRYPTION_PADDING_NONE)
                .setKeySize(256)
                .build()
        )
        return generator.generateKey()
    }
}
