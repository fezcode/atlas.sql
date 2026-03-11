# Atlas SQL

![Banner Image](./banner-image.png)

**atlas.sql** is a lightweight, keyboard-centric terminal user interface (TUI) for interacting with SQL databases. Part of the **Atlas Suite**, it provides a streamlined experience for running queries and visualizing results directly in your terminal.

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)

## ✨ Features

- 🗄️ **Multi-Engine Support:** Seamlessly connect to SQLite and PostgreSQL.
- 🌈 **Syntax Highlighting:** Write and view queries with real-time SQL syntax highlighting.
- 📊 **Interactive Tables:** Browse query results in a clean, scrollable table view.
- ⌨️ **Vim Bindings:** Navigate query history and results without leaving the keyboard.
- 💾 **Local First:** Lightweight and fast, perfect for quick data exploration.
- 📦 **Cross-Platform:** Binaries available for Windows, Linux, and macOS.

## 🚀 Installation

### From Source
```bash
git clone https://github.com/fezcode/atlas.sql
cd atlas.sql
go build -o atlas.sql .
```

## ⌨️ Usage

### SQLite
```bash
./atlas.sql sqlite://relative/path/to/db.sqlite
./atlas.sql sqlite:///absolute/path/to/db.sqlite
```

### PostgreSQL
```bash
./atlas.sql "postgres://user:password@localhost:5432/dbname?sslmode=disable"
```

## 🕹️ Controls

| Key | Action |
|-----|--------|
| `Enter` | Execute SQL query |
| `↑/↓` or `k/j` | Navigate query results (Table mode) |
| `Tab` | Switch between Query Input and Table Result |
| `Ctrl+C` | Quit |
| `Esc` | Clear Input / Exit focus |

## 🏗️ Building for all platforms

The project uses **gobake** to generate binaries for all platforms:

```bash
gobake build
```
Binaries will be placed in the `build/` directory.

## 📄 License
MIT License - see [LICENSE](LICENSE) for details.
