# udp_forward
udp_forward can send udp data from one source to mutli destinations

## Usage

```sh
udp_forward can send udp data from one source to mutli destinations

Usage:
  udp_forward [flags]

Flags:
  -d, --destinations stringArray   destinations for udp data, e.g., udp:192.168.1.2:9000 or unix:/path/to/unix.sock
  -h, --help                       help for udp_forward
  -l, --listen string              listen for udp data, e.g., udp:0.0.0.0:514 or unix:/path/to/unix.sock
  -v, --verbose                    print info log
      --version                    show version
      --vv                         more verbose, print debug log
```

e.g.
- `udp_forward -l udp:0.0.0.0:9001 -d udp:192.168.1.2:9001`
- `udp_forward -l udp:0.0.0.0:9001 -d unix:/tmp/test.sock`
- `udp_forward -l unix:/tmp/test.sock -d udp:192.168.1.2:9001`
- `udp_forward -l udp:0.0.0.0:9001 -d udp:192.168.1.2:9001 -v`
    - `-v`, print info log
    - `-vv`, print debug log

## Install

### Build from source

```sh
git clone https://github.com/PengShaw/udp_forward.git
cd udp_forward
go build
```
