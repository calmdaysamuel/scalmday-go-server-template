// Copyright (c) 2025 Samuel Calmday. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/tls"
	"scalmday-go-server-template/config"
	"scalmday-go-server-template/internal/server"

	"github.com/palantir/witchcraft-go-server/v2/witchcraft"
)

func main() {
	if err := witchcraft.NewServer().
		WithInitFunc(server.Init).
		WithSelfSignedCertificate().
		WithECVKeyProvider(witchcraft.ECVKeyNoOp()).
		WithRuntimeConfigType(config.RuntimeConfig{}).
		WithInstallConfigType(config.InstallConfig{}).
		WithClientAuth(tls.NoClientCert).
		Start(); err != nil {
		panic(err)
	}
}
