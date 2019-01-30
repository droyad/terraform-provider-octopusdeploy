package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/MattHodge/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployDeployKubernetesContainersAction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployDeploymentProcessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeployKubernetesContainersAction(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeployKubernetesContainers(),
				),
			},
		},
	})
}

func testAccDeployKubernetesContainersAction() string {
	return testAccBuildTestActionTerraform(`
		deploy_kubernetes_containers_action {
          name = "Test"

          deployment_name = "thedeployment"
          replicas = "1"          // Optional
          progression_deadline = "1800"          // Optional

          label {
              name = "Label1"
              value = "LabelValue"
          }

          deployment_strategy = "MaxUnavailable"          // Octopus.Action.KubernetesContainers.DeploymentStyle  Recreate, RollingUpdate, BlueGreen (No Variable)
          rolling_update_max_unavailable = "3"          // Octopus.Action.KubernetesContainers.MaxUnavailable
          rolling_update_max_surge = "3"          // Octopus.Action.KubernetesContainers.MaxSurge
          deployment_wait = "NoWait"          // Wait (default) or NoWait Octopus.Action.KubernetesContainers.DeploymentWait only if style is Recreate or RollingUpdate

          config_map_volume {
              name = "configmap"              //reference_name_type = "CustomResource" // LinkedResource (default) or CustomResource, Don't add as n option
              external_config_map_name = "external_config_map"              // If specified set reference_name_type to CustomResource

              item {
                  key = "Item1"
                  path = "path1"
              }
          }

          secret_volume {
              name = "secret"
              external_secret_resource = "external secret"              // If specified set reference_name_type to CustomResource

              item {
                  key = "Item1"
                  path = "path1"
              }
          }

          empty_dir {
              name = "emptydir"
              medium = "Memory"              // optional
          }

          host_path_volume {
              name = "hostpath"
              type = "FileOrCreate"              // optional
              path = "/tmp"
          }

          persistent_volume_claim_volume {
              name = "persistent"
              claim_name = "Abc"
          }

          container {
              name = "container1"
              feed_id = "Feeds-1"              // optional, default to built in
              package_id = "MyContainer"
              is_init_container = "False"              // default False, sets InitContainer property
              command = "run.ps1"              // optional
              command_arguments = "-console"

              // Should these be a seperate group?
              cpu_request = "3m"              // optional
              cpu_limit = "199"              // optional
              memory_request = "55G"              // optional
              memory_limit = "2"              // optional

              // Should these be a seperate group? (security Context)
              privilege_escalation = "True"              // True, False or variable
              privileged_mode = "True"              // True, False or variable
              read_only_root_filesystem = "True"              // True, False or variable
              run_as_non_root = "True"              // True, False or variable
              run_as_user = "#{User}"
              run_as_group = "#{Group}"

              add_posix_capabilities = [
                  "CAP_AUDIT_READ"]
              drop_posix_capabilities = [
                  "CAP_AUDIT_WRITE"]

              // All optional
              se_linux_level = "level1"
              se_linux_role = "roleA"
              se_linux_type = "setype"
              se_linux_user = "seuser"

              port {
                  name = "port1"                  // optional
                  port = "#{Port}"
                  protocol = "UDP"                  // optional
              }

              volume_mount {
                  name = "configmap"
                  mount_path = "/usr"
                  sub_path = "bin"                  // optional
                  read_only = "True"                  // optional
              }

              liveness_probe {                  // optional 1
                  failure_threshold = "4"                  // optional
                  timeout = "5"                  // optional
                  initial_dely = "6"                  // optional
                  period = "7"                  // optional


                  // only zero or one of the three of the following can be specified
                  health_check_command = "the_command"

                  http_get {
                      path = "/status"
                      host = "myhost"                      // optional
                      port = "80"
                      scheme = "HTTPS"                      // optional HTTP, HTTPS or variable

                      http_header {
                          name = "Head1"
                          value = "Value1"
                      }
                  }

                  tcp_socket {
                      host = "myhost"                      // optional
                      port = "80"
                  }
              }

              readiness_probe {                  // same as liveness_probe
                  failure_threshold = "4"                  // optional
                  timeout = "5"                  // optional
                  initial_dely = "6"                  // optional
                  period = "7"                  // optional


                  // only zero or one of the three of the following can be specified
                  health_check_command = "the_command"

                  http_get {
                      path = "/status"
                      host = "myhost"                      // optional
                      port = "80"
                      scheme = "HTTPS"                      // optional HTTP, HTTPS or variable

                      http_header {
                          name = "Head1"
                          value = "Value1"
                      }
                  }

                  tcp_socket {
                      host = "myhost"
                      // optional
                      port = "80"
                  }
              }

              environment_variable {
                  name = "Env1"
                  value = "value1"
              }

              config_map_environment_variable {
                  name = "a"
                  config_map_name = "b"
                  key = "c"
              }

              secret_environment_variable {
                  name = "s1"
                  secret_name = "s2"
                  key = "s3"
              }
          }

          // zero or one
          pod_security_context {
              // all items are optional
              fs_group = "1"
              // optional

              run_as_user = "#{User}"
              run_as_group = "#{Group}"
              suplimental_groups = "1,2"
              run_as_non_root = "True"
              // Default to "False"

              se_linux_level = "level1"
              se_linux_role = "roleA"
              se_linux_type = "setype"
              se_linux_user = "seuser"

              sysctl {
                  name = "sysctl1"
                  value = "value1"
              }
          }

          // zero or more
          pod_affinity_rule {
              preferred_affinity_weight = "2"
              // Optional, If set, podAffinityDetails.Type should be "Preferred", otherwise it should be "Required""
              topology_key = "topkey"
              namespaces = "D,E"
              // optional

              in_rule {
                  label_key = "key"
                  operation = "In"
                  values = "value"
              }

              exist_rule {
                  label_key = "existkey"
                  operation = "Exists"
              }
          }

          // zero or more, almost identical to pod rule
          node_affinity_rule {
              preferred_affinity_weight = "2"
              // Optional, If set, nodeAffinityDetails.Type should be "Preferred", otherwise it should be "Required""

              in_rule {
                  label_key = "key"
                  operation = "In"
                  values = "value"
              }

              exist_rule {
                  label_key = "existkey"
                  operation = "Exists"
              }
          }

          // zero or more
          pod_annotation {
              name = "podann"
              value = "value1"
          }

          // zero or more
          deployment_annotation {
              name = "depann"
              value = "value1"
          }

          ingress {
              // optional
              name = "ingress"
              service_name = "ingress-service"
              service_type = "ClusterIP"               // Default to "NodePort", other values are "ClusterIP" and "LoadBalancer", no variable
              cluster_ip_address = "4.4.4.4"              // optional
              load_balancer_ip_address = "3.3.3.3"              // optional

              // zero or more
              annotation {
                  name = "ingannot"
                  value = "value1"
              }

              // one or more
              host_rule {
                  host = "myhost"

                  // one or more
                  path {
                      path = "/path"
                      service_port = "777"
                  }
              }

              // one or more
              service_port {
                  name = "svcport"
                  port = "999"
                  target_port = "888"
                  //optional
                  node_port = "444"
                  // Optional
                  protocol = "UDP"
                  // Default to TCP
              }

              // zero or more
              service_annotation {
                  name = "ingannot"
                  value = "value1"
              }
          }

          config_map_name = "CFGMAP"

          config_map_item {
              key = "cfgmapitem"
              value = "value1"
          }

          secret_name = "scrt"

          scrt_item {
              key = "scrtitem"
              value = "value1"
          }

          custom_resource_yaml = <<EOF
kind: NetworkPolicy
metadata:
    name: test-network-policy"
EOF

      }
	`)
}

func testAccCheckDeployKubernetesContainers() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)

		process, err := getDeploymentProcess(s, client);
		if err != nil {
			return err
		}

		action := process.Steps[0].Actions[0]

		if action.ActionType != expectedActionType {
			return fmt.Errorf("Action type is incorrect: %s, expected: %s", action.ActionType, expectedActionType)
		}

		if(len(action.Packages) == 0) {
			return fmt.Errorf("No package")
		}

		if action.Properties["Octopus.Action.WindowsService.CreateOrUpdateService"] != "True" {
			return fmt.Errorf("Windows Service feature is not enabled")
		}

		if action.Properties["Octopus.Action.WindowsService.ServiceName"] != "MyService" {
			return fmt.Errorf("Service Name is incorrect: %s", action.Properties["Octopus.Action.WindowsService.ServiceName"])
		}

		if action.Properties["Octopus.Action.WindowsService.DisplayName"] != "My Service" {
			return fmt.Errorf("Display Name is incorrect: %s", action.Properties["Octopus.Action.WindowsService.DisplayName"])
		}

		if action.Properties["Octopus.Action.WindowsService.Description"] != "Do stuff" {
			return fmt.Errorf("Description is incorrect: %s", action.Properties["Octopus.Action.WindowsService.Description"])
		}

		if action.Properties["Octopus.Action.WindowsService.ExecutablePath"] != "MyService.exe" {
			return fmt.Errorf("Executable Path is incorrect: %s", action.Properties["Octopus.Action.WindowsService.ExecutablePath"])
		}

		if action.Properties["Octopus.Action.WindowsService.Arguments"] != "-arg" {
			return fmt.Errorf("Arguments is incorrect: %s", action.Properties["Octopus.Action.WindowsService.Arguments"])
		}

		if action.Properties["Octopus.Action.WindowsService.ServiceAccount"] != "_CUSTOM" {
			return fmt.Errorf("Service Account is incorrect: %s", action.Properties["Octopus.Action.WindowsService.ServiceAccount"])
		}

		if action.Properties["Octopus.Action.WindowsService.CustomAccountName"] != "User" {
			return fmt.Errorf("Custom Account Name is incorrect: %s", action.Properties["Octopus.Action.WindowsService.CustomAccountName"])
		}

		if action.Properties["Octopus.Action.WindowsService.CustomAccountPassword"] != "Password" {
			return fmt.Errorf("Custom Account Password is incorrect: %s", action.Properties["Octopus.Action.WindowsService.CustomAccountPassword"])
		}

		if action.Properties["Octopus.Action.WindowsService.StartMode"] != "manual" {
			return fmt.Errorf("Start Mode is incorrect: %s", action.Properties["Octopus.Action.WindowsService.StartMode"])
		}

		if action.Properties["Octopus.Action.WindowsService.Dependencies"] != "OtherService" {
			return fmt.Errorf("Dependencies is incorrect: %s", action.Properties["Octopus.Action.WindowsService.Dependencies"])
		}

		return nil;
	}
}
