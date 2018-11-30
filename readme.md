box
===

- `box` reads choices from STDIN, 
- On `box`, the user selects one by cursor (HJKL,C-n&C-p&C-f&C-b)
- `box` outputs chosen one to STDOUT.

<img src="box0.png" />

How To Build
============

        git clone https://github.com/zetamatta/go-box
        cd $GOPATH/src/github.com/zetamatta/go-box/v2
        go get github.com/mattn/go-runewidth
        go get github.com/mattn/go-tty
        go build
        cd box
        go build
