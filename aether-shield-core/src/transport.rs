#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum TransportState {
    NotConfigured,
    Ready,
    Active,
    Failed,
}

pub trait PacketTransport: Send {
    fn state(&self) -> TransportState;
    fn start(&mut self) -> Result<(), &'static str>;
    fn stop(&mut self);
}

#[derive(Debug, Default)]
pub struct UnconfiguredTransport;

impl PacketTransport for UnconfiguredTransport {
    fn state(&self) -> TransportState {
        TransportState::NotConfigured
    }

    fn start(&mut self) -> Result<(), &'static str> {
        Err("packet transport not configured")
    }

    fn stop(&mut self) {}
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn unconfigured_transport_never_claims_active() {
        let mut transport = UnconfiguredTransport;
        assert_eq!(transport.state(), TransportState::NotConfigured);
        assert!(transport.start().is_err());
        assert_eq!(transport.state(), TransportState::NotConfigured);
    }
}
