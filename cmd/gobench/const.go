// Copyright 2017-2018 zhvala
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"time"
)

// Program info
const (
	// AppVersion gobench version
	AppVersion = "1.00"
	// Copyright copyright info
	Copyright = "Copyright (c) zhvala 2017-2018, Apache 2.0 Open Source Software."
)

// HTTP Method supported
const (
	// GET HTTP GET
	GET = "GET"
	// POST HTTP POST
	POST = "POST"
	// HEAD HTTP HEAD
	HEAD = "HEAD"
	// OPTION HTTP OPTION
	OPTION = "OPTION"
	// TRACE HTTP TRACE
	TRACE = "TRACE"
)

const (
	// HTTP http1.1 is used as a default version in golang http.client
	HTTP = iota
	// HTTP2 http2
	HTTP2
)

const (
	// HTTPPrefix http url prefix
	HTTPPrefix = "http://"
	// HTTPSPrefix https url prefix
	HTTPSPrefix = "http://"
)

// Time value
const (
	cTimeMax        = time.Duration(1<<63 - 1)
	cTimeMin        = time.Duration(-1 << 63)
	cDefaultTimeout = 1000 // millisecond
)
