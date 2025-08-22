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
	"time"

	"github.com/Netcracker/qubership-apihub-agents-backend/client"
	"github.com/Netcracker/qubership-apihub-agents-backend/secctx"
	"github.com/Netcracker/qubership-apihub-agents-backend/view"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

type CleanupService interface {
	CreateDraftsCleanupJob(schedule string) error
}

func NewCleanupService(apihubClient client.ApihubClient) CleanupService {
	cronInstance := cron.New()
	cronInstance.Start()
	return &cleanupServiceImpl{
		apihubClient: apihubClient,
		cronInstance: cronInstance,
	}
}

type cleanupServiceImpl struct {
	apihubClient client.ApihubClient
	cronInstance *cron.Cron
}

type draftsCleanupJob struct {
	schedule     string
	apihubClient client.ApihubClient
}

func (c *cleanupServiceImpl) CreateDraftsCleanupJob(schedule string) error {
	job := draftsCleanupJob{
		schedule:     schedule,
		apihubClient: c.apihubClient,
	}
	_, err := c.cronInstance.AddJob(schedule, &job)
	if err != nil {
		log.Warnf("Drafts cleanup job wasn't added for schedule - %s. With error - %s", schedule, err)
		return err
	}
	log.Infof("Drafts cleanup job was created with schedule - %s", schedule)

	return nil
}

func (j draftsCleanupJob) Run() {
	//TODO: implement calculation of the timeout based on the schedule
	ctx := secctx.MakeSysadminContext(context.Background())
	workspaces, err := j.apihubClient.GetPackages(ctx, view.PackagesSearchReq{
		Kind: string(view.KindWorkspace),
	})
	if err != nil {
		log.Errorf("[SnapshotsCleanup] failed to get workspaces")
		return
	}
	//TODO: add a configurable TTL parameter
	retention := time.Now().AddDate(0, -6, 0)
	log.Infof("[SnapshotsCleanup] Starting deleting snapshots published earlier than %s", retention)
	for _, workspace := range workspaces.Packages {
		snapshotsGroupId := workspace.Id + "." + view.DefaultSnapshotsGroupAlias
		pkg, err := j.apihubClient.GetPackageById(ctx, snapshotsGroupId)
		if err != nil {
			log.Errorf("[SnapshotsCleanup] Failed to check '%s' runenv group existence: %s", snapshotsGroupId, err.Error())
			continue
		}
		if pkg == nil {
			continue
		}

		jobId, err := j.apihubClient.DeleteVersionsRecursively(ctx, snapshotsGroupId, view.DeleteVersionsRecursivelyReq{OlderThanDate: retention})
		if err != nil {
			log.Errorf("[SnapshotsCleanup] Failed to delete old snapshots in group %s: %s", snapshotsGroupId, err.Error())
			continue
		}
		if jobId == "" {
			log.Infof("[SnapshotsCleanup] Snapshots group %s not found", snapshotsGroupId)
			continue
		}
		log.Infof("[SnapshotsCleanup] Cleanup snapshots for group %s has been successfully started with job id %s", snapshotsGroupId, jobId)
	}
}
