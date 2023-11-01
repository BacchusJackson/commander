# Installation

The Commander CLI can be installed in a number of ways thanks to the fantastic tooling provided by Go.

## Pre-compiled binaries

Executable binaries can be downloaded from the GitHub Releases page. 
Download the binary for your platform (Windows, macOS, or Linux) and extract the archive. 

## Build from source using Golang

If you have the latest version of (Go)[https://go.dev/] installed on your system, you can build Commander from source.
The easiest way is to use `make build` which will put the binary in the folder `bin` at the project root.

```shell
git clone https://github.com/BacchusJackson/commander.git
make build
./bin/commander --help
```
