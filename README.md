IRC
======================
[![Build Status](https://travis-ci.org/mengzhuo/irc.svg?branch=master)](https://travis-ci.org/mengzhuo/irc)
[![GoDoc](https://godoc.org/github.com/mengzhuo/irc?status.svg)](http://godoc.org/github.com/mengzhuo/irc)

Fast IRC decode/encode library with zero allocation

This project is inspired by 

* https://github.com/sorcix/irc
* https://github.com/valyala/fasthttp


### Benchmarks
```
benchmark                          old ns/op     new ns/op     delta
BenchmarkParseMessage_short-4      313           170           -45.69%
BenchmarkParseMessage_medium-4     485           225           -53.61%
BenchmarkParseMessage_long-4       627           442           -29.51%

benchmark                          old allocs     new allocs     delta
BenchmarkParseMessage_short-4      2              0              -100.00%
BenchmarkParseMessage_medium-4     3              0              -100.00%
BenchmarkParseMessage_long-4       3              0              -100.00%

benchmark                          old bytes     new bytes     delta
BenchmarkParseMessage_short-4      96            0             -100.00%
BenchmarkParseMessage_medium-4     160           0             -100.00%
BenchmarkParseMessage_long-4       240           0             -100.00%
```
