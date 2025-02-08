#!/bin/sh

input_lines(){
	echo '2025-02-07T22:09:14.012345Z INFO [user-service] 127.0.0.1 GET HTTP/1.1 200 "curl/8.7" 127.0.0.1   /health
	echo '2025-02-08T22:09:14.012345Z INFO [user-service] 127.0.0.1 GET HTTP/1.1 200 "curl/8.7" 192.168.2.9 /api
}

input_lines |
  ./extractip4 |
  python3 -c 'import ipaddress; import operator; import struct; import sys; import functools; functools.reduce(
  	lambda state, f: f(state),
	[
		functools.partial(map, struct.Struct(">I").unpack),
		functools.partial(map, operator.itemgetter(0)),
		functools.partial(map, ipaddress.IPv4Address),
		functools.partial(map, print),
		lambda prints: sum(1 for _ in prints),
	],
	iter(
		functools.partial(sys.stdin.buffer.read, 4),
		b"",
	),
  )'
