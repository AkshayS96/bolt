# ⚡ Bolt – Developer Swiss Army Knife CLI

A fast, single-binary toolbox for everyday developer tasks — built in Go.

UUIDs, JSON formatting, JWT decoding, hashing, encoding, HTTP requests, timestamps, passwords, and **50+ more commands** in one CLI.

---

## 📦 Installation

### Go Install (Recommended)

```bash
go install github.com/akshaysolanki/bolt@latest
```

### Homebrew (macOS / Linux)

```bash
brew tap akshaysolanki/bolt
brew install bolt
```

### APT (Debian / Ubuntu)

```bash
curl -sL https://github.com/akshaysolanki/bolt/releases/latest/download/bolt_linux_amd64.deb -o bolt.deb
sudo dpkg -i bolt.deb
```

### Windows (Scoop)

```powershell
scoop bucket add akshaysolanki https://github.com/akshaysolanki/scoop-bolt
scoop install bolt
```

### Manual Download

Grab the latest binary from [**Releases**](https://github.com/akshaysolanki/bolt/releases) for your platform.

---

## 🚀 Quick Start

```bash
bolt uuid                          # Generate UUID v4
bolt json format data.json         # Pretty print JSON
bolt hash sha256 hello             # SHA256 hash
bolt base64 encode "hello world"   # Base64 encode
bolt jwt decode <token>            # Decode JWT
bolt time now                      # Current timestamp
bolt password strong               # Generate strong password
bolt http get https://api.com      # HTTP GET request
echo '{"a":1}' | bolt json format  # Pipe support ✓
```

---

## 📖 Commands

### 🔑 ID Generators
| Command | Description |
|---------|-------------|
| `bolt uuid` | Generate a UUID v4 |
| `bolt uuid short` | Short unique ID (8 chars) |
| `bolt nanoid` | Generate a NanoID |
| `bolt cuid` | Generate a CUID-like unique ID |
| `bolt random [length]` | Random alphanumeric string |

### 📦 Data & Encoding
| Command | Description |
|---------|-------------|
| `bolt json format [file]` | Pretty print JSON |
| `bolt json minify [file]` | Minify JSON |
| `bolt json query <path> [file]` | Query JSON (GJSON syntax) |
| `bolt json validate [file]` | Validate JSON syntax |
| `bolt json to-yaml [file]` | Convert JSON → YAML |
| `bolt json from-yaml [file]` | Convert YAML → JSON |
| `bolt base64 encode/decode` | Base64 encode/decode |
| `bolt url encode/decode` | URL encode/decode |
| `bolt hex encode/decode` | Hex encode/decode |

### 🔒 Security
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

### ⏰ Time & Date
| Command | Description |
|---------|-------------|
| `bolt time now` | Current time (ISO, Unix, UTC, Local) |
| `bolt time unix` | Current Unix timestamp |
| `bolt time iso` | Current ISO 8601 timestamp |
| `bolt time convert <ts>` | Convert Unix ↔ ISO |
| `bolt time diff <d1> <d2>` | Difference between dates |

### 🌐 HTTP & Network
| Command | Description |
|---------|-------------|
| `bolt http get <url>` | GET request |
| `bolt http post <url>` | POST request (body from stdin) |
| `bolt http headers <url>` | Show response headers |
| `bolt http json <url>` | GET with pretty JSON |
| `bolt dns lookup <domain>` | DNS lookup (A, MX, NS, TXT) |
| `bolt ping <host>` | TCP ping |
| `bolt ip` | Show local IP addresses |
| `bolt port check <port>` | Check if port is in use |
| `bolt port kill <port>` | Kill process on port |

### ✏️ Text & Strings
| Command | Description |
|---------|-------------|
| `bolt slug <text>` | Slugify text |
| `bolt case camel/snake/kebab/pascal` | Case conversion |
| `bolt trim <text>` | Trim whitespace |
| `bolt length <text>` | Count chars, bytes, words |
| `bolt regex match <pattern> <text>` | Find regex matches |
| `bolt regex test <pattern> <text>` | Test regex match |
| `bolt regex extract <pattern> <text>` | Extract capture groups |
| `bolt regex replace <p> <r> <text>` | Regex replace |

### 📁 Files & System
| Command | Description |
|---------|-------------|
| `bolt file hash <path>` | SHA256 hash of file |
| `bolt file size <path>` | File size |
| `bolt file lines <path>` | Count lines |
| `bolt file stats <path>` | File metadata |
| `bolt diff <file1> <file2>` | Colored file diff |

### 🧰 Utilities
| Command | Description |
|---------|-------------|
| `bolt color hex2rgb <hex>` | Hex → RGB |
| `bolt color rgb2hex <r> <g> <b>` | RGB → Hex |
| `bolt color random` | Random color (Hex + RGB + HSL) |
| `bolt lorem [n]` | Lorem ipsum paragraphs |
| `bolt lorem -w [n]` | Lorem ipsum words |
| `bolt cron explain "<expr>"` | Explain cron expression |

---

## 🔧 Pipe Support

All commands that accept input support stdin piping:

```bash
cat data.json | bolt json format
echo "hello" | bolt base64 encode
cat file.txt | bolt hash sha256
```

---

## 🏗️ Build from Source

```bash
git clone https://github.com/akshaysolanki/bolt.git
cd bolt
go build -o bolt .
./bolt --version
```

### Cross-compile

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o bolt-linux .

# Windows
GOOS=windows GOARCH=amd64 go build -o bolt.exe .

# macOS ARM (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o bolt-darwin-arm64 .
```

---

## 🤝 Contributing

1. Fork the repo
2. Create a new command file in `cmd/`
3. Register with `rootCmd.AddCommand()` in `init()`
4. Submit a PR

---

## 📄 License

MIT License © [Akshay Solanki](https://x.com/__akshaysolanki)

---

<p align="center">
  Made with ❤️ by <a href="https://x.com/__akshaysolanki">@__akshaysolanki</a>
</p>
