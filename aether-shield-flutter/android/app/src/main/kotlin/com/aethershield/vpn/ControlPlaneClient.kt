package com.aethershield.vpn

import org.json.JSONObject
import java.net.HttpURLConnection
import java.net.URL

data class ReadySession(
    val id: String,
    val endpoint: String,
    val serverPublicKey: String,
    val clientAddress: String
)

class ControlPlaneClient(private val baseUrl: String) {
    fun createSession(clientPublicKey: String, region: String = ""): String {
        val body = JSONObject()
            .put("clientPublicKey", clientPublicKey)
            .put("region", region)
        val response = request("POST", "/v1/tunnel/session", body.toString())
        if (response.code != 202) error("session create failed: ${response.code}")
        return JSONObject(response.body).getString("sessionId")
    }

    fun getReadySession(sessionId: String): ReadySession? {
        val response = request("GET", "/v1/tunnel/session/$sessionId", null)
        if (response.code != 200) error("session status failed: ${response.code}")
        val json = JSONObject(response.body)
        return when (json.getString("state")) {
            "READY" -> ReadySession(
                id = json.getString("id"),
                endpoint = json.getString("endpoint"),
                serverPublicKey = json.getString("serverPublicKey"),
                clientAddress = json.getString("clientAddress")
            )
            "FAILED" -> error("session provisioning failed: ${json.optString("error")}")
            else -> null
        }
    }

    private fun request(method: String, path: String, body: String?): Response {
        val connection = URL(baseUrl.trimEnd('/') + path).openConnection() as HttpURLConnection
        connection.requestMethod = method
        connection.connectTimeout = 10_000
        connection.readTimeout = 10_000
        connection.setRequestProperty("Content-Type", "application/json")
        if (body != null) {
            connection.doOutput = true
            connection.outputStream.use { it.write(body.toByteArray(Charsets.UTF_8)) }
        }
        val code = connection.responseCode
        val stream = if (code in 200..299) connection.inputStream else connection.errorStream
        val text = stream?.bufferedReader()?.use { it.readText() } ?: ""
        connection.disconnect()
        return Response(code, text)
    }

    private data class Response(val code: Int, val body: String)
}
