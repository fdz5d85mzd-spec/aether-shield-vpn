use base64::{engine::general_purpose::STANDARD, Engine as _};
use x25519_dalek::{PublicKey, StaticSecret};

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct KeyPair {
    private: [u8; 32],
    public: [u8; 32],
}

impl KeyPair {
    pub fn generate() -> Result<Self, getrandom::Error> {
        let mut private = [0u8; 32];
        getrandom::getrandom(&mut private)?;
        clamp_wireguard_secret(&mut private);
        let secret = StaticSecret::from(private);
        let public = PublicKey::from(&secret).to_bytes();
        Ok(Self { private, public })
    }

    pub fn from_private_bytes(mut private: [u8; 32]) -> Self {
        clamp_wireguard_secret(&mut private);
        let secret = StaticSecret::from(private);
        let public = PublicKey::from(&secret).to_bytes();
        Self { private, public }
    }

    pub fn private_bytes(&self) -> [u8; 32] {
        self.private
    }

    pub fn public_base64(&self) -> String {
        STANDARD.encode(self.public)
    }
}

fn clamp_wireguard_secret(key: &mut [u8; 32]) {
    key[0] &= 248;
    key[31] &= 127;
    key[31] |= 64;
}

impl Drop for KeyPair {
    fn drop(&mut self) {
        self.private.fill(0);
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn generated_keys_are_nonempty_and_public_is_base64() {
        let pair = KeyPair::generate().unwrap();
        assert_ne!(pair.private_bytes(), [0u8; 32]);
        assert_eq!(pair.public_base64().len(), 44);
    }

    #[test]
    fn private_material_is_clamped() {
        let pair = KeyPair::from_private_bytes([0xff; 32]);
        let private = pair.private_bytes();
        assert_eq!(private[0] & 7, 0);
        assert_eq!(private[31] & 128, 0);
        assert_ne!(private[31] & 64, 0);
    }
}
