( **English** / [Japanese](./release_note_ja.md) )

- **Bumped the major version to v3.**
- **Removed deprecated functions and methods:**
  - `BoxT`, `Choce`, `ChoiceMulti`, `ChooseMulti`, `Choose`, `PrintNoLastLineFeed`, `Print`
- **Renamed and reorganized functions/methods:**
  - `PrintX` → `Print`: prints items in a tabular format **with** a trailing newline
  - `PrintNoLastLineFeedX` → `Print`: prints items in a tabular format **without** a trailing newline
  - `NewBox` → `New`: constructor for the `Box` type
  - `AnsiCutter` → unexported: regular expression for removing ANSI escape sequences
- **Reorganized source file layout:**
  - `box.go`: definitions related to the `Box` type
  - `main.go`: global definitions
  - `internal/ansi`: ANSI escape sequences
  - `internal/key`: key code definitions
- **Updated dependencies to the latest versions:**
  - `mattn/go-tty` v0.0.7
  - `mattn/go-colorable` v0.1.14
  - `mattn/go-runewidth` v0.0.19
- **Removed the unused `context.Context` parameter across the package.**
- **Extracted the grid-based display feature** from the `Box` type into a new `Grid` type under the `grid` subpackage.
- **Virtualized terminal I/O**, following the same design as **go-readline-ny**, allowing arbitrary terminal input backends.
  Set the `Box.Tty` field to one of the following before calling `Init()`.
  The `New()` function still defaults to `tty8` for backward compatibility.
  - `&tty8.Tty{}` from "[github.com/nyaosorg/go-ttyadapter][adapter]/tty8" → uses [github.com/mattn/go-tty][go-tty]
  - `&tty10.Tty{}` from "[github.com/nyaosorg/go-ttyadapter][adapter]/tty10" → uses [golang.org/x/term][xterm]
  - `&auto.Pilot{Text: []string{...}}` from "[github.com/nyaosorg/go-ttyadapter][adapter]/auto" → pseudo terminal for programmatic input
- **Removed the internal workaround for Unicode *Ambiguous-width* characters** on Windows Terminal.
  Since `github.com/mattn/go-runewidth` v0.0.17 now auto-detects Windows Terminal and adjusts widths appropriately.
  (See [mattn/go-runewidth#85](https://github.com/mattn/go-runewidth/pull/85))

[adapter]: https://github.com/nyaosorg/go-ttyadapter
[go-tty]: https://github.com/mattn/go-tty
[xterm]: https://pkg.go.dev/golang.org/x/term

v2.2.1
======
Apr 19, 2024

- Implement `[(*Box)]Select{Index,String}Context`
- Restore PrintNoLastLineFeed with `Deprecated:`
- Set `Deprecated:` to `BoxT`, `New`, `Print`, `(*Box) Print`, and `(*Box) PrintNoLineFeed`

v2.2.0
======
Apr 10, 2024

- Implement new functions and methods that returns error instead of calling panic on error
- Make TAB-Key same as RIGHT, and SHIFT-TAB as LEFT
- Fix: box.exe could not be built
- Rename BoxT to Box
- Add single selection mode
- Add test

v2.1.3
=======
Feb 20, 2022

- Fix: the import-path was old one in the test-code.  
  (as a result, `go get -u` downloaded both zetamatta/go-box and nyaosorg/go-box )

v2.1.2
=======
Sep 13, 2021

- Change owner: zetamatta to nyaosorg
- Fix: import "github.com/zetamatta/go-box/v2" remained

v2.0.8
=======
Feb 22, 2021

- Support Windows Terminal

v2.0.4
=======
Apr 13, 2019

Do not use `ESC[0A`
