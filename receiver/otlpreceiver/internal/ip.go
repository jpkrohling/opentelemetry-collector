// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal // import "go.opentelemetry.io/collector/receiver/otlpreceiver/internal"

import "net"

// ParseIP parses the given string for an IP address. The input string might contain the port,
// but must not contain a protocol or path. Suitable for getting the IP part of a client connection.
func ParseIP(source string) string {
	ipstr, _, err := net.SplitHostPort(source)
	if err == nil {
		return ipstr
	}
	ip := net.ParseIP(source)
	if ip != nil {
		return ip.String()
	}
	return ""
}
