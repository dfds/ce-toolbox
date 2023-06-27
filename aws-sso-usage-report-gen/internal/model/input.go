package model

import "time"

type InputData []AccessEntry

type AccessEntry struct {
	ID                                   string        `json:"id"`
	CreatedDateTime                      time.Time     `json:"createdDateTime"`
	UserDisplayName                      string        `json:"userDisplayName"`
	UserPrincipalName                    string        `json:"userPrincipalName"`
	UserID                               string        `json:"userId"`
	AppID                                string        `json:"appId"`
	AppDisplayName                       string        `json:"appDisplayName"`
	IPAddress                            string        `json:"ipAddress"`
	IPAddressFromResourceProvider        interface{}   `json:"ipAddressFromResourceProvider"`
	ClientAppUsed                        string        `json:"clientAppUsed"`
	UserAgent                            string        `json:"userAgent"`
	CorrelationID                        string        `json:"correlationId"`
	ConditionalAccessStatus              string        `json:"conditionalAccessStatus"`
	OriginalRequestID                    string        `json:"originalRequestId"`
	IsInteractive                        bool          `json:"isInteractive"`
	TokenIssuerName                      string        `json:"tokenIssuerName"`
	TokenIssuerType                      string        `json:"tokenIssuerType"`
	ClientCredentialType                 string        `json:"clientCredentialType"`
	ProcessingTimeInMilliseconds         int           `json:"processingTimeInMilliseconds"`
	RiskDetail                           string        `json:"riskDetail"`
	RiskLevelAggregated                  string        `json:"riskLevelAggregated"`
	RiskLevelDuringSignIn                string        `json:"riskLevelDuringSignIn"`
	RiskState                            string        `json:"riskState"`
	RiskEventTypesV2                     []interface{} `json:"riskEventTypes_v2"`
	ResourceDisplayName                  string        `json:"resourceDisplayName"`
	ResourceID                           string        `json:"resourceId"`
	ResourceTenantID                     string        `json:"resourceTenantId"`
	HomeTenantID                         string        `json:"homeTenantId"`
	HomeTenantName                       string        `json:"homeTenantName"`
	AuthenticationMethodsUsed            []interface{} `json:"authenticationMethodsUsed"`
	AuthenticationRequirement            string        `json:"authenticationRequirement"`
	SignInIdentifier                     string        `json:"signInIdentifier"`
	SignInIdentifierType                 interface{}   `json:"signInIdentifierType"`
	ServicePrincipalName                 interface{}   `json:"servicePrincipalName"`
	SignInEventTypes                     []string      `json:"signInEventTypes"`
	ServicePrincipalID                   string        `json:"servicePrincipalId"`
	FederatedCredentialID                interface{}   `json:"federatedCredentialId"`
	UserType                             string        `json:"userType"`
	FlaggedForReview                     bool          `json:"flaggedForReview"`
	IsTenantRestricted                   bool          `json:"isTenantRestricted"`
	AutonomousSystemNumber               int           `json:"autonomousSystemNumber"`
	CrossTenantAccessType                string        `json:"crossTenantAccessType"`
	ServicePrincipalCredentialKeyID      interface{}   `json:"servicePrincipalCredentialKeyId"`
	ServicePrincipalCredentialThumbprint string        `json:"servicePrincipalCredentialThumbprint"`
	UniqueTokenIdentifier                string        `json:"uniqueTokenIdentifier"`
	IncomingTokenType                    string        `json:"incomingTokenType"`
	AuthenticationProtocol               string        `json:"authenticationProtocol"`
	ResourceServicePrincipalID           string        `json:"resourceServicePrincipalId"`
	AuthenticationAppDeviceDetails       interface{}   `json:"authenticationAppDeviceDetails"`
	Status                               struct {
		ErrorCode         int    `json:"errorCode"`
		FailureReason     string `json:"failureReason"`
		AdditionalDetails string `json:"additionalDetails"`
	} `json:"status"`
	DeviceDetail struct {
		DeviceID        string `json:"deviceId"`
		DisplayName     string `json:"displayName"`
		OperatingSystem string `json:"operatingSystem"`
		Browser         string `json:"browser"`
		IsCompliant     bool   `json:"isCompliant"`
		IsManaged       bool   `json:"isManaged"`
		TrustType       string `json:"trustType"`
	} `json:"deviceDetail"`
	Location struct {
		City            string `json:"city"`
		State           string `json:"state"`
		CountryOrRegion string `json:"countryOrRegion"`
		GeoCoordinates  struct {
			Altitude  interface{} `json:"altitude"`
			Latitude  float64     `json:"latitude"`
			Longitude float64     `json:"longitude"`
		} `json:"geoCoordinates"`
	} `json:"location"`
	MfaDetail struct {
		AuthMethod interface{} `json:"authMethod"`
		AuthDetail interface{} `json:"authDetail"`
	} `json:"mfaDetail"`
	AppliedConditionalAccessPolicies     []interface{} `json:"appliedConditionalAccessPolicies"`
	AuthenticationContextClassReferences []struct {
		ID     string `json:"id"`
		Detail string `json:"detail"`
	} `json:"authenticationContextClassReferences"`
	AuthenticationProcessingDetails []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"authenticationProcessingDetails"`
	NetworkLocationDetails []struct {
		NetworkType  string   `json:"networkType"`
		NetworkNames []string `json:"networkNames"`
	} `json:"networkLocationDetails"`
	AuthenticationDetails []struct {
		AuthenticationStepDateTime     time.Time   `json:"authenticationStepDateTime"`
		AuthenticationMethod           string      `json:"authenticationMethod"`
		AuthenticationMethodDetail     interface{} `json:"authenticationMethodDetail"`
		Succeeded                      bool        `json:"succeeded"`
		AuthenticationStepResultDetail string      `json:"authenticationStepResultDetail"`
		AuthenticationStepRequirement  string      `json:"authenticationStepRequirement"`
	} `json:"authenticationDetails"`
	AuthenticationRequirementPolicies []struct {
		RequirementProvider string `json:"requirementProvider"`
		Detail              string `json:"detail"`
	} `json:"authenticationRequirementPolicies"`
	SessionLifetimePolicies []struct {
		ExpirationRequirement string `json:"expirationRequirement"`
		Detail                string `json:"detail"`
	} `json:"sessionLifetimePolicies"`
	PrivateLinkDetails struct {
		PolicyID       string `json:"policyId"`
		PolicyName     string `json:"policyName"`
		ResourceID     string `json:"resourceId"`
		PolicyTenantID string `json:"policyTenantId"`
	} `json:"privateLinkDetails"`
	AppliedEventListeners                    []interface{} `json:"appliedEventListeners"`
	AuthenticationAppPolicyEvaluationDetails []interface{} `json:"authenticationAppPolicyEvaluationDetails"`
}
