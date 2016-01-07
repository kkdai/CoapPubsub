PubSub over CoAP(Constrained Application Protocol 
==================

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/kkdai/CoapPubsub/master/LICENSE)  [![GoDoc](https://godoc.org/github.com/kkdai/CoapPubsub?status.svg)](https://godoc.org/github.com/kkdai/CoapPubsub)  [![Build Status](https://travis-ci.org/kkdai/CoapPubsub.svg?branch=master)](https://travis-ci.org/kkdai/CoapPubsub)


Comparison between CoAP, XMPP, Restful HTTP, MQTT
---------------

| Protocol  |      CoAP      |  XMPP |RESTful HTTP | MQTT |
|----------|:-------------:|------:|------:|------:|
| Transport |  UDP | TCP| TCP | TCP |
| Messaging |    Request/Response    |   Publish/Subscribe Request/Response |   Request/Response |   Publish/Subscribe  |
| LLN Suitability (1000s nodes) | Excellent |    Excellent |    Excellent |    Excellent |
| Success Storied | Utility Field Area Networks |    Remote management of consumer white goods |    Smart Energy Profile 2 (premise energy management/home services) |    Extending enterprise messaging into IoT applications |
    
    


Features
---------------






Install
---------------
`go get github.com/kkdai/CoapPubsub`


Usage
---------------

```go


```



Benchmark
---------------

```
```

Inspired
---------------

- [MQTT and CoAP, IoT Protocols](https://eclipse.org/community/eclipse_newsletter/2014/february/article2.php)
- [RFC 7252](http://tools.ietf.org/html/rfc7252)
- [https://github.com/dustin/go-coap](https://github.com/dustin/go-coap)
- [CoAP an introduction](http://www.herjulf.se/download/coap-2013-fall.pdf)

Project52
---------------

It is one of my [project 52](https://github.com/kkdai/project52).


License
---------------

This package is licensed under MIT license. See LICENSE for details.

