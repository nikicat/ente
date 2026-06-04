package models

// AdminUserDetails is the full per-user information returned by the
// `GET /admin/user` endpoint of museum.
type AdminUserDetails struct {
	User         User              `json:"user"`
	Subscription AdminSubscription `json:"subscription"`
	Details      AdminUsageDetails `json:"details"`
	Tokens       []AdminTokenInfo  `json:"tokens"`
	AuthCodes    int64             `json:"authCodes"`
}

// AdminUsageDetails mirrors museum's details.UserDetailsResponse.
type AdminUsageDetails struct {
	Usage        int64             `json:"usage"`
	FileCount    *int64            `json:"fileCount"`
	StorageBonus int64             `json:"storageBonus"`
	Subscription AdminSubscription `json:"subscription"`
	FamilyData   *AdminFamilyData  `json:"familyData"`
}

type AdminSubscription struct {
	ProductID       string `json:"productID"`
	Storage         int64  `json:"storage"`
	ExpiryTime      int64  `json:"expiryTime"`
	PaymentProvider string `json:"paymentProvider"`
	OriginalTxnID   string `json:"originalTransactionID"`
}

type AdminFamilyData struct {
	Storage    int64               `json:"storage"`
	ExpiryTime int64               `json:"expiryTime"`
	Members    []AdminFamilyMember `json:"members"`
}

type AdminFamilyMember struct {
	Email   string `json:"email"`
	Usage   int64  `json:"usage"`
	IsAdmin bool   `json:"isAdmin"`
	Status  string `json:"status"`
}

type AdminTokenInfo struct {
	CreationTime int64  `json:"creationTime"`
	LastUsedTime int64  `json:"lastUsedTime"`
	UA           string `json:"ua"`
	IsDeleted    bool   `json:"isDeleted"`
	App          string `json:"app"`
}
