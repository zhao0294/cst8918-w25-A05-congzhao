package azure

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
)

// ListVirtualMachinesForResourceGroup retrieves a list of virtual machine names in the specified resource group and subscription.
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

// GetVirtualMachineNics retrieves the names of network interfaces (NICs) associated with a specified virtual machine.
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
		// nicRef.ID format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/networkInterfaces/{nicName}
		parts := strings.Split(*nicRef.ID, "/")
		if len(parts) > 0 {
			nicName := parts[len(parts)-1]
			nicNames = append(nicNames, nicName)
		}
	}

	return nicNames, nil
}

// VMImage struct represents the image information of a virtual machine.
type VMImage struct {
	Publisher string
	Offer     string
	SKU       string
	Version   string
}

// GetVirtualMachineImage retrieves the image reference of a specified virtual machine.
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
