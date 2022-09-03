# [socketbuffer](https://github.com/danlapid/socketbuffer/)

Sync files between different hosts over a one way network link such as a Data Diode

[![Test Status](https://github.com/danlapid/socketbuffer/actions/workflows/test.yml/badge.svg)](https://github.com/danlapid/socketbuffer/actions?query=workflow%3ATest)
[![Coverage Status](https://coveralls.io/repos/github/danlapid/socketbuffer/badge.svg?branch=main)](https://coveralls.io/github/danlapid/socketbuffer?branch=main)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

## Exports

- GetReadBuffer: returns the size of the rx buffer
- GetAvailableBytes: returns the amount of bytes on the listening socket

Works on linux, windows and mac
