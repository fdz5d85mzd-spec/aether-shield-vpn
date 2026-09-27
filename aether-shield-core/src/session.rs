use crate::transport::{PacketTransport, TransportState};

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum SessionState {
    Stopped,
    TunReady,
    TransportReady,
    Connected,
    Failed,
}

pub struct TunnelSession<T: PacketTransport> {
    state: SessionState,
    tun_fd: Option<i32>,
    transport: T,
}

impl<T: PacketTransport> TunnelSession<T> {
    pub fn new(transport: T) -> Self {
        Self {
            state: SessionState::Stopped,
            tun_fd: None,
            transport,
        }
    }

    pub fn state(&self) -> SessionState {
        self.state
    }

    pub fn attach_tun(&mut self, fd: i32) -> Result<(), &'static str> {
        if fd < 0 {
            self.state = SessionState::Failed;
            return Err("invalid tun fd");
        }
        self.tun_fd = Some(fd);
        self.state = SessionState::TunReady;
        Ok(())
    }

    pub fn connect(&mut self) -> Result<(), &'static str> {
        if self.tun_fd.is_none() {
            self.state = SessionState::Failed;
            return Err("tun not attached");
        }
        if self.transport.state() == TransportState::NotConfigured {
            return Err("transport not configured");
        }
        self.state = SessionState::TransportReady;
        self.transport.start()?;
        if self.transport.state() != TransportState::Active {
            self.state = SessionState::Failed;
            return Err("transport failed to become active");
        }
        self.state = SessionState::Connected;
        Ok(())
    }

    pub fn stop(&mut self) {
        self.transport.stop();
        self.tun_fd = None;
        self.state = SessionState::Stopped;
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::transport::UnconfiguredTransport;

    #[test]
    fn tun_ready_is_not_connected() {
        let mut session = TunnelSession::new(UnconfiguredTransport);
        session.attach_tun(9).unwrap();
        assert_eq!(session.state(), SessionState::TunReady);
        assert!(session.connect().is_err());
        assert_eq!(session.state(), SessionState::TunReady);
    }
}
