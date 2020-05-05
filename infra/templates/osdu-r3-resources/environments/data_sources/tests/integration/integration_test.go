//  Copyright © Microsoft Corporation
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	cosmosIntegTests "github.com/microsoft/cobalt/infra/modules/providers/azure/cosmosdb/tests/integration"
	redisIntegTests "github.com/microsoft/cobalt/infra/modules/providers/azure/redis-cache/tests/integration"
	storageIntegTests "github.com/microsoft/cobalt/infra/modules/providers/azure/storage-account/tests/integration"
	esIntegTestConfig "github.com/microsoft/cobalt/infra/modules/providers/elastic/elastic-cloud-enterprise/tests"
	esIntegTests "github.com/microsoft/cobalt/infra/modules/providers/elastic/elastic-cloud-enterprise/tests/integration"
	"github.com/microsoft/cobalt/test-harness/infratests"
)

var subscription = os.Getenv("ARM_SUBSCRIPTION_ID")
var tfOptions = &terraform.Options{
	TerraformDir: "../../",
	BackendConfig: map[string]interface{}{
		"storage_account_name": os.Getenv("TF_VAR_remote_state_account"),
		"container_name":       os.Getenv("TF_VAR_remote_state_container"),
	},
}

// Runs a suite of test assertions to validate that a provisioned data source environment
// is fully functional.
func TestDataEnvironment(t *testing.T) {
	esIntegTestConfig.ESVersion = "6.8.3"
	testFixture := infratests.IntegrationTestFixture{
		GoTest:                t,
		TfOptions:             tfOptions,
		ExpectedTfOutputCount: 12,
		TfOutputAssertions: []infratests.TerraformOutputValidation{
			esIntegTests.CheckClusterHealth("elastic_cluster_properties"),
			esIntegTests.CheckClusterVersion("elastic_cluster_properties"),
			esIntegTests.CheckClusterIndexing("elastic_cluster_properties"),
			redisIntegTests.CheckRedisWriteOperations("redis_hostname", "redis_primary_access_key", "redis_port"),
			redisIntegTests.InspectProvisionedCache("redis_name", "data_resource_group_name"),
			storageIntegTests.InspectStorageAccount("storage_account", "storage_account_containers", "data_resource_group_name"),
			cosmosIntegTests.InspectProvisionedCosmosDBAccount("data_resource_group_name", "cosmosdb_account_name", "cosmosdb_properties"),
		},
	}
	infratests.RunIntegrationTests(&testFixture)
}
