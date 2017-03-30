box
===

- `box` reads choices from STDIN, 
- On `box`, the user selects one by cursor (HJKL,C-n&C-p&C-f&C-b)
- `box` outputs chosen one to STDOUT.

<img src="box0.png" />

How To Build
============

On Windows
----------

	git clone https://github.com/zetamatta/go-box
	cd go-box
	go get github.com/mattn/go-runewidth
	go get github.com/zetamatta/go-getch
	go build

On UNIX(tested on FreeBSD)
-------------------------

	git clone https://github.com/zetamatta/go-box
	cd go-box
	go get github.com/mattn/go-runewidth
	go get github.com/mattn/go-tty
	go build
