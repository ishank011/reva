// Copyright 2018-2020 CERN
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
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package gateway

import (
	"context"
	"strings"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/status"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/pkg/errors"
)

func (s *svc) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	c, err := pool.GetUserProviderServiceClient(s.c.UserProviderEndpoint)
	if err != nil {
		return &user.GetUserResponse{
			Status: status.NewInternal(ctx, err, "error getting auth client"),
		}, nil
	}

	res, err := c.GetUser(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling GetUser")
	}

	return res, nil
}

func (s *svc) FindUsers(ctx context.Context, req *user.FindUsersRequest) (*user.FindUsersResponse, error) {
	c, err := pool.GetUserProviderServiceClient(s.c.UserProviderEndpoint)
	if err != nil {
		return &user.FindUsersResponse{
			Status: status.NewInternal(ctx, err, "error getting auth client"),
		}, nil
	}

	res, err := c.FindUsers(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling GetUser")
	}

	return res, nil
}

func (s *svc) GetUserGroups(ctx context.Context, req *user.GetUserGroupsRequest) (*user.GetUserGroupsResponse, error) {
	c, err := pool.GetUserProviderServiceClient(s.c.UserProviderEndpoint)
	if err != nil {
		return &user.GetUserGroupsResponse{
			Status: status.NewInternal(ctx, err, "error getting auth client"),
		}, nil
	}

	res, err := c.GetUserGroups(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling GetUser")
	}

	return res, nil
}

func (s *svc) IsInGroup(ctx context.Context, req *user.IsInGroupRequest) (*user.IsInGroupResponse, error) {
	c, err := pool.GetUserProviderServiceClient(s.c.UserProviderEndpoint)
	if err != nil {
		return &user.IsInGroupResponse{
			Status: status.NewInternal(ctx, err, "error getting auth client"),
		}, nil
	}

	res, err := c.IsInGroup(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling GetUser")
	}

	return res, nil
}

func (s *svc) resolveUIDToUser(ctx context.Context, uid *user.UserId) (*user.UserId, error) {
	if !strings.HasPrefix(uid.OpaqueId, "uid:") {
		return uid, nil
	}
	id := strings.TrimPrefix(uid.OpaqueId, "uid:")
	getUserReq := &user.GetUserRequest{
		Opaque: &types.Opaque{
			Map: map[string]*types.OpaqueEntry{
				"uid": &types.OpaqueEntry{
					Decoder: "plain",
					Value:   []byte(id),
				},
			},
		},
	}
	getUserRes, err := s.GetUser(ctx, getUserReq)
	if err != nil {
		return nil, err
	}
	if getUserRes.Status.Code != rpc.Code_CODE_OK {
		return nil, errors.New("error resolving UID to user ID")
	}

	return getUserRes.User.Id, nil
}
