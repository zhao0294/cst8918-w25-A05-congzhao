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

	// 1. 列出虚拟机
	vmList, err := azure.ListVirtualMachinesForResourceGroup(resourceGroup, subscriptionID)
	require.NoError(t, err)
	require.Greater(t, len(vmList), 0, "No VM found in resource group")
	vmName := vmList[0]

	// 2. 获取并验证NIC
	nics, err := azure.GetVirtualMachineNics(subscriptionID, resourceGroup, vmName)
	require.NoError(t, err)
	nicCheck := len(nics) > 0

	// 3. 获取并验证虚拟机镜像（Ubuntu版本）
	vmImage, err := azure.GetVirtualMachineImage(subscriptionID, resourceGroup, vmName)
	require.NoError(t, err)
	expectedUbuntuOffer := "0001-com-ubuntu-server-focal" // 根据你的实际Ubuntu版本调整
	ubuntuCheck := vmImage.Offer == expectedUbuntuOffer

	// 4. SSH验证Apache是否运行
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

	// 统一打印测试总结
	t.Log("\n======== Test Summary ========")
	t.Logf("NIC exists and connected: %v, NICs: %v", nicCheck, nics)
	t.Logf("VM running expected Ubuntu version: %v, Offer: %s", ubuntuCheck, vmImage.Offer)
	t.Logf("Apache HTTP response check: %v", curlCheck)
	t.Logf("Apache process running check: %v", psCheck)

	// 断言所有检查是否通过
	assert.True(t, nicCheck, "NIC check failed")
	assert.True(t, ubuntuCheck, "Ubuntu version check failed")
	assert.True(t, curlCheck, "Apache HTTP check failed")
	assert.True(t, psCheck, "Apache process check failed")

	fmt.Println("✅ Test completed with summary above.")
}
