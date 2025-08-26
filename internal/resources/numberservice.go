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

package resources

import (
	"context"
	werror "github.com/palantir/witchcraft-go-error"
	"scalmday-go-server-template/config"
	"scalmday-go-server-template/internal/generated/applicationapi/com/scalmday/mynumberservice"
	"scalmday-go-server-template/internal/permissions"

	"github.com/palantir/pkg/bearertoken"
	"github.com/palantir/pkg/refreshable/v2"
)

type numberService struct {
	runtime      refreshable.Refreshable[config.RuntimeConfig]
	userProvider permissions.UserProvider
}

func (n *numberService) Get(ctx context.Context, authHeader bearertoken.Token) (mynumberservice.MyNumber, error) {
	if _, err := n.userProvider.GetUser(ctx, authHeader); err != nil {
		return 0, werror.WrapWithContextParams(ctx, err, "user not found")
	}
	return mynumberservice.MyNumber(n.runtime.Current().MyFavoriteNumber), nil
}

func NewNumberService(ctx context.Context, userProvider permissions.UserProvider, runtime refreshable.Refreshable[config.RuntimeConfig]) mynumberservice.MyNumberService {
	return &numberService{runtime: runtime, userProvider: userProvider}
}
