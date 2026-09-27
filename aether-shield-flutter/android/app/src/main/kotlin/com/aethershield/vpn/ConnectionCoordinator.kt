package com.aethershield.vpn

import kotlinx.coroutines.delay

enum class ClientConnectionState {
    IDLE,
    REQUESTING_SESSION,
    WAITING_FOR_NODE,
    READY_TO_CONFIGURE,
    FAILED
}

class ConnectionCoordinator(
    private val control: ControlPlaneClient
) {
    var state: ClientConnectionState = ClientConnectionState.IDLE
        private set

    suspend fun provision(
        clientPublicKey: String,
        region: String = "",
        maxAttempts: Int = 20,
        delayMs: Long = 1_500
    ): ReadySession {
        require(clientPublicKey.isNotBlank())
        state = ClientConnectionState.REQUESTING_SESSION
        return try {
            val sessionId = control.createSession(clientPublicKey, region)
            state = ClientConnectionState.WAITING_FOR_NODE
            repeat(maxAttempts) {
                val ready = control.getReadySession(sessionId)
                if (ready != null) {
                    state = ClientConnectionState.READY_TO_CONFIGURE
                    return ready
                }
                delay(delayMs)
            }
            error("session provisioning timed out")
        } catch (error: Throwable) {
            state = ClientConnectionState.FAILED
            throw error
        }
    }
}
