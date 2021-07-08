package sims

const (
	IntervalFlag      = "interval"
	EventGatewayFlag  = "events-gateway"
	SessionRangeFlag  = "session-range"
	SessionPrefixFlag = "session-prefix"
	ClientIdFlag      = "client-id"
	ClientSecretFlag  = "client-secret"
	ConsentLevelsFlag = "consent-levels"
	QuietFlag         = "quiet"
)

var BillingId string

func SetBillingId(billingId string) {
	BillingId = billingId

}
