# topojson - TopoJSON implementation in Go

[![Build Status](https://travis-ci.org/rubenv/topojson.svg?branch=master)](https://travis-ci.org/rubenv/topojson) [![GoDoc](https://godoc.org/github.com/rubenv/topojson?status.png)](https://godoc.org/github.com/rubenv/topojson)

Implements the TopoJSON specification:
https://github.com/mbostock/topojson-specification

Uses the GeoJSON implementation of paulmach:
https://github.com/paulmach/orb/geojson

Large parts are a port of the canonical JavaScript implementation, big chunks
of the test suite are ported as well:
https://github.com/mbostock/topojson

## Installation
```
go get github.com/rubenv/topojson
```

Import into your application with:

```go
import "github.com/rubenv/topojson"
```

## Usage

```go
topology := topojson.NewTopology(fc, nil)
```

Optionally pass options as the second argument.

This generates a Topology which can be encoded to JSON.

## License

    Copyright (c) 2016-2018, Ruben Vermeersch
    Copyright (c) 2012-2016, Michael Bostock
    All rights reserved.

    Redistribution and use in source and binary forms, with or without
    modification, are permitted provided that the following conditions are met:

    * Redistributions of source code must retain the above copyright notice, this
      list of conditions and the following disclaimer.

    * Redistributions in binary form must reproduce the above copyright notice,
      this list of conditions and the following disclaimer in the documentation
      and/or other materials provided with the distribution.

    * The name Michael Bostock may not be used to endorse or promote products
      derived from this software without specific prior written permission.

    THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
    AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
    IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
    DISCLAIMED. IN NO EVENT SHALL MICHAEL BOSTOCK BE LIABLE FOR ANY DIRECT,
    INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
    BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
    DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY
    OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
    NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
    EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
