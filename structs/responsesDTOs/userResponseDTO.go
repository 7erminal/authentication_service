package responsesDTOs

import (
	"time"
)

type Currencies struct {
	CurrencyId   string
	Symbol       string
	Currency     string
	Active       int
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    int
	ModifiedBy   int
}

type Countries struct {
	CountryId       string
	Country         string
	Description     string
	CountryCode     string
	DefaultCurrency int64
	DateCreated     time.Time
	DateModified    time.Time
	CreatedBy       int
	ModifiedBy      int
}

type Branches struct {
	BranchId     string
	Branch       string
	Country      int64
	Location     string
	PhoneNumber  string
	Active       int
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    string
	ModifiedBy   string
}

type Shops struct {
	ShopId              string `orm:"auto"`
	ShopName            string
	ShopDescription     string `orm:"size(255)"`
	ShopAssistantName   string `orm:"size(100)"`
	ShopAssistantNumber string `orm:"size(100)"`
	PhoneNumber         string
	Email               string
	Image               string    `orm:"size(100);omitempty"`
	DateCreated         time.Time `orm:"type(datetime)"`
	DateModified        time.Time `orm:"type(datetime)"`
	CreatedBy           string
	ModifiedBy          string
	Active              int
}

type UserExtraDetails struct {
	UserDetailsId string
	Branch        *Branches
	Shop          *Shops
	Nickname      string
	DateCreated   time.Time
	DateModified  time.Time
	CreatedBy     string
	ModifiedBy    string
	Active        int
}

type Actions struct {
	ActionId     string
	Action       string
	Description  string
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    int
	ModifiedBy   int
	Active       int
}

type Permissions struct {
	PermissionId          string
	Permission            string
	PermissionCode        string
	PermissionDescription string
	DateCreated           time.Time
	DateModified          time.Time
	CreatedBy             string
	ModifiedBy            string
	Active                int
}

type Role_permissions struct {
	RolePermissionId string
	Role             *Roles
	Permission       *Permissions
	Action           *Actions
	DateCreated      time.Time
	DateModified     time.Time
	CreatedBy        string
	ModifiedBy       string
	Active           int
}

type Roles struct {
	RoleId          string
	Role            string
	Description     string
	DateCreated     time.Time
	DateModified    time.Time
	CreatedBy       string
	ModifiedBy      string
	Active          int
	RolePermissions []*Role_permissions
}

type Users struct {
	UserId        string
	UserType      int
	UserDetails   *UserExtraDetails
	ImagePath     string
	FullName      string
	Username      string
	Password      string
	Email         string
	PhoneNumber   string
	Gender        string
	Dob           time.Time
	Address       string
	IdType        string
	IdNumber      string
	MaritalStatus string
	Role          *Roles
	Active        int
	IsVerified    bool
	DateCreated   time.Time
	DateModified  time.Time
	CreatedBy     string
	ModifiedBy    string
}

type UserResp struct {
	UserId        int64
	ImagePath     string
	UserType      int
	FullName      string
	Username      string
	Password      string
	Email         string
	PhoneNumber   string
	Gender        string
	Dob           time.Time
	Address       string
	IdType        string
	IdNumber      string
	MaritalStatus string
	Active        int
	Role          *Roles
	IsVerified    bool
	DateCreated   time.Time
	DateModified  time.Time
	CreatedBy     string
	ModifiedBy    string
	Branch        *Branches
}

type UserResponseDTO struct {
	StatusCode int
	Result     *Users
	StatusDesc string
}

type UserTokenResponseDTO struct {
	IsValid bool
	User    *AuthenticatedUser
}

type RoleApiResponseDTO struct {
	StatusCode int
	Role       *Roles
	StatusDesc string
}

type UserPermission struct {
	PermissionCode string
	ActionCode     string
}

type Identification_types struct {
	IdentificationTypeId string
	Name                 string
	Code                 string
	DateCreated          time.Time
	DateModified         time.Time
	CreatedBy            string
	ModifiedBy           string
	Active               int
}

type Customer_categories struct {
	CustomerCategoryId string
	Category           string
	Description        string
	DateCreated        time.Time
	DateModified       time.Time
	CreatedBy          string
	ModifiedBy         string
	Active             int
}

type Customer_emergency_contacts struct {
	CustomerEmergencyContactId string
	Name                       string
	Contact                    string
	Customer                   *Customers
	DateCreated                time.Time
	DateModified               time.Time
	CreatedBy                  string
	ModifiedBy                 string
}

type Customer_guarantors struct {
	CustomerGuarantorId string
	Name                string
	Contact             string
	Customer            *Customers
	DateCreated         time.Time
	DateModified        time.Time
	CreatedBy           string
	ModifiedBy          string
}

type Customers struct {
	CustomerId           string
	CustomerNumber       string
	FullName             string
	ImagePath            string
	Email                string
	PhoneNumber          string
	Gender               string
	Location             string
	IdentificationType   *Identification_types
	IdentificationNumber string
	Branch               *Branches
	Shop                 *Shops
	CustomerCategory     *Customer_categories
	Nickname             string
	Dob                  time.Time
	DateCreated          time.Time
	DateModified         time.Time
	CreatedBy            string
	ModifiedBy           string
	Active               int
	LastTxnDate          time.Time
	EmergencyContacts    []*Customer_emergency_contacts
	Guarantors           []*Customer_guarantors
}

type CustomerTokenResponseDTO struct {
	IsValid  bool
	Customer *AuthenticatedCustomer
}

type CustomerResponseDTO struct {
	StatusCode int
	Result     *Customers
	StatusDesc string
}

type LoginResponseDTO struct {
	StatusCode   int
	AccessToken  string
	RefreshToken string
	User         *Users
	StatusDesc   string
}

type AuthenticatedUser struct {
	UserID      string
	Username    string
	RoleID      string
	RoleName    string
	ExpiryTime  int64
	BranchID    string
	Permissions []UserPermission
	Shop        string
}

type AuthenticatedCustomer struct {
	CustomerId       string
	Username         string
	Number           string
	CustomerCategory string
	ExpiryTime       int64
	BranchID         string
	Shop             string
}
