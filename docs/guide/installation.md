# Installation

### Go Install (Recommended)

```bash
go install github.com/AkshayS96/bolt@latest
```

### Homebrew (macOS / Linux)

```bash
brew tap AkshayS96/bolt
brew install bolt
```

### APT (Debian / Ubuntu)

```bash
curl -sL https://github.com/AkshayS96/bolt/releases/latest/download/bolt_linux_amd64.deb -o bolt.deb
sudo dpkg -i bolt.deb
```

### Windows (Scoop)

```powershell
scoop bucket add AkshayS96 https://github.com/AkshayS96/scoop-bolt
scoop install bolt
```

### Manual Download

Grab the latest binary from [**Releases**](https://github.com/AkshayS96/bolt/releases) for your platform.

---

## 🏗️ Build from Source

```bash
git clone https://github.com/AkshayS96/bolt.git
cd bolt
go build -o bolt .
./bolt --version
```
