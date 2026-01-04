# SubMonsif - Subdomain Engine................

![banner](https://img.shields.io/badge/Subdomain--Finder-Made%20By%20MrMonsif-orange)

## Overview

**SubMonsif** is a powerful and advanced subdomain enumeration tool written in Go. It supports active and passive subdomain discovery from multiple real sources and can use bruteforce with customizable wordlists.
Ideal for bug bounty, penetration testing, and recon workflows.

---

## Features

* Fast, multi-threaded subdomain enumeration
* Uses real public APIs: **crt.sh, HackerTarget, VirusTotal, OTX, Shodan**
* Supports bruteforce discovery
* Removes duplicate results automatically
* Custom wordlist support
* Works on **Linux** and **Windows**

---

## Installation

### Requirements

* [Go](https://go.dev/doc/install) 1.18+
* Git (for cloning repo)

### Build (Linux/MacOS)

```bash
git clone https://github.com/monsifhmouri/SubMonsif.git
cd SubMonsif
go build -o SubMonsif main.go
chmod +x SubMonsif
```

### Build (Windows)

```cmd
git clone https://github.com/monsifhmouri/SubMonsif.git
cd SubMonsif
go build -o SubMonsif.exe main.go
```

---

## Usage

### Linux

```bash
./SubMonsif -d example.com -t 100 -brute -v
./SubMonsif -dl domains.txt -t 200 -o results.txt
```

### Windows

```cmd
SubMonsif.exe -d example.com -t 100 -brute -v
SubMonsif.exe -dl domains.txt -t 200 -o results.txt
```

### Options

| Flag     | Description                     |
| -------- | ------------------------------- |
| `-d`     | Target domain (example.com)     |
| `-dl`    | File with list of domains       |
| `-t`     | Number of threads (default 100) |
| `-brute` | Enable bruteforce mode          |
| `-v`     | Verbose mode                    |
| `-o`     | Output file for results         |

---

## Supported Passive Sources

* [crt.sh](https://crt.sh/)
* [HackerTarget](https://hackertarget.com/)
* [VirusTotal](https://www.virustotal.com/)
* [AlienVault OTX](https://otx.alienvault.com/)
* [Shodan](https://www.shodan.io/)

---

## How to Get Your API Keys

Some sources require free registration to use their APIs.
You **must** put your API keys directly in the correct functions in `providers/passive.go`:

| Source     | Register & Get Key                                                               | Where to put the key in code         |
| ---------- | -------------------------------------------------------------------------------- | ------------------------------------ |
| VirusTotal | [https://www.virustotal.com/gui/join-us](https://www.virustotal.com/gui/join-us) | Replace value in `getFromVirusTotal` |
| OTX        | [https://otx.alienvault.com/api/](https://otx.alienvault.com/api/)               | Replace value in `getFromOTX`        |
| Shodan     | [https://account.shodan.io/register](https://account.shodan.io/register)         | Replace value in `getFromShodan`     |

**Edit `providers/passive.go`:**
Replace the placeholder `apiKey := "............."` with your actual key in each relevant function:

```go
// Example:
apiKey := "YOUR_VIRUSTOTAL_API_KEY"
// and the same for OTX and Shodan
```

---

## Example (Output)

```
[Brute] Found: www.target.com
[+] Found: api.target.com
...
[+] Found 5523 subdomains
```

---

## Credits

* Author: [MrMonsif](https://github.com/monsifhmouri)
* Inspired by subfinder and other recon tools

---

## Warning

**Never commit or push your API keys to any public repository!**
If you plan to share or publish your code on GitHub, remove or replace all API keys before pushing.

---

## License

MIT License


