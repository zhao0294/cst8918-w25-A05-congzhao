package test

import (
	"fmt"
	"os"
)

// GetTargetAzureSubscriptionE retrieves the Azure Subscription ID from the environment variable ARM_SUBSCRIPTION_ID.
func GetTargetAzureSubscriptionE() (string, error) {
	subID := os.Getenv("ARM_SUBSCRIPTION_ID")
	if subID == "" {
		return "", fmt.Errorf("ARM_SUBSCRIPTION_ID environment variable is not set")
	}
	return subID, nil
}
