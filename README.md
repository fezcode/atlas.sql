# Atlas SQL

![Banner Image](./banner-image.png)

**atlas.sql** is a lightweight, keyboard-centric terminal user interface (TUI) for interacting with SQL databases. Part of the **Atlas Suite**, it provides a streamlined experience for running queries, visualizing results, and managing data directly in your terminal.

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)
![Version](https://img.shields.io/badge/version-0.2.0-blue)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)

## ✨ Features

- 🗄️ **Multi-Engine Support:** Seamlessly connect to SQLite and PostgreSQL.
- 🌈 **Syntax Highlighting:** Write and view queries with real-time feedback.
- 📊 **Interactive Tables:** Navigate query results with arrow-key based horizontal and vertical scrolling.
- 🔍 **Detail View:** Press `v` to see the full content of a selected row in a scrollable vertical view—perfect for long text fields.
- 📏 **Dynamic Resizing:** Adjust column widths on the fly using `+` and `-` keys.
- 📋 **CSV Export:** Copy entire result sets to your clipboard instantly with `c`.
- 🛡️ **Safety First:** Automatically appends `LIMIT 500` to SELECT statements if no limit is specified.
- 📝 **Full DDL/DML:** Support for `CREATE`, `INSERT`, `UPDATE`, and `DELETE` with "rows affected" reporting.
- 🎓 **Built-in Tutorial:** Access a quick SQL cheat sheet directly inside the help menu.
- 📦 **Cross-Platform:** Binaries available for Windows, Linux, and macOS.

## 🚀 Installation

### From Source
```bash
git clone https://github.com/fezcode/atlas.sql
cd atlas.sql
go build -o atlas.sql .
```

## ⌨️ Usage

Run the tool by providing a connection string:

### SQLite
```bash
./atlas.sql sqlite://path/to/db.sqlite
```

### PostgreSQL
```bash
./atlas.sql "postgres://user:password@localhost:5432/dbname?sslmode=disable"
```

## 🕹️ Controls

| Key | Action |
|-----|--------|
| `Enter` | Execute SQL query |
| `Tab` | Switch focus between Input and Table |
| `↑ / ↓` | Navigate result rows |
| `← / →` | Scroll table columns horizontally |
| `+ / -` | Increase / Decrease column width |
| `v` | Toggle Detail View for selected row |
| `c` | Copy entire result set as CSV to clipboard |
| `h` | Toggle the interactive Help & SQL Tutorial |
| `Ctrl+T` | Quick-list all tables in the database |
| `Ctrl+S` | Quick-list all schemas/attached databases |
| `Esc` | Clear Input / Exit focus / Close menus |
| `q / Ctrl+C` | Quit |

## 🏗️ Building for all platforms

The project uses **gobake** to generate binaries for all platforms:

```bash
gobake build
```
Binaries will be placed in the `build/` directory.

## 📄 License
MIT License - see [LICENSE](LICENSE) for details.
