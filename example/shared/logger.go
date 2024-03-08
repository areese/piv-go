// Copyright 2024 Apple, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shared

// LogI hides ZeroLogger from bits that don't need to care about it.
type LogI interface {
	VerboseMsg(message string)
	VerboseMsgf(format string, args ...interface{})
	// InfoMsgf will only log if quiet flag is NOT set.
	InfoMsg(message string)
	InfoMsgf(format string, args ...interface{})
	DebugMsgf(format string, args ...interface{})
	DebugMsg(message string)
	IsDebugEnabled() bool
	ErrorMsg(err error, message string)
	ErrorMsgf(err error, format string, args ...interface{})
}

type NopLogger struct{}

var _ LogI = (*NopLogger)(nil)

func (n *NopLogger) VerboseMsg(string)                                       {}
func (n *NopLogger) VerboseMsgf(string, ...interface{})                      {}
func (n *NopLogger) InfoMsg(string)                                          {}
func (n *NopLogger) InfoMsgf(string, ...interface{})                         {}
func (n *NopLogger) DebugMsgf(string, ...interface{})                        {}
func (n *NopLogger) DebugMsg(string)                                         {}
func (n *NopLogger) IsDebugEnabled() bool                                    { return false }
func (n *NopLogger) ErrorMsg(error, string)                                  {}
func (n *NopLogger) ErrorMsgf(err error, format string, args ...interface{}) {}

func Nop(l LogI) LogI {
	if l == nil {
		return &NopLogger{}
	}

	return l
}
