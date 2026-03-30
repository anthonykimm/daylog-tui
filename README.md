# daylog

A terminal UI for tracking daily tasks, git commits, and Linear issues.

![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)

## Features

- **Task management** — add, complete, delete, and hide tasks
- **Git commit tracking** — auto-detects today's commits by author
- **Linear integration** — OAuth sync, link issues to tasks, mark done bidirectionally
- **Daily summary** — auto-generated from tasks and commits, editable, copy to clipboard
- **Calendar navigation** — browse past days in read-only mode
- **SQLite storage** — XDG-compliant paths (`~/.local/share/daylog/`)

## Install

Requires [Go 1.26+](https://go.dev/dl/).

```
go install github.com/anthonykimm/daylog@latest
```

Then run:

```
daylog
```

## Keybindings

| Key | Action |
|---|---|
| `a` | Add task |
| `space` / `enter` | Toggle complete / link issue |
| `d` | Delete or hide |
| `r` | Refresh commits and issues |
| `u` | Toggle hidden items |
| `c` | Copy summary to clipboard |
| `i` | Edit summary |
| `g` | Go to date (calendar) |
| `o` | Open Linear issue in browser |
| `L` | Linear setup / menu |
| `1`/`2`/`3`/`4` | Switch pane |
| `j`/`k` | Navigate |
| `q` | Quit |

## Linear Setup

1. Press `L` to start setup
2. Enter your Linear OAuth Client ID and Secret
3. Authorize in browser — daylog listens on `localhost:19284` for the callback
4. Issues appear in the Linear pane; press `enter` to link one as a task

## License

MIT
