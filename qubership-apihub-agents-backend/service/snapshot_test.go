// Copyright 2024-2025 NetCracker Technology Corporation
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

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Netcracker/qubership-apihub-agents-backend/view"
)

// mockApihubClient is a mock implementation of ApihubClient for testing
type mockApihubClient struct {
	getPackageByIdFunc func(ctx context.Context, id string) (*view.SimplePackage, error)
	createPackageFunc  func(ctx context.Context, pkg view.PackageCreateRequest) (string, error)
}

func (m *mockApihubClient) GetPackageById(ctx context.Context, id string) (*view.SimplePackage, error) {
	if m.getPackageByIdFunc != nil {
		return m.getPackageByIdFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockApihubClient) CreatePackage(ctx context.Context, pkg view.PackageCreateRequest) (string, error) {
	if m.createPackageFunc != nil {
		return m.createPackageFunc(ctx, pkg)
	}
	return "", nil
}

// Implement other required methods with empty implementations
func (m *mockApihubClient) CheckAuthToken(ctx context.Context, token string) (bool, error) {
	return false, nil
}

func (m *mockApihubClient) GetApiKeyByKey(ctx context.Context, apiKey string) (*view.ApihubApiKeyView, error) {
	return nil, nil
}

func (m *mockApihubClient) GetPatByPAT(ctx context.Context, token string) (*view.PersonalAccessTokenExtAuthView, error) {
	return nil, nil
}

func (m *mockApihubClient) GetPackageByServiceName(ctx context.Context, workspaceId string, serviceName string) (*view.PackagesInfo, error) {
	return nil, nil
}

func (m *mockApihubClient) GetPackages(ctx context.Context, searchReq view.PackagesSearchReq) (*view.Packages, error) {
	return nil, nil
}

func (m *mockApihubClient) GetUserPackagesPromoteStatuses(ctx context.Context, packagesReq view.PackagesReq) (view.AvailablePackagePromoteStatuses, error) {
	return nil, nil
}

func (m *mockApihubClient) GetVersion(ctx context.Context, id, version string) (*view.VersionContent, error) {
	return nil, nil
}

func (m *mockApihubClient) Publish(ctx context.Context, config view.BuildConfig, src []byte, clientBuild bool, builderId string, saveSources bool, dependencies []string) (string, error) {
	return "", nil
}

func (m *mockApihubClient) GetVersions(ctx context.Context, packageId string, searchReq view.VersionSearchRequest) (*view.PublishedVersionsView, error) {
	return nil, nil
}

func (m *mockApihubClient) DeleteVersionsRecursively(ctx context.Context, packageId string, req view.DeleteVersionsRecursivelyReq) (string, error) {
	return "", nil
}

func (m *mockApihubClient) GetVersionReferences(ctx context.Context, id, version string) (*view.VersionReferences, error) {
	return nil, nil
}

func (m *mockApihubClient) GetVersionRestOperationsWithData(ctx context.Context, packageId string, version string, limit int, page int) (*view.RestOperations, error) {
	return nil, nil
}

func (m *mockApihubClient) GetPublishStatuses(ctx context.Context, packageId string, publishIds []string) ([]view.PublishStatusResponse, error) {
	return nil, nil
}

func (m *mockApihubClient) GetApiKeyById(ctx context.Context, apiKeyId string) (*view.ApihubApiKeyView, error) {
	return nil, nil
}

func (m *mockApihubClient) GetUserById(ctx context.Context, userId string) (*view.User, error) {
	return nil, nil
}

func (m *mockApihubClient) GetSystemInfo(ctx context.Context) (*view.ApihubSystemInfo, error) {
	return nil, nil
}

func TestCreateGroupHierarchy(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		targetGroupId  string
		mockClient     *mockApihubClient
		expectedError  string
		expectedCalls  map[string]int  // expected number of calls to GetPackageById for each groupId
		expectedCreate map[string]bool // groups that should be created
	}{
		{
			name:          "success - all groups already exist",
			targetGroupId: "workspace.snapshots.cloud.namespace",
			mockClient: &mockApihubClient{
				getPackageByIdFunc: func(ctx context.Context, id string) (*view.SimplePackage, error) {
					return &view.SimplePackage{
						Id:   id,
						Kind: string(view.KindGroup),
					}, nil
				},
			},
			expectedCalls: map[string]int{
				"workspace":                           1,
				"workspace.snapshots":                 1,
				"workspace.snapshots.cloud":           1,
				"workspace.snapshots.cloud.namespace": 1,
			},
			expectedCreate: map[string]bool{},
		},
		{
			name:          "success - workspace exists, all other groups need to be created",
			targetGroupId: "workspace.snapshots.cloud.namespace",
			mockClient: &mockApihubClient{
				getPackageByIdFunc: func(ctx context.Context, id string) (*view.SimplePackage, error) {
					if id == "workspace" {
						return &view.SimplePackage{
							Id:   id,
							Kind: string(view.KindWorkspace),
						}, nil
					}
					// All other groups don't exist
					return nil, nil
				},
				createPackageFunc: func(ctx context.Context, pkg view.PackageCreateRequest) (string, error) {
					return pkg.ParentId + "." + pkg.Alias, nil
				},
			},
			expectedCalls: map[string]int{
				"workspace":                           1,
				"workspace.snapshots":                 1,
				"workspace.snapshots.cloud":           1,
				"workspace.snapshots.cloud.namespace": 1,
			},
			expectedCreate: map[string]bool{
				"workspace.snapshots":                 true,
				"workspace.snapshots.cloud":           true,
				"workspace.snapshots.cloud.namespace": true,
			},
		},
		{
			name:          "success - some groups exist, some need to be created",
			targetGroupId: "workspace.snapshots.cloud.namespace",
			mockClient: &mockApihubClient{
				getPackageByIdFunc: func(ctx context.Context, id string) (*view.SimplePackage, error) {
					if id == "workspace" {
						return &view.SimplePackage{
							Id:   id,
							Kind: string(view.KindWorkspace),
						}, nil
					}
					if id == "workspace.snapshots" {
						return &view.SimplePackage{
							Id:   id,
							Kind: string(view.KindGroup),
						}, nil
					}
					// cloud and namespace don't exist
					return nil, nil
				},
				createPackageFunc: func(ctx context.Context, pkg view.PackageCreateRequest) (string, error) {
					return pkg.ParentId + "." + pkg.Alias, nil
				},
			},
			expectedCalls: map[string]int{
				"workspace":                           1,
				"workspace.snapshots":                 1,
				"workspace.snapshots.cloud":           1,
				"workspace.snapshots.cloud.namespace": 1,
			},
			expectedCreate: map[string]bool{
				"workspace.snapshots.cloud":           true,
				"workspace.snapshots.cloud.namespace": true,
			},
		},
		{
			name:          "error - workspace does not exist",
			targetGroupId: "workspace.snapshots.cloud.namespace",
			mockClient: &mockApihubClient{
				getPackageByIdFunc: func(ctx context.Context, id string) (*view.SimplePackage, error) {
					if id == "workspace" {
						return nil, nil
					}
					return nil, nil
				},
			},
			expectedError: "workspace workspace does not exist",
		},
		{
			name:          "error - GetPackageById returns error for workspace",
			targetGroupId: "workspace.snapshots.cloud.namespace",
			mockClient: &mockApihubClient{
				getPackageByIdFunc: func(ctx context.Context, id string) (*view.SimplePackage, error) {
					if id == "workspace" {
						return nil, errors.New("network error")
					}
					return nil, nil
				},
			},
			expectedError: "unable to get group workspace: network error",
		},
		{
			name:          "error - GetPackageById returns error for group",
			targetGroupId: "workspace.snapshots.cloud.namespace",
			mockClient: &mockApihubClient{
				getPackageByIdFunc: func(ctx context.Context, id string) (*view.SimplePackage, error) {
					if id == "workspace" {
						return &view.SimplePackage{
							Id:   id,
							Kind: string(view.KindWorkspace),
						}, nil
					}
					if id == "workspace.snapshots" {
						return nil, errors.New("database error")
					}
					return nil, nil
				},
			},
			expectedError: "unable to get group workspace.snapshots: database error",
		},
		{
			name:          "error - CreatePackage returns error",
			targetGroupId: "workspace.snapshots.cloud.namespace",
			mockClient: &mockApihubClient{
				getPackageByIdFunc: func(ctx context.Context, id string) (*view.SimplePackage, error) {
					if id == "workspace" {
						return &view.SimplePackage{
							Id:   id,
							Kind: string(view.KindWorkspace),
						}, nil
					}
					return nil, nil
				},
				createPackageFunc: func(ctx context.Context, pkg view.PackageCreateRequest) (string, error) {
					return "", errors.New("permission denied")
				},
			},
			expectedError: "permission denied",
		},
		{
			name:          "error - existing package is not a group",
			targetGroupId: "workspace.snapshots.cloud.namespace",
			mockClient: &mockApihubClient{
				getPackageByIdFunc: func(ctx context.Context, id string) (*view.SimplePackage, error) {
					if id == "workspace" {
						return &view.SimplePackage{
							Id:   id,
							Kind: string(view.KindWorkspace),
						}, nil
					}
					if id == "workspace.snapshots" {
						return &view.SimplePackage{
							Id:   id,
							Kind: string(view.KindPackage), // Wrong kind
						}, nil
					}
					return nil, nil
				},
			},
			expectedError: "package workspace.snapshots exists but is not a group (kind: package)",
		},
		{
			name:          "success - single level group (workspace only)",
			targetGroupId: "workspace",
			mockClient: &mockApihubClient{
				getPackageByIdFunc: func(ctx context.Context, id string) (*view.SimplePackage, error) {
					return &view.SimplePackage{
						Id:   id,
						Kind: string(view.KindWorkspace),
					}, nil
				},
			},
			expectedCalls: map[string]int{
				"workspace": 1,
			},
			expectedCreate: map[string]bool{},
		},
		{
			name:          "success - two level hierarchy",
			targetGroupId: "workspace.snapshots",
			mockClient: &mockApihubClient{
				getPackageByIdFunc: func(ctx context.Context, id string) (*view.SimplePackage, error) {
					if id == "workspace" {
						return &view.SimplePackage{
							Id:   id,
							Kind: string(view.KindWorkspace),
						}, nil
					}
					return nil, nil
				},
				createPackageFunc: func(ctx context.Context, pkg view.PackageCreateRequest) (string, error) {
					return pkg.ParentId + "." + pkg.Alias, nil
				},
			},
			expectedCalls: map[string]int{
				"workspace":           1,
				"workspace.snapshots": 1,
			},
			expectedCreate: map[string]bool{
				"workspace.snapshots": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Track calls to GetPackageById
			getPackageByIdCalls := make(map[string]int)
			createPackageCalls := make(map[string]bool)

			// Wrap the mock functions to track calls
			originalGetPackageById := tt.mockClient.getPackageByIdFunc
			tt.mockClient.getPackageByIdFunc = func(ctx context.Context, id string) (*view.SimplePackage, error) {
				getPackageByIdCalls[id]++
				if originalGetPackageById != nil {
					return originalGetPackageById(ctx, id)
				}
				return nil, nil
			}

			originalCreatePackage := tt.mockClient.createPackageFunc
			tt.mockClient.createPackageFunc = func(ctx context.Context, pkg view.PackageCreateRequest) (string, error) {
				groupId := pkg.ParentId + "." + pkg.Alias
				createPackageCalls[groupId] = true
				if originalCreatePackage != nil {
					return originalCreatePackage(ctx, pkg)
				}
				return "", nil
			}

			service := &snapshotServiceImpl{
				apihubClient: tt.mockClient,
			}

			err := service.createGroupHierarchy(ctx, tt.targetGroupId)

			// Check error
			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error %q, but got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error %q, but got %q", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got %v", err)
				}
			}

			// Check GetPackageById calls
			if tt.expectedCalls != nil {
				for groupId, expectedCount := range tt.expectedCalls {
					actualCount := getPackageByIdCalls[groupId]
					if actualCount != expectedCount {
						t.Errorf("Expected %d calls to GetPackageById for %q, but got %d", expectedCount, groupId, actualCount)
					}
				}
			}

			// Check CreatePackage calls
			if tt.expectedCreate != nil {
				for groupId, shouldCreate := range tt.expectedCreate {
					wasCreated := createPackageCalls[groupId]
					if shouldCreate && !wasCreated {
						t.Errorf("Expected CreatePackage to be called for %q, but it wasn't", groupId)
					}
					if !shouldCreate && wasCreated {
						t.Errorf("Expected CreatePackage not to be called for %q, but it was", groupId)
					}
				}
			}
		})
	}
}
