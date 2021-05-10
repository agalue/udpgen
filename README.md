# udpgen - UDP Packet Generator

## Overview

An application inspired by the C++ version of the [udpgen](https://github.com/OpenNMS/udpgen) tool written in [Go](https://golang.org/) to help the execution of stress tests against the OpenNMS, generating large volumes of traffic.

Unlike its C++ sibling, this one can be compiled and executed on any operating system except for Windows due to Syslog compatibility; however, it might not be as fast as the original version.

It currently supports packet generation via UDP for:
* SNMP Traps
* Syslog Messages
* Netflow 5 flows
* Netflow 9 flows (coming soon)

## Building

### Local

Make sure to have Go version [1.16.x](https://golang.org/dl/) installed on your system.

```bash=
go install
```

### Docker

```bash=
docker build -t agalue/udp .
```

## Usage

### Generate Syslog Message

Generate 100000 Syslog messages per second targeted at 172.23.1.1:514.

```bash=
udpgen -r 100000 -h 172.23.1.1 -p 514
```

Or via Docker:

```bash=
docker run -it --rm agalue/udpgen -r 100000 -h 172.23.1.1 -p 514
```

### Generate SNMP Traps

Generate 20000 SNMPv2 traps per second targeted at 127.0.0.1:1162.

```bash=
udpgen -x snmp -r 20000 -h 127.0.0.1 -p 1162
```

### Generate Netflow 5 flows

Generate 200000 Netflow 5 flows per second targeted at 192.168.0.1:8877.

```sh
udpgen -x netflow5 -r 200000 -h 192.168.0.1 -p 8877
```

### TODO

* Make the Syslog messages and SNMP Trap content dynamic.
* Implement Netflow 9.
