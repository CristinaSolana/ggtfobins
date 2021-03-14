# Get GTFOBINS
Get info from [GTFOBins](https://gtfobins.github.io/) about a given exploit for given commands

![Image of GGTFOBINS](ggtfobins-screenshot.jpg)

## Install
`go get github.com/CristinaSolana/ggtfobins`

## Usage
`ggtfobins  --exploit suid --bins cpan,bash`

## Command not found error
Run `export PATH=$PATH:$(go env GOPATH)/bin`

## Available Exploits
- bind-shell
- capabilities
- command
- file-download
- file-read
- file-upload
- file-write
- library-load
- limited-suid
- non-interactive-bind-shell
- non-interactive-reverse-shell
- reverse-shell
- shell
- sudo
- suid

---

[Contribute to GTFOBins](https://gtfobins.github.io/contribute/)

---

Follow GTFOBins' creators:
- [norbemi](https://twitter.com/norbemi)
- [cyrus_and](https://twitter.com/cyrus_and)

Follow [me](https://twitter.com/nightshiftc)
