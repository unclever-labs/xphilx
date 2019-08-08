# xPHILx

<img src="https://raw.githubusercontent.com/unclever-labs/xphilx/master/xphilx.jpg" width="375" height="466"/>

## Introduction

A little agent to exfiltrate layer 7 payloads to S3 from a running server.

(Lol, sorry about the dramatic title. Just a cool Golang tool to scrape payloads if you can't update an application's logic itself)

## Contents

- [Install](#install)
- [How it Works](#how-it-works)
- [Usage](#usage)
- [Credits](#credits)


## Install

The easiest way to install is using `go get`:

```bash
go get -u -v github.com/unclever-labs/xphilx
```

You can install for OS X using homebrew:

```bash
brew install unclever-labs/unclever-labs/xphilx
```

You can also compile binaries yourself using:

```bash
./build.sh
# Look in directory: bin
```

## How it Works

<img src="https://raw.githubusercontent.com/unclever-labs/xphilx/master/xphilx.gif" width="1233" height="679"/>

- `./terminal-1.sh` is running `xphilx` listening on lo0 port 80
- `./terminal-2.sh` is running `python -m SimpleHTTPServer 80` serving port 80 filesystem
- `./terminal-3.sh` is running 3x `curl -XPOST -d '.....' localhost`

`xphilx` is an agent listening on a network interface and port. It captures the TCP stream bytes and runs them through Golang's `http.ReadRequest()` to generate a layer 7 request payload. The body of the payload then gets pushed to an S3 bucket.

## Usage

Open a new terminal and run:

```bash
./terminal-1.sh
```

Open a new terminal and run:

```bash
./terminal-2.sh
```

Open a new terminal and run:

```bash
./terminal-3.sh
```

## Credits

- Photo from @cctildef on IG
- A lot of copy pasta from: https://github.com/google/gopacket/blob/master/examples/httpassembly/main.go

**Please correct me if I'm wrong with my explanation of `How it Works` by opening an Issue or PR**
