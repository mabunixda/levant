// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/levant/test/acctest"
	"github.com/hashicorp/nomad/api"
)

func TestDeploy_basic(t *testing.T) {
	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: acctest.DeployTestStepRunner{
					FixtureName: "deploy_basic.nomad",
				},
				Check: acctest.CheckDeploymentStatus("successful"),
			},
		},
		CleanupFunc: acctest.CleanupPurgeJob,
	})
}

func TestDeploy_namespace(t *testing.T) {
	acctest.Test(t, acctest.TestCase{
		SetupFunc: func(s *acctest.TestState) error {
			if _, err := s.Nomad.Namespaces().Register(&api.Namespace{Name: "test"}, nil); err != nil {
				return fmt.Errorf("could not create test namespace: %w", err)
			}
			return nil
		},
		Steps: []acctest.TestStep{
			{
				Runner: acctest.DeployTestStepRunner{
					FixtureName: "deploy_namespace.nomad",
				},
				Check: acctest.CheckDeploymentStatus("successful"),
			},
		},
		CleanupFunc: func(s *acctest.TestState) error {
			if err := acctest.CleanupPurgeJob(s); err != nil {
				return err
			}

			if _, err := s.Nomad.Namespaces().Delete("test", nil); err != nil {
				return fmt.Errorf("could not delete namespace test: %w", err)
			}
			return nil
		},
	})
}

func TestDeploy_driverError(t *testing.T) {
	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: acctest.DeployTestStepRunner{
					FixtureName: "deploy_driver_error.nomad",
				},
				ExpectErr: true,
				CheckErr: func(_ error) bool {
					// this is a bit pointless without the error bubbled up from levant
					return true
				},
			},
			{
				// allows us to check a job was registered and previous step error wasn't a parse failure etc.
				Check: acctest.CheckDeploymentStatus("failed"),
			},
		},
		CleanupFunc: acctest.CleanupPurgeJob,
	})
}

func TestDeploy_allocError(t *testing.T) {
	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: acctest.DeployTestStepRunner{
					FixtureName: "deploy_alloc_error.nomad",
				},
				ExpectErr: true,
				CheckErr: func(_ error) bool {
					// this is a bit pointless without the error bubbled up from levant
					return true
				},
			},
			{
				// allows us to check a job was registered and previous step error wasn't a parse failure etc.
				Check: acctest.CheckDeploymentStatus("failed"),
			},
		},
		CleanupFunc: acctest.CleanupPurgeJob,
	})
}

func TestDeploy_count(t *testing.T) {
	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: acctest.DeployTestStepRunner{
					FixtureName: "deploy_count.nomad",
					Vars: map[string]interface{}{
						"count": "3",
					},
				},
				Check: acctest.CheckDeploymentStatus("successful"),
			},
			{
				Runner: acctest.DeployTestStepRunner{
					FixtureName: "deploy_count.nomad",
					Vars: map[string]interface{}{
						"count": "1",
					},
				},
				Check: acctest.CheckDeploymentStatus("successful"),
			},
			{
				// expect levant to read counts from the api
				Check: acctest.CheckTaskGroupCount("test", 3),
			},
			{
				Runner: acctest.DeployTestStepRunner{
					FixtureName: "deploy_count.nomad",
					Vars: map[string]interface{}{
						"count": "1",
					},
					ForceCount: true,
				},
				Check: acctest.CheckDeploymentStatus("successful"),
			},
			{
				Check: acctest.CheckTaskGroupCount("test", 1),
			},
		},
		CleanupFunc: acctest.CleanupPurgeJob,
	})
}

func TestDeploy_canary(t *testing.T) {
	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: acctest.DeployTestStepRunner{
					FixtureName: "deploy_canary.nomad",
					Canary:      10,
					Vars: map[string]interface{}{
						"env_version": "1",
					},
				},
				Check: acctest.CheckDeploymentStatus("successful"),
			},
			{
				Runner: acctest.DeployTestStepRunner{
					FixtureName: "deploy_canary.nomad",
					Canary:      10,
					Vars: map[string]interface{}{
						"env_version": "2",
					},
				},
				Check: acctest.CheckDeploymentStatus("successful"),
			},
		},
		CleanupFunc: acctest.CleanupPurgeJob,
	})
}

func TestDeploy_lifecycle(t *testing.T) {
	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: acctest.DeployTestStepRunner{
					FixtureName: "deploy_lifecycle.nomad",
				},
				Check: acctest.CheckDeploymentStatus("successful"),
			},
		},
		CleanupFunc: acctest.CleanupPurgeJob,
	})
}

func TestDeploy_taskScalingStanza(t *testing.T) {
	acctest.Test(t, acctest.TestCase{
		Steps: []acctest.TestStep{
			{
				Runner: acctest.DeployTestStepRunner{
					FixtureName: "deploy_task_scaling.nomad",
				},
				Check: acctest.CheckDeploymentStatus("successful"),
			},
		},
		CleanupFunc: acctest.CleanupPurgeJob,
	})
}
