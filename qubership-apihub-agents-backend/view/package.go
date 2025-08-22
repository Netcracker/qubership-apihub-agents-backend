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

package view

import "time"

type PackageKind string

const KindPackage PackageKind = "package"
const KindWorkspace PackageKind = "workspace"
const KindGroup PackageKind = "group"
const KindDashbord PackageKind = "dashboard"

type SimplePackage struct {
	Id                    string              `json:"packageId"`
	Alias                 string              `json:"alias"`
	ParentId              string              `json:"parentId"`
	Kind                  string              `json:"kind"`
	Name                  string              `json:"name"`
	Description           string              `json:"description"`
	IsFavorite            bool                `json:"isFavorite"`
	ServiceName           string              `json:"serviceName,omitempty"`
	ImageUrl              string              `json:"imageUrl"`
	Parents               []ParentPackageInfo `json:"parents"`
	UserRole              string              `json:"userRole"`
	DefaultRole           string              `json:"defaultRole"`
	DeletionDate          *time.Time          `json:"-"`
	DeletedBy             string              `json:"-"`
	CreatedBy             string              `json:"-"`
	CreatedAt             time.Time           `json:"createdAt,omitempty"`
	ReleaseVersionPattern string              `json:"releaseVersionPattern"`
	DefaultReleaseVersion string              `json:"defaultReleaseVersion"`
}

type SimplePackages struct {
	Packages []SimplePackage `json:"packages"`
}

type ParentPackageInfo struct {
	Id       string `json:"packageId"`
	Alias    string `json:"alias"`
	ParentId string `json:"parentId"`
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	ImageUrl string `json:"imageUrl"`
}

type PackageCreateRequest struct {
	ParentId              string `json:"parentId"` //required
	Kind                  string `json:"kind"`     //required
	Name                  string `json:"name"`     //required
	Alias                 string `json:"alias"`    //required
	Description           string `json:"description"`
	ServiceName           string `json:"serviceName"`
	ImageUrl              string `json:"imageUrl"`
	DefaultRole           string `json:"defaultRole"`
	ReleaseVersionPattern string `json:"releaseVersionPattern"`
	ExcludeFromSearch     *bool  `json:"excludeFromSearch"`
}

// TODO: update!
type PackageResponse struct {
	Id          string `json:"id"`
	Alias       string `json:"alias" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Kind        string `json:"kind" validate:"required"`
	ParentId    string `json:"parentId" validate:"required"`
	Description string `json:"description"`
	ServiceName string `json:"serviceName"`
	ImageUrl    string `json:"imageUrl"`
}

type AvailablePackagePromoteStatuses map[string][]string // map[packageId][]version status

type PackagesReq struct {
	Packages []string `json:"packages"`
}

type PackagesSearchReq struct {
	ServiceName        string
	TextFilter         string
	Kind               string
	ParentId           string
	ShowAllDescendants bool
	ShowParents        bool
	Page               int
	Limit              int
}
