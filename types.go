package smsir

import (
	"encoding/json"
	"errors"
	"fmt"
)

// API business status codes returned in the JSON envelope "status" field.
// The full table, including codes without a constant, is in docs/statusCode.md.
const (
	StatusError              = 0
	StatusOK                 = 1
	StatusInvalidAPIKey      = 10
	StatusInactiveAPIKey     = 11
	StatusIPRestricted       = 12
	StatusAccountDisabled    = 13
	StatusAccountSuspended   = 14
	StatusUpgradePlan        = 15
	StatusBadParameter       = 16
	StatusRateLimited        = 20
	StatusInvalidLine        = 101
	StatusInsufficientCredit = 102
	StatusEmptyText          = 103
	StatusInvalidMobile      = 104
	StatusTooManyMobiles     = 105
	StatusTemplateNotFound   = 113
	StatusParamTooLong       = 114
	StatusBlacklisted        = 115
	StatusMissingParameter   = 116
)

// Template review status (GET /v1/templates and Status filter).
const (
	TemplateStatusPending  byte = 1
	TemplateStatusApproved byte = 2
	TemplateStatusRejected byte = 3
)

// Template type codes.
const (
	TemplateTypeOTP         byte = 1
	TemplateTypeInfo        byte = 2
	TemplateTypeOrderStatus byte = 3
	TemplateTypeOther       byte = 100
)

// Delivery status codes from SMS.ir message reports.
const (
	DeliveryReachedPhone      = 1
	DeliveryNotReachedPhone   = 2
	DeliveryReachedTelecom    = 3
	DeliveryNotReachedTelecom = 4
	DeliveryReachedOperator   = 5
	DeliveryFailed            = 6
	DeliveryBlacklist         = 7
	DeliveryUnknown           = 8
)

type apiEnvelope struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// APIError is a non-success SMS.ir API response (business or HTTP failure).
type APIError struct {
	HTTPStatus int
	Status     int
	Message    string
}

func (e *APIError) Error() string {
	if e == nil {
		return "smsir: unknown error"
	}
	if e.Message != "" {
		return fmt.Sprintf("smsir: status=%d http=%d: %s", e.Status, e.HTTPStatus, e.Message)
	}

	return fmt.Sprintf("smsir: status=%d http=%d", e.Status, e.HTTPStatus)
}

// AsAPIError unwraps err as [*APIError] when possible.
func AsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr, true
	}

	return nil, false
}

// ErrNotConfigured is returned when the API key is empty.
var ErrNotConfigured = errors.New("smsir: not configured")

// BulkSendRequest is the body for POST /v1/send/bulk.
type BulkSendRequest struct {
	LineNumber   int64
	MessageText  string
	Mobiles      []string
	SendDateTime *int64 // unix UTC; nil = send now
}

// LikeToLikeRequest is the body for POST /v1/send/likeToLike.
type LikeToLikeRequest struct {
	LineNumber   int64
	MessageTexts []string
	Mobiles      []string
	SendDateTime *int64
}

// BulkSendResult is returned by bulk and like-to-like sends.
// MessageIDs[i] == 0 means the corresponding mobile is blacklisted.
type BulkSendResult struct {
	PackID     string  `json:"packId"`
	MessageIDs []int64 `json:"messageIds"`
	Cost       float64 `json:"cost"`
	HTTPStatus int     `json:"-"`
}

// VerifyParameter is a template placeholder replacement (#NAME# → Value).
type VerifyParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// VerifyRequest is the body for POST /v1/send/verify.
type VerifyRequest struct {
	Mobile     string
	TemplateID int
	Parameters []VerifyParameter
}

// VerifyResult is returned by verify send.
type VerifyResult struct {
	MessageID  int64   `json:"messageId"`
	Cost       float64 `json:"cost"`
	HTTPStatus int     `json:"-"`
}

// CancelScheduledResult is returned by DELETE /v1/send/scheduled/{packId}.
type CancelScheduledResult struct {
	ReturnedCreditAmount float64 `json:"returnedCreditAmount"`
	SMSCount             int     `json:"smsCount"`
	HTTPStatus           int     `json:"-"`
}

// LivePack is one row from GET /v1/send/pack (today's send packs).
type LivePack struct {
	PackID           string `json:"packId"`
	RecipientCount   int    `json:"recipientCount"`
	CreationDateTime int64  `json:"creationDateTime"`
}

// MessageReport is one sent message from GET /v1/send/{messageId}.
// A null deliveryState decodes as 0.
type MessageReport struct {
	MessageID        int64      `json:"messageId"`
	Mobile           FlexString `json:"mobile"`
	MessageText      string     `json:"messageText"`
	SendDateTime     int64      `json:"sendDateTime"`
	LineNumber       int64      `json:"lineNumber"`
	Cost             float64    `json:"cost"`
	DeliveryState    int        `json:"deliveryState"`
	DeliveryDateTime *int64     `json:"deliveryDateTime"`
}

func (m *MessageReport) UnmarshalJSON(data []byte) error {
	var aux struct {
		MessageID        int64      `json:"messageId"`
		Mobile           FlexString `json:"mobile"`
		MessageText      string     `json:"messageText"`
		SendDateTime     int64      `json:"sendDateTime"`
		LineNumber       int64      `json:"lineNumber"`
		Cost             float64    `json:"cost"`
		DeliveryState    *int       `json:"deliveryState"`
		DeliveryDateTime *int64     `json:"deliveryDateTime"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*m = MessageReport{
		MessageID:        aux.MessageID,
		Mobile:           aux.Mobile,
		MessageText:      aux.MessageText,
		SendDateTime:     aux.SendDateTime,
		LineNumber:       aux.LineNumber,
		Cost:             aux.Cost,
		DeliveryDateTime: aux.DeliveryDateTime,
	}
	if aux.DeliveryState != nil {
		m.DeliveryState = *aux.DeliveryState
	}

	return nil
}

// ReceivedMessage is an inbound SMS row.
type ReceivedMessage struct {
	Mobile           FlexString `json:"mobile"`
	MessageText      string     `json:"messageText"`
	ReceivedDateTime int64      `json:"receivedDateTime"`
	LineNumber       int64      `json:"lineNumber"`
}

// TemplateParameter describes a #NAME# placeholder in a template.
// Name is written without the surrounding # and is at most 40 characters.
type TemplateParameter struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TemplateRequest is the body for AddTemplate and ModifyTemplate.
// Type 0 omits the field: create then uses TemplateTypeOther (100), and update keeps the current type.
// A nil or empty Parameters list is omitted. On update, that leaves the current descriptions in place.
type TemplateRequest struct {
	Title        string
	TemplateText string
	Type         byte
	Parameters   []TemplateParameter
}

// Template is a verify/fast-send template from GET /v1/templates.
type Template struct {
	ID               int                 `json:"id"`
	Title            string              `json:"title"`
	TemplateText     string              `json:"templateText"`
	Type             byte                `json:"type"`
	Parameters       []TemplateParameter `json:"parameters"`
	Status           byte                `json:"status"`
	RejectionReason  *string             `json:"rejectionReason"`
	CreationDateTime int64               `json:"creationDateTime"`
}
