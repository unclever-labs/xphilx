# xPHILx

<img src="https://raw.githubusercontent.com/unclever-labs/xphilx/master/xphilx.jpg" width="375" height="466"/>

## Introduction

A little agent to exfiltrate layer 7 payloads to S3 from a running server

## Contents

- [How it Works](#how-it-works)
- [Usage](#usage)
- [Credits](#credits)


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
