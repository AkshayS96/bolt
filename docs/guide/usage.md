# Usage

## 🚀 Quick Start

Here are a few quick examples of `bolt` in action:

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

## 🔧 Pipe Support

All commands that accept input support `stdin` piping, which makes `bolt` extremely powerful when combined with other shell utilities:

```bash
cat data.json | bolt json format
echo "hello" | bolt base64 encode
cat file.txt | bolt hash sha256
```
