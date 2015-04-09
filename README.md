# ADRS - A DNS(Recursive) server

A implementation of recursive DNS Server in [the Go programming language](https://golang.org).  

**NOTE: THIS PROJECT IS STILL UNDER DEVELOPMENT!**

### Status

[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](LICENSE.md)
[![Build Status](https://travis-ci.org/chzyer/adrs.svg?branch=master)](https://travis-ci.org/chzyer/adrs)
[![Coverage Status](https://coveralls.io/repos/chzyer/adrs/badge.svg?branch=master)](https://coveralls.io/r/chzyer/adrs?branch=master)
[![GoDoc](https://godoc.org/github.com/chzyer/adrs?status.svg)](https://godoc.org/github.com/chzyer/adrs)  

[![Join the chat at https://gitter.im/chzyer/adrs](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/chzyer/adrs?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)


### Features (road map)
* Support protocol like HTTP, TCP, UDP.
* Support routers which is based on (sub)domain to direct to foreign name servers.
* Support wrong records detection.
* Support custom(internal) domain resolves (may needs a web dashboard).
* Support distributed deployment.
* Support load balances like round-robin.

### Progress now

* Support protocol UDP
