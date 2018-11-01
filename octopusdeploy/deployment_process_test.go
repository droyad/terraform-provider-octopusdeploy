package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/MattHodge/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployDeploymentProcessBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployDeploymentProcessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentProcessBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployDeploymentProcess(),
				),
			},
		},
	})
}

func testAccDeploymentProcessBasic() string {
	return `
		resource "octopusdeploy_lifecycle" "test" {
			name = "Test Lifecycle"
		}

		resource "octopusdeploy_project_group" "test" {
			name = "Test Group"
		}

		resource "octopusdeploy_project" "test" {
			name             = "Test Project"
			lifecycle_id     = "${octopusdeploy_lifecycle.test.id}"
			project_group_id = "${octopusdeploy_project_group.test.id}"
		}

		resource "octopusdeploy_deployment_process" "test''" {
			project_id = "${octopusdeploy_project.test.id}"

			step {
				name = "Test"

				action {
					name = "Test"
					action_type = "Octopus.Script"

					property {
						key = "Octopus.Action.RunOnServer"
						value = "true"
					}
				}
			}
		}
		`
}

func testAccCheckOctopusDeployDeploymentProcessDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)

	if err := destroyProjectHelper(s, client); err != nil {
		return err
	}
	if err := destroyHelperProjectGroup(s, client); err != nil {
		return err
	}
	if err := destroyHelperLifecycle(s, client); err != nil {
		return err
	}
	return nil
}

func testAccCheckOctopusDeployDeploymentProcess() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)

		process, err := getDeploymentProcess(s, client);
		if err != nil {
			return err
		}

		expectedNumberOfSteps := 1
		numberOfSteps := len(process.Steps)
		if numberOfSteps != expectedNumberOfSteps {
			return fmt.Errorf("Deployment process has %d steps instead of the expected %d", numberOfSteps, expectedNumberOfSteps)
		}

		if process.Steps[0].Actions[0].Properties["Octopus.Action.RunOnServer"] != "true" {
			return fmt.Errorf("The RunOnServer property has not been set to true on the deployment process")
		}

		return nil;
	}
}


func getDeploymentProcess(s *terraform.State, client *octopusdeploy.Client) (*octopusdeploy.DeploymentProcess, error) {
	for _, r := range s.RootModule().Resources {
		if r.Type == "octopusdeploy_deployment_process" {
			return client.DeploymentProcess.Get(r.Primary.ID);
		}
	}
	return nil, fmt.Errorf("No deployment process found in the terraform resources")
}
