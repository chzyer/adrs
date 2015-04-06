# ADRS (A Dns Recursive Server)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](LICENSE.md)
[![Build Status](https://travis-ci.org/chzyer/adrs.svg?branch=master)](https://travis-ci.org/chzyer/adrs)
[![Coverage Status](https://coveralls.io/repos/chzyer/adrs/badge.svg?branch=master)](https://coveralls.io/r/chzyer/adrs?branch=master)
[![GoDoc](https://godoc.org/github.com/chzyer/adrs?status.svg)](https://godoc.org/github.com/chzyer/adrs)

A DNS recursion server implements by golang.  
Which will supports protocols like HTTP, TCP, UDP, then use redis to store cached records.

**NOTE: THIS PROJECT IS STILL UNDER DEVELOPMENT!**

### Feature
* supports protocol like HTTP, TCP, UDP.
* supports routers which based on (sub)domain to direct to foreign name servers.
* supports wrong records detection.
* supports custom(internal) domain resolves (may needs a web dashboard).

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

To be continue.
