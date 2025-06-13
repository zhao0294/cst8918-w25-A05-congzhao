package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zhao0294/cst8918-w25-A05-congzhao/azure"
)

func TestAzureVMWebServer(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../",
	}

	terraform.InitAndApply(t, terraformOptions)

	publicIP := terraform.Output(t, terraformOptions, "public_ip")
	resourceGroup := terraform.Output(t, terraformOptions, "resource_group_name")
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	require.NotEmpty(t, subscriptionID)

	// 1. list Virtual Machines in the resource group
	vmList, err := azure.ListVirtualMachinesForResourceGroup(resourceGroup, subscriptionID)
	require.NoError(t, err)
	require.Greater(t, len(vmList), 0, "No VM found in resource group")
	vmName := vmList[0]

	// 2. Get and verify the NICs associated with the VM
	nics, err := azure.GetVirtualMachineNics(subscriptionID, resourceGroup, vmName)
	require.NoError(t, err)
	nicCheck := len(nics) > 0

	// 3. Get and verify the Ubuntu version of the VM
	vmImage, err := azure.GetVirtualMachineImage(subscriptionID, resourceGroup, vmName)
	require.NoError(t, err)
	expectedUbuntuOffer := "0001-com-ubuntu-server-focal" // default Ubuntu 20.04 LTS offer
	ubuntuCheck := vmImage.Offer == expectedUbuntuOffer

	// 4. Verify Apache HTTP server is running (using SSH)
	privateKeyBytes, err := os.ReadFile(os.ExpandEnv("$HOME/.ssh/id_rsa"))
	require.NoError(t, err)

	keyPair := ssh.KeyPair{PrivateKey: string(privateKeyBytes)}

	vmSSH := ssh.Host{
		Hostname:    publicIP,
		SshUserName: "azureadmin",
		SshKeyPair:  &keyPair,
	}

	curlOutput, err := ssh.CheckSshCommandE(t, vmSSH, "curl -I http://localhost")
	require.NoError(t, err)
	curlCheck := strings.Contains(curlOutput, "200 OK")

	psOutput, err := ssh.CheckSshCommandE(t, vmSSH, "ps -ef | grep apache2")
	require.NoError(t, err)
	psCheck := strings.Contains(psOutput, "apache2")

	// print summary of checks
	t.Log("\n======== Test Summary ========")
	t.Logf("NIC exists and connected: %v, NICs: %v", nicCheck, nics)
	t.Logf("VM running expected Ubuntu version: %v, Offer: %s", ubuntuCheck, vmImage.Offer)
	t.Logf("Apache HTTP response check: %v", curlCheck)
	t.Logf("Apache process running check: %v", psCheck)

	// assert checks
	assert.True(t, nicCheck, "NIC check failed")
	assert.True(t, ubuntuCheck, "Ubuntu version check failed")
	assert.True(t, curlCheck, "Apache HTTP check failed")
	assert.True(t, psCheck, "Apache process check failed")

	fmt.Println("✅ Test completed with summary above.")
}
