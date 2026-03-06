# Commands

Bolt comes packed with over 55+ utility commands.

## 🔑 ID Generators
| Command | Description |
|---------|-------------|
| `bolt uuid` | Generate a UUID v4 |
| `bolt uuid short` | Short unique ID (8 chars) |
| `bolt nanoid` | Generate a NanoID |
| `bolt cuid` | Generate a CUID-like unique ID |
| `bolt random [length]` | Random alphanumeric string |

## 📦 Data & Encoding
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

## 🔒 Security
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

## ⏰ Time & Date
| Command | Description |
|---------|-------------|
| `bolt time now` | Current time (ISO, Unix, UTC, Local) |
| `bolt time unix` | Current Unix timestamp |
| `bolt time iso` | Current ISO 8601 timestamp |
| `bolt time convert <ts>` | Convert Unix ↔ ISO |
| `bolt time diff <d1> <d2>` | Difference between dates |

## 🌐 HTTP & Network
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

## ✏️ Text & Strings
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

## 📁 Files & System
| Command | Description |
|---------|-------------|
| `bolt file hash <path>` | SHA256 hash of file |
| `bolt file size <path>` | File size |
| `bolt file lines <path>` | Count lines |
| `bolt file stats <path>` | File metadata |
| `bolt diff <file1> <file2>` | Colored file diff |

## 🖼️ Image Processing
| Command | Description |
|---------|-------------|
| `bolt img resize --width <w> <in> <out>` | Resize an image |
| `bolt img crop --ratio <ratio> <in> <out>` | Crop an image |
| `bolt img info <in>` | Image metadata (size, format) |
| `bolt img placeholder <WxH> <out>` | Generate basic placeholder |
| `bolt img blur --sigma <s> <in> <out>` | Apply Gaussian blur |

## 🧰 Utilities
| Command | Description |
|---------|-------------|
| `bolt color hex2rgb <hex>` | Hex → RGB |
| `bolt color rgb2hex <r> <g> <b>` | RGB → Hex |
| `bolt color random` | Random color (Hex + RGB + HSL) |
| `bolt lorem [n]` | Lorem ipsum paragraphs |
| `bolt lorem -w [n]` | Lorem ipsum words |
| `bolt cron explain "<expr>"` | Explain cron expression |
