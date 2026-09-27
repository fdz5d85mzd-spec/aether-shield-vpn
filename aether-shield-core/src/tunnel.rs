#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum TunnelState {
    Stopped,
    Starting,
    Running,
    Stopping,
    Failed,
}

#[derive(Debug)]
pub struct TunnelLifecycle {
    state: TunnelState,
    tun_fd: Option<i32>,
}

impl Default for TunnelLifecycle {
    fn default() -> Self {
        Self {
            state: TunnelState::Stopped,
            tun_fd: None,
        }
    }
}

impl TunnelLifecycle {
    pub fn state(&self) -> TunnelState {
        self.state
    }

    pub fn start(&mut self, tun_fd: i32) -> Result<(), &'static str> {
        if tun_fd < 0 {
            self.state = TunnelState::Failed;
            return Err("invalid tun fd");
        }
        if self.state != TunnelState::Stopped {
            return Err("tunnel not stopped");
        }
        self.state = TunnelState::Starting;
        self.tun_fd = Some(tun_fd);
        // Packet forwarding is intentionally not started until a transport engine is integrated.
        self.state = TunnelState::Running;
        Ok(())
    }

    pub fn stop(&mut self) {
        self.state = TunnelState::Stopping;
        self.tun_fd = None;
        self.state = TunnelState::Stopped;
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn lifecycle_rejects_bad_fd() {
        let mut tunnel = TunnelLifecycle::default();
        assert!(tunnel.start(-1).is_err());
        assert_eq!(tunnel.state(), TunnelState::Failed);
    }

    #[test]
    fn lifecycle_starts_and_stops() {
        let mut tunnel = TunnelLifecycle::default();
        assert!(tunnel.start(7).is_ok());
        assert_eq!(tunnel.state(), TunnelState::Running);
        tunnel.stop();
        assert_eq!(tunnel.state(), TunnelState::Stopped);
    }
}
