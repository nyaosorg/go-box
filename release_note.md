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
