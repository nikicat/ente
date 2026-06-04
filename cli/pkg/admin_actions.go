package pkg

import (
	"context"
	"fmt"
	"github.com/ente-io/cli/internal"
	"github.com/ente-io/cli/internal/api"
	"github.com/ente-io/cli/pkg/model"
	"github.com/ente-io/cli/utils"
	"log"
	"strings"
	"time"
)

func (c *ClICtrl) GetUserId(ctx context.Context, params model.AdminActionForUser) error {
	accountCtx, err := c.buildAdminContext(ctx, params.AdminEmail)
	if err != nil {
		return err
	}
	id, err := c.Client.GetUserIdFromEmail(accountCtx, params.UserEmail)
	if err != nil {
		return err
	}
	fmt.Println(id.User.ID)
	return nil
}

func (c *ClICtrl) ListUsers(ctx context.Context, params model.AdminActionForUser) error {
	accountCtx, err := c.buildAdminContext(ctx, params.AdminEmail)
	if err != nil {
		return err
	}
	users, err := c.Client.ListUsers(accountCtx)
	if err != nil {
		if apiErr, ok := err.(*api.ApiError); ok && apiErr.StatusCode == 400 && strings.Contains(apiErr.Message, "Token is too old") {
			fmt.Printf("Error: old admin token, please re-authenticate using `ente account add` \n")
			return nil
		}
		return err
	}
	for _, user := range users {
		fmt.Printf("Email: %s, ID: %d, Created: %s\n", user.Email, user.ID, time.UnixMicro(user.CreationTime).Format("2006-01-02"))
	}
	return nil
}

func (c *ClICtrl) GetUserDetails(ctx context.Context, params model.AdminActionForUser) error {
	accountCtx, err := c.buildAdminContext(ctx, params.AdminEmail)
	if err != nil {
		return err
	}
	details, err := c.Client.GetUserDetails(accountCtx, params.UserEmail)
	if err != nil {
		if apiErr, ok := err.(*api.ApiError); ok && apiErr.StatusCode == 400 && strings.Contains(apiErr.Message, "Token is too old") {
			fmt.Printf("Error: old admin token, please re-authenticate using `ente account add` \n")
			return nil
		}
		return err
	}

	sub := details.Subscription
	usage := details.Details

	fmt.Printf("User\n")
	fmt.Printf("  Email:    %s\n", details.User.Email)
	fmt.Printf("  ID:       %d\n", details.User.ID)
	if details.User.CreationTime > 0 {
		fmt.Printf("  Created:  %s\n", time.UnixMicro(details.User.CreationTime).Format("2006-01-02 15:04:05"))
	}

	fmt.Printf("Usage\n")
	fmt.Printf("  Used:     %s\n", utils.ByteCountDecimalGIB(usage.Usage))
	if usage.FileCount != nil {
		fmt.Printf("  Files:    %d\n", *usage.FileCount)
	}
	if usage.StorageBonus > 0 {
		fmt.Printf("  Bonus:    %s\n", utils.ByteCountDecimalGIB(usage.StorageBonus))
	}

	fmt.Printf("Subscription\n")
	fmt.Printf("  Plan:     %s\n", sub.ProductID)
	fmt.Printf("  Storage:  %s\n", utils.ByteCountDecimalGIB(sub.Storage))
	if sub.ExpiryTime > 0 {
		fmt.Printf("  Expiry:   %s\n", time.UnixMicro(sub.ExpiryTime).Format("2006-01-02"))
	}
	if sub.PaymentProvider != "" {
		fmt.Printf("  Provider: %s\n", sub.PaymentProvider)
	}

	if usage.FamilyData != nil && len(usage.FamilyData.Members) > 0 {
		fmt.Printf("Family (storage %s)\n", utils.ByteCountDecimalGIB(usage.FamilyData.Storage))
		for _, m := range usage.FamilyData.Members {
			adminMarker := ""
			if m.IsAdmin {
				adminMarker = " (admin)"
			}
			fmt.Printf("  %s%s - %s [%s]\n", m.Email, adminMarker, utils.ByteCountDecimalGIB(m.Usage), m.Status)
		}
	}

	fmt.Printf("Other\n")
	fmt.Printf("  Auth codes:    %d\n", details.AuthCodes)
	activeTokens := 0
	for _, t := range details.Tokens {
		if !t.IsDeleted {
			activeTokens++
		}
	}
	fmt.Printf("  Active tokens: %d\n", activeTokens)

	return nil
}

func (c *ClICtrl) DeleteUser(ctx context.Context, params model.AdminActionForUser) error {
	accountCtx, err := c.buildAdminContext(ctx, params.AdminEmail)
	if err != nil {
		return err
	}
	err = c.Client.DeleteUser(accountCtx, params.UserEmail)
	if err != nil {
		if apiErr, ok := err.(*api.ApiError); ok && apiErr.StatusCode == 400 && strings.Contains(apiErr.Message, "Token is too old") {
			fmt.Printf("Error: old admin token, please re-authenticate using `ente account add` \n")
			return nil
		}
		return err
	}
	fmt.Println("Successfully deleted user")
	return nil
}

func (c *ClICtrl) Disable2FA(ctx context.Context, params model.AdminActionForUser) error {
	accountCtx, err := c.buildAdminContext(ctx, params.AdminEmail)
	if err != nil {
		return err
	}
	userDetails, err := c.Client.GetUserIdFromEmail(accountCtx, params.UserEmail)
	if err != nil {
		return err
	}
	err = c.Client.Disable2Fa(accountCtx, userDetails.User.ID)
	if err != nil {
		if apiErr, ok := err.(*api.ApiError); ok && apiErr.StatusCode == 400 && strings.Contains(apiErr.Message, "Token is too old") {
			fmt.Printf("Error: Old admin token, please re-authenticate using `ente account add` \n")
			return nil
		}
		return err
	}
	fmt.Println("Successfully disabled 2FA for user")
	return nil
}

func (c *ClICtrl) DisablePasskeys(ctx context.Context, params model.AdminActionForUser) error {
	accountCtx, err := c.buildAdminContext(ctx, params.AdminEmail)
	if err != nil {
		return err
	}
	userDetails, err := c.Client.GetUserIdFromEmail(accountCtx, params.UserEmail)
	if err != nil {
		return err
	}
	err = c.Client.DisablePassKeyMFA(accountCtx, userDetails.User.ID)
	if err != nil {
		if apiErr, ok := err.(*api.ApiError); ok && apiErr.StatusCode == 400 && strings.Contains(apiErr.Message, "Token is too old") {
			fmt.Printf("Error: Old admin token, please re-authenticate using `ente account add` \n")
			return nil
		}
		return err
	}
	fmt.Println("Successfully disabled passkey for user")
	return nil
}

func (c *ClICtrl) UpdateFreeStorage(ctx context.Context, params model.AdminActionForUser, noLimit bool) error {
	accountCtx, err := c.buildAdminContext(ctx, params.AdminEmail)
	if err != nil {
		return err
	}
	userDetails, err := c.Client.GetUserIdFromEmail(accountCtx, params.UserEmail)
	if err != nil {
		return err
	}
	if noLimit {
		// set storage to 100TB and expiry to + 100 years
		err := c.Client.UpdateFreePlanSub(accountCtx, userDetails, 100*1024*1024*1024*1024, time.Now().AddDate(100, 0, 0).UnixMicro())
		if err != nil {
			return err
		} else {
			fmt.Println("Successfully updated storage and expiry date for user")
		}
		return nil
	}
	storageSize, err := internal.GetStorageSize("Enter a storage size (e.g.'5MB', '10GB', '2Tb'): ")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	dateStr, err := internal.GetUserInput("Enter sub expiry date in YYYY-MM-DD format  (e.g.'2040-12-31')")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	date, err := _parseDateOrDateTime(dateStr)
	if err != nil {
		return err
	}

	fmt.Printf("Updating storage for user %s to %s (old %s) with new expiry %s (old %s) \n",
		params.UserEmail,
		utils.ByteCountDecimalGIB(storageSize), utils.ByteCountDecimalGIB(userDetails.Subscription.Storage),
		date.Format("2006-01-02"),
		time.UnixMicro(userDetails.Subscription.ExpiryTime).Format("2006-01-02"))
	// press y to confirm
	confirmed, _ := internal.ConfirmAction("Are you sure you want to update the storage ('y' or 'n')?")
	if !confirmed {
		return nil
	} else {
		err := c.Client.UpdateFreePlanSub(accountCtx, userDetails, storageSize, date.UnixMicro())
		if err != nil {
			return err
		} else {
			fmt.Println("Successfully updated storage and expiry date for user")
		}
	}

	return nil
}

func (c *ClICtrl) SendTestMail(ctx context.Context, params model.AdminActionForUser, to, from, fromName string) error {
	accountCtx, err := c.buildAdminContext(ctx, params.AdminEmail)
	if err != nil {
		return err
	}
	err = c.Client.SendTestMail(accountCtx, to, from, fromName)
	if err != nil {
		if apiErr, ok := err.(*api.ApiError); ok && apiErr.StatusCode == 400 && strings.Contains(apiErr.Message, "Token is too old") {
			fmt.Printf("Error: old admin token, please re-authenticate using `ente account add` \n")
			return nil
		}
		return err
	}
	fmt.Printf("Successfully sent test email to %s\n", to)
	return nil
}

func (c *ClICtrl) buildAdminContext(ctx context.Context, adminEmail string) (context.Context, error) {
	accounts, err := c.GetAccounts(ctx)
	if err != nil {
		return nil, err
	}
	if len(accounts) == 0 {
		return nil, fmt.Errorf("no accounts found, use `account add` to add an account")
	}
	var acc *model.Account
	for _, a := range accounts {
		if a.Email == adminEmail {
			acc = &a
			break
		}
	}
	if (len(accounts) > 1) && (acc == nil) {
		return nil, fmt.Errorf("multiple accounts found, specify the admin email using --admin-user")
	}
	if acc == nil && len(accounts) == 1 {
		acc = &accounts[0]
		fmt.Printf("Assuming %s as the Admin \n------------\n", acc.Email)
	}

	if acc == nil {
		return nil, fmt.Errorf("account not found for %s, use `account list` to list accounts", adminEmail)
	}
	secretInfo, err := c.KeyHolder.LoadSecrets(*acc)
	if err != nil {
		return nil, err
	}
	accountCtx := c.buildRequestContext(ctx, *acc, model.Filter{})
	c.Client.AddToken(acc.AccountKey(), secretInfo.TokenStr())
	return accountCtx, nil
}

func _parseDateOrDateTime(input string) (time.Time, error) {
	var layout string
	if strings.Contains(input, " ") {
		// If the input contains a space, assume it's a date-time format
		layout = "2006-01-02 15:04:05"
	} else {
		// If there's no space, assume it's just a date
		layout = "2006-01-02"
	}
	return time.Parse(layout, input)
}
