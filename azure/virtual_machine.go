package azure

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
)

// ListVirtualMachinesForResourceGroup 返回指定订阅和资源组中的所有虚拟机名称
func ListVirtualMachinesForResourceGroup(resourceGroup, subscriptionID string) ([]string, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get credential: %w", err)
	}

	ctx := context.Background()
	vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create vm client: %w", err)
	}

	pager := vmClient.NewListPager(resourceGroup, nil)

	vmNames := []string{}
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get VM page: %w", err)
		}

		for _, vm := range page.Value {
			if vm.Name != nil {
				vmNames = append(vmNames, *vm.Name)
			}
		}
	}

	return vmNames, nil
}

// GetVirtualMachineNics 查询指定订阅和资源组里某虚拟机关联的NIC名称列表
func GetVirtualMachineNics(subscriptionID, resourceGroup, vmName string) ([]string, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get credential: %w", err)
	}

	ctx := context.Background()
	vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create vm client: %w", err)
	}

	vmResp, err := vmClient.Get(ctx, resourceGroup, vmName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get VM: %w", err)
	}

	nicNames := []string{}
	if vmResp.Properties == nil || vmResp.Properties.NetworkProfile == nil || vmResp.Properties.NetworkProfile.NetworkInterfaces == nil {
		return nil, fmt.Errorf("vm network interfaces not found")
	}

	for _, nicRef := range vmResp.Properties.NetworkProfile.NetworkInterfaces {
		if nicRef.ID == nil {
			continue
		}
		// nicRef.ID 格式: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/networkInterfaces/{nicName}
		parts := strings.Split(*nicRef.ID, "/")
		if len(parts) > 0 {
			nicName := parts[len(parts)-1]
			nicNames = append(nicNames, nicName)
		}
	}

	return nicNames, nil
}

// VMImage 结构体用于描述虚拟机镜像信息
type VMImage struct {
	Publisher string
	Offer     string
	SKU       string
	Version   string
}

// GetVirtualMachineImage 查询指定 VM 使用的镜像信息
func GetVirtualMachineImage(subscriptionID, resourceGroup, vmName string) (*VMImage, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get credential: %w", err)
	}

	ctx := context.Background()
	vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create vm client: %w", err)
	}

	vmResp, err := vmClient.Get(ctx, resourceGroup, vmName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get VM: %w", err)
	}

	if vmResp.Properties == nil || vmResp.Properties.StorageProfile == nil || vmResp.Properties.StorageProfile.ImageReference == nil {
		return nil, fmt.Errorf("vm image reference not found")
	}

	imgRef := vmResp.Properties.StorageProfile.ImageReference

	image := &VMImage{
		Publisher: "",
		Offer:     "",
		SKU:       "",
		Version:   "",
	}
	if imgRef.Publisher != nil {
		image.Publisher = *imgRef.Publisher
	}
	if imgRef.Offer != nil {
		image.Offer = *imgRef.Offer
	}
	if imgRef.SKU != nil {
		image.SKU = *imgRef.SKU
	}
	if imgRef.Version != nil {
		image.Version = *imgRef.Version
	}

	return image, nil

}
