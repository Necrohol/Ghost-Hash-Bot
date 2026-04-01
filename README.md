# 👻 Ghost-Hash-Bot 💀

**Spooky fast, diskless, and a tad bit blind.** Ghost-Hash-Bot is a DevOps-focused Gentoo Manifest generator designed to chew through dozens of ebuilds in parallel. It streams sources directly from the wire into RAM-based hashers (BLAKE2B/SHA512) and then into the void (`io.Discard`). No disk bloat, no `/tmp` pollution, and no Portage "hoopla."

## 🚀 The "Lunch Break" Workflow
Fixed a broken manifest from a Windows machine or a phone? Use the cloud:
1. **Phantom-Bot** sniffs the eclass logic from **Codeberg**.
2. **Gen-Helpers** writes a Go-native fetcher.
3. **Ghost-Hash** streams the bits and spits out a valid `Manifest` line.

---

## 🏗️ Project Architecture (K.I.S.S.)

The bot is split into independent Go files to allow for hot-plugging new eclass logic without touching the core.

| File | Purpose |
| :--- | :--- |
| `main.go` | The Core CLI (urfave/cli). Handles flags and user flow. |
| `ghost.go` | The IO Engine. Streams URL -> MultiWriter -> /dev/null. |
| `int.go` | The Integration Hub. Maps eclasses to their generated helpers. |
| `regex-rename.go` | The Sniffer. Extracts `SRC_URI` and `->` logic from eclasses. |
| `gen-helpers.go` | The Generator. Turns Sniffed logic into `helpers_[eclass].go`. |

---

## 🛠️ RAD Assembly Line

### 1. Sniff an Eclass
If you need to support a new eclass (e.g., `zig-utils`), point the sniffer at your local repo or the Codeberg mirror:
```bash
go run regex-rename.go /var/db/repos/gentoo/eclass/zig-utils.eclass
# Ghost-Hash-Bot
Ghost-Hash Bot
