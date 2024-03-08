//  Copyright © 2024 Apple Inc. All rights reserved.
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

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/areese/piv-go/piv"
)

var (
	ErrFileNotFound           = errors.New("file not found")
	ErrNoCardsSelected        = errors.New("no yubikeys found")
	ErrNotYetImplemented      = errors.New("not yet implemented")
	ErrPathIsCurrentDirectory = errors.New("path is [.]")
	ErrPathIsRootDirectory    = errors.New("path is [/]")
)

// CloseYubikey closes a key and logs any errors.
// Deprecated: Use Close
// nolint:varnamelen
func CloseYubikey(logger LogI, yk *piv.YubiKey) {
	Close(context.Background(), logger, yk)
}

func Close(ctx context.Context, logger LogI, yk *piv.YubiKey) {
	logger = Nop(logger)
	if nil == yk {
		logger.DebugMsg("closing nil yubikey")

		return
	}

	if logger.IsDebugEnabled() {
		serial, _ := yk.Serial()

		logger.DebugMsgf("closing yubikey [%d]", serial)
	}

	closeErr := yk.Close()
	if closeErr != nil {
		logger.ErrorMsg(closeErr, "failed to close Yubikey")
	}
}

func ValidateFileFlag(ctx context.Context, logger LogI, commandName, fileName string) error {
	logger = Nop(logger)

	var err error

	if fileName == "" {
		err = fmt.Errorf("%w: empty filename passed", ErrFileNotFound)
		logger.ErrorMsg(err, "File flag was empty.")

		return err
	}

	err = IsValidFileArg(logger, commandName, []string{fileName})
	if err != nil {
		err = fmt.Errorf("%w: file [%s] was not found", err, fileName)
		logger.ErrorMsg(err, "Failed to find file.")

		return err
	}

	return nil
}

// MakeJSONString dumps a struct to json as a helper.
func MakeJSONString(data interface{}) string {
	prettifiedOSJSON, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Sprintf(`{ "error": "%s"}`, err.Error())
	}

	return string(prettifiedOSJSON)
}
