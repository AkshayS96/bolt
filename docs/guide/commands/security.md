# 🔒 Security
| Command | Description |
|---------|-------------|
| `bolt jwt decode <token>` | Decode JWT (header + payload) |
| `bolt jwt header <token>` | Show JWT header |
| `bolt jwt payload <token>` | Show JWT payload |
| `bolt jwt exp <token>` | Show expiration & validity |
| `bolt jwt verify <token> --secret` | Verify HMAC signature |
| `bolt hash md5/sha1/sha256 <text>` | Hash text |
| `bolt hash file <path>` | Hash a file (SHA256) |
| `bolt password generate` | Generate password |
| `bolt password strong` | Strong password with symbols |
| `bolt entropy <text>` | Calculate password entropy |
