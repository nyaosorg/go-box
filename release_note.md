( **English** / [Japanese](./release_note_ja.md) )

* **Bumped the major version to v3.**
* **Removed deprecated functions and methods:**
  * `BoxT`
  * `Choce`
  * `ChoiceMulti`
  * `ChooseMulti`
  * `Choose`
  * `PrintNoLastLineFeed`
  * `Print`
* **Renamed and reorganized functions/methods:**
  * `PrintX` → `Print`: prints items in a tabular format **with** a trailing newline
  * `PrintNoLastLineFeedX` → `Print`: prints items in a tabular format **without** a trailing newline
  * `NewBox` → `New`: constructor for the `Box` type
  * `AnsiCutter` → made unexported: regular expression for removing ANSI escape sequences
* **Reorganized function placements across source files:**
  * `box.go`: definitions related to the `Box` type
  * `main.go`: global definitions
  * `tty.go`: terminal-related code
  * `internal/ansi`: ANSI escape sequence definitions
  * `internal/key`: key code definitions
* **Updated imported packages to the latest versions:**
  * `mattn/go-tty` → v0.0.7
  * `mattn/go-colorable` → v0.1.14
  * `mattn/go-runewidth` → v0.0.19

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
