#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct CoreCapabilities {
    pub wireguard_engine_compiled: bool,
    pub packet_forwarding_available: bool,
    pub post_quantum_available: bool,
}

pub const fn capabilities() -> CoreCapabilities {
    CoreCapabilities {
        wireguard_engine_compiled: cfg!(feature = "wireguard-engine"),
        // A feature flag alone is not proof that a runtime engine is usable.
        packet_forwarding_available: false,
        post_quantum_available: false,
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn default_build_does_not_claim_packet_forwarding() {
        let caps = capabilities();
        assert!(!caps.packet_forwarding_available);
        assert!(!caps.post_quantum_available);
    }
}
