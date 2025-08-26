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

package server

import (
	"context"
	werror "github.com/palantir/witchcraft-go-error"
	"scalmday-go-server-template/config"
	"scalmday-go-server-template/internal/database"
	"scalmday-go-server-template/internal/generated/applicationapi/com/scalmday/mynumberservice"
	"scalmday-go-server-template/internal/permissions"
	"scalmday-go-server-template/internal/resources"

	"github.com/palantir/pkg/refreshable"
	"github.com/palantir/witchcraft-go-server/v2/witchcraft"
)

func Init(ctx context.Context, info witchcraft.InitInfo) (gracefulShutdownFunc func(), initializationErr error) {
	runtimeConfig := refreshable.ToV2[config.RuntimeConfig](info.RuntimeConfig)
	installConfig := info.InstallConfig.(config.InstallConfig)
	_, closer, err := database.Init(ctx, installConfig)
	if err != nil {
		return nil, err
	}

	userProvider, err := permissions.NewUserProvider(ctx, installConfig)
	if err != nil {
		return nil, werror.WrapWithContextParams(ctx, err, "failed to initialize user provider")
	}
	err = mynumberservice.RegisterRoutesMyNumberService(info.Router, resources.NewNumberService(ctx, userProvider, runtimeConfig))
	if err != nil {
		return nil, err
	}
	return func() {
		_ = closer(ctx)
	}, nil
}
