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

package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"

	"github.com/areese/piv-go/example/shared"
	"github.com/areese/piv-go/piv"
)

func main() {
}

//nolint:deadcode
func Encrypt(fileName string, logger shared.LogI) error {
	ctx := context.Background()

	pgpConfig := shared.Config{
		CardSelection: nil,
		Debug:         false,
		Verbose:       false,
		Trace:         false,
		Quiet:         false,
		ShowPublic:    false,
		Base64Encoded: false,
	}

	// do the general setup
	yubikeyClient, fileBytes, err := pgpConfig.DevEncryptDecryptSetup(ctx, logger, "encrypt", fileName)
	defer func() { yubikeyClient.Close(ctx, logger) }()

	if err != nil {
		return err
	}

	// we have file bytes, now we can encrypt it.
	var pubkey *rsa.PublicKey

	pubkey, err = yubikeyClient.ReadPublicKey(ctx, logger, piv.AsymmetricConfidentiality)
	if err != nil {
		logger.ErrorMsg(err, "Failed to load public key")

		return err
	}

	fileSize := len(fileBytes)
	// nolint:gomnd // -11 must be a header size, this is from the rsa.go source.
	maxSize := pubkey.Size() - 11

	if fileSize > maxSize {
		logger.InfoMsgf("Only encrypting [%d] bytes out of [%d] from file [%s]", maxSize, fileSize, fileName)
		fileSize = maxSize
	}

	var encData []byte

	encData, err = rsa.EncryptPKCS1v15(rand.Reader, pubkey, fileBytes[:fileSize])
	if err != nil {
		logger.ErrorMsg(err, "Failed to load public key")

		return err
	}

	logger.InfoMsg(base64.StdEncoding.EncodeToString(encData))

	return nil
}
