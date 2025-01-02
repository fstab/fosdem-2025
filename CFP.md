Zero-Code Distributed Traces for any programming language
---------------------------------------------------------

We live in a world of many programming languages, even more application development frameworks and many different versions of these technologies. Wouldn’t it be nice if we could automatically get OpenTelemetry distributed tracing to work for all of these applications, in a manner that’s performant, requires zero-code changes, low overhead and doesn’t need maintaining instrumentation support for every single application development framework or version?

We’d like to present two novel approaches of using eBPF to achieve OpenTelemetry trace context propagation across services residing on different nodes, without the need for any external services for trace matching or ordering.

The first approach uses TCP/IP Level 4 embedding of trace context information which works for any protocol, even when encryption is enabled. We’ll show you how metadata can be transmitted and embedded in various parts of the TCP/IP packet to send context to the other side, outside of the application protocol.

The second approach uses TCP/IP Level 7 protocol manipulation to embed the additional trace context information. We’ll show two different Level 7 protocol extension approaches we’ve implemented and describe some of their pros and cons.
