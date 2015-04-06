# ADRS - A DNS(Recursive) Server
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](LICENSE.md)
[![Build Status](https://travis-ci.org/chzyer/adrs.svg?branch=master)](https://travis-ci.org/chzyer/adrs)
[![Coverage Status](https://coveralls.io/repos/chzyer/adrs/badge.svg?branch=master)](https://coveralls.io/r/chzyer/adrs?branch=master)
[![GoDoc](https://godoc.org/github.com/chzyer/adrs?status.svg)](https://godoc.org/github.com/chzyer/adrs)

A implementation of recursive DNS Server in [the Go programming language](https://golang.org).   
`ADRS` will support protocols like HTTP, TCP, UDP, and use [redis](http://redis.io) to store records cached.

**NOTE: THIS PROJECT IS STILL UNDER DEVELOPMENT!**

### Feature
* Support protocol like HTTP, TCP, UDP.
* Support routers which is based on (sub)domain to direct to foreign name servers.
* Support wrong records detection.
* Support custom(internal) domain resolves (may needs a web dashboard).

### Topology
```
                 Local Host                        |  Foreign
                                                   |
    +---------+               +----------+         |  +--------+
    |         | user queries  |          |queries  |  |        |
    |  User   |-------------->|          |---------|->|Foreign |
    | Program |               | Resolver |         |  |  Name  |
    |         |<--------------|          |<--------|--| Server |
    |         | user responses|          |responses|  |        |
    +---------+               +----------+         |  +--------+
                                |     A            |
                cache additions |     | references |
                                V     |            |
                              +----------+         |
                              |  Shared  |         |
                              | database |         |
                              +----------+         |
                                A     |            |
      +---------+     refreshes |     | references |
     /         /|               |     V            |
    +---------+ |             +----------+         |  +--------+
    |         | |             |          |responses|  |        |
    |         | |             |   Name   |---------|->|Foreign |
    |  Master |-------------->|  Server  |         |  |Resolver|
    |  files  | |             |          |<--------|--|        |
    |         |/              |          | queries |  +--------+
    +---------+               +----------+         |
                                A     |maintenance |  +--------+
                                |     +------------|->|        |
                                |      queries     |  |Foreign |
                                |                  |  |  Name  |
                                +------------------|--| Server |
                             maintenance responses |  +--------+
```
