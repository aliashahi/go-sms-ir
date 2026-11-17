package smsir

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL   = "https://api.sms.ir"
	defaultUserAgent = "go-sms-ir/1.0"
	maxBodyBytes     = 1 << 20
)

// templatePlaceholder matches a #NAME# token in template text.
var templatePlaceholder = regexp.MustCompile(`#[^#\s]+#`)

// Client talks to api.sms.ir using an X-API-KEY.
type Client struct {
	apiKey     string
	lineNumber int64
	baseURL    string
	userAgent  string
	httpClient *http.Client
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides the API host (default https://api.sms.ir).
func WithBaseURL(base string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(strings.TrimSpace(base), "/")
	}
}

// WithHTTPClient sets a custom HTTP client (timeouts, proxies, transport).
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		if hc != nil {
			c.httpClient = hc
		}
	}
}

// WithUserAgent overrides the User-Agent header (default go-sms-ir/1.0).
func WithUserAgent(ua string) Option {
	return func(c *Client) {
		if strings.TrimSpace(ua) != "" {
			c.userAgent = strings.TrimSpace(ua)
		}
	}
}

// WithTimeout sets the default HTTP client timeout. Ignored if WithHTTPClient was used.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		if d > 0 {
			c.httpClient = &http.Client{Timeout: d}
		}
	}
}

// New creates a Client. lineNumber is the default sender for bulk / like-to-like.
// Verify/OTP sends do not require a line number.
func New(apiKey string, lineNumber int64, opts ...Option) *Client {
	c := &Client{
		apiKey:     strings.TrimSpace(apiKey),
		lineNumber: lineNumber,
		baseURL:    defaultBaseURL,
		userAgent:  defaultUserAgent,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Configured reports whether an API key is set.
func (c *Client) Configured() bool {
	return c != nil && c.apiKey != ""
}

// LineNumber returns the default sender line configured on the client.
func (c *Client) LineNumber() int64 {
	if c == nil {
		return 0
	}

	return c.lineNumber
}

// SendSMS sends one free-text message via the bulk endpoint.
func (c *Client) SendSMS(ctx context.Context, mobile, message string) (*BulkSendResult, error) {
	return c.SendBulk(ctx, BulkSendRequest{
		LineNumber:  c.lineNumber,
		MessageText: message,
		Mobiles:     []string{mobile},
	})
}

// SendBulk sends one message to many mobiles (max 100).
// POST /v1/send/bulk
func (c *Client) SendBulk(ctx context.Context, req BulkSendRequest) (*BulkSendResult, error) {
	if !c.Configured() {
		return nil, ErrNotConfigured
	}
	if req.LineNumber == 0 {
		req.LineNumber = c.lineNumber
	}
	if req.LineNumber == 0 {
		return nil, fmt.Errorf("smsir: line number is required")
	}
	if strings.TrimSpace(req.MessageText) == "" {
		return nil, fmt.Errorf("smsir: message text is required")
	}
	if len(req.Mobiles) == 0 {
		return nil, fmt.Errorf("smsir: mobiles are required")
	}
	if len(req.Mobiles) > 100 {
		return nil, fmt.Errorf("smsir: at most 100 mobiles are allowed")
	}

	mobiles := make([]string, 0, len(req.Mobiles))
	for _, m := range req.Mobiles {
		mobiles = append(mobiles, FormatMobile(m))
	}
	body := map[string]any{
		"lineNumber":  req.LineNumber,
		"messageText": req.MessageText,
		"mobiles":     mobiles,
	}
	if req.SendDateTime != nil {
		body["sendDateTime"] = *req.SendDateTime
	}

	var out BulkSendResult
	httpStatus, err := c.doJSON(ctx, http.MethodPost, "/v1/send/bulk", nil, body, &out)
	out.HTTPStatus = httpStatus
	if err != nil {
		return &out, err
	}

	return &out, nil
}

// SendLikeToLike sends paired messages to mobiles (lengths must match, max 100).
// POST /v1/send/likeToLike
func (c *Client) SendLikeToLike(ctx context.Context, req LikeToLikeRequest) (*BulkSendResult, error) {
	if !c.Configured() {
		return nil, ErrNotConfigured
	}
	if req.LineNumber == 0 {
		req.LineNumber = c.lineNumber
	}
	if req.LineNumber == 0 {
		return nil, fmt.Errorf("smsir: line number is required")
	}
	if len(req.Mobiles) == 0 || len(req.MessageTexts) == 0 {
		return nil, fmt.Errorf("smsir: mobiles and messageTexts are required")
	}
	if len(req.Mobiles) != len(req.MessageTexts) {
		return nil, fmt.Errorf("smsir: mobiles and messageTexts length must match")
	}
	if len(req.Mobiles) > 100 {
		return nil, fmt.Errorf("smsir: at most 100 mobiles are allowed")
	}

	mobiles := make([]string, 0, len(req.Mobiles))
	for _, m := range req.Mobiles {
		mobiles = append(mobiles, FormatMobile(m))
	}
	body := map[string]any{
		"lineNumber":   req.LineNumber,
		"messageTexts": req.MessageTexts,
		"mobiles":      mobiles,
	}
	if req.SendDateTime != nil {
		body["sendDateTime"] = *req.SendDateTime
	}

	var out BulkSendResult
	httpStatus, err := c.doJSON(ctx, http.MethodPost, "/v1/send/likeToLike", nil, body, &out)
	out.HTTPStatus = httpStatus
	if err != nil {
		return &out, err
	}

	return &out, nil
}

// SendVerify sends a high-priority templated message (OTP / verify).
// Service lines are used; messages can reach numbers that blocked ads.
// POST /v1/send/verify
func (c *Client) SendVerify(ctx context.Context, req VerifyRequest) (*VerifyResult, error) {
	if !c.Configured() {
		return nil, ErrNotConfigured
	}
	if strings.TrimSpace(req.Mobile) == "" {
		return nil, fmt.Errorf("smsir: mobile is required")
	}
	if req.TemplateID == 0 {
		return nil, fmt.Errorf("smsir: templateId is required")
	}
	if len(req.Parameters) == 0 {
		return nil, fmt.Errorf("smsir: parameters are required")
	}

	body := map[string]any{
		"mobile":     FormatMobile(req.Mobile),
		"templateId": req.TemplateID,
		"parameters": req.Parameters,
	}

	var out VerifyResult
	httpStatus, err := c.doJSON(ctx, http.MethodPost, "/v1/send/verify", nil, body, &out)
	out.HTTPStatus = httpStatus
	if err != nil {
		return &out, err
	}

	return &out, nil
}

// CancelScheduled deletes a scheduled pack.
// DELETE /v1/send/scheduled/{packId}
func (c *Client) CancelScheduled(ctx context.Context, packID string) (*CancelScheduledResult, error) {
	if !c.Configured() {
		return nil, ErrNotConfigured
	}
	packID = strings.TrimSpace(packID)
	if packID == "" {
		return nil, fmt.Errorf("smsir: packId is required")
	}

	var out CancelScheduledResult
	path := "/v1/send/scheduled/" + url.PathEscape(packID)
	httpStatus, err := c.doJSON(ctx, http.MethodDelete, path, nil, nil, &out)
	out.HTTPStatus = httpStatus
	if err != nil {
		return &out, err
	}

	return &out, nil
}

// MessageReport returns delivery info for one message id.
// GET /v1/send/{messageId}
func (c *Client) MessageReport(ctx context.Context, messageID int64) (*MessageReport, error) {
	var out MessageReport
	path := "/v1/send/" + strconv.FormatInt(messageID, 10)
	if _, err := c.doJSON(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// TodayPacks lists today's send packs.
// GET /v1/send/pack
// pageSize max is 100 (API default 100). pageNumber defaults to 1 when omitted.
func (c *Client) TodayPacks(ctx context.Context, pageSize, pageNumber int) ([]LivePack, error) {
	q := pageQuery(pageSize, pageNumber)
	var out []LivePack
	if _, err := c.doJSON(ctx, http.MethodGet, "/v1/send/pack", q, nil, &out); err != nil {
		return nil, err
	}

	return out, nil
}

// PackReport returns messages in a pack.
// GET /v1/send/pack/{packId}
func (c *Client) PackReport(ctx context.Context, packID string) ([]MessageReport, error) {
	var out []MessageReport
	path := "/v1/send/pack/" + url.PathEscape(strings.TrimSpace(packID))
	if _, err := c.doJSON(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}

	return out, nil
}

// TodayReport lists today's sent messages.
// GET /v1/send/live
func (c *Client) TodayReport(ctx context.Context, pageSize, pageNumber int) ([]MessageReport, error) {
	q := pageQuery(pageSize, pageNumber)
	var out []MessageReport
	if _, err := c.doJSON(ctx, http.MethodGet, "/v1/send/live", q, nil, &out); err != nil {
		return nil, err
	}

	return out, nil
}

// ArchiveReport lists archived sent messages.
// GET /v1/send/archive
func (c *Client) ArchiveReport(ctx context.Context, fromDate, toDate *int64, pageSize, pageNumber int) ([]MessageReport, error) {
	q := pageQuery(pageSize, pageNumber)
	if fromDate != nil {
		q.Set("fromDate", strconv.FormatInt(*fromDate, 10))
	}
	if toDate != nil {
		q.Set("toDate", strconv.FormatInt(*toDate, 10))
	}
	var out []MessageReport
	if _, err := c.doJSON(ctx, http.MethodGet, "/v1/send/archive", q, nil, &out); err != nil {
		return nil, err
	}

	return out, nil
}

// LatestReceived returns the latest inbound messages.
// GET /v1/receive/latest
func (c *Client) LatestReceived(ctx context.Context, count int) ([]ReceivedMessage, error) {
	q := url.Values{}
	if count > 0 {
		q.Set("count", strconv.Itoa(count))
	}
	var out []ReceivedMessage
	if _, err := c.doJSON(ctx, http.MethodGet, "/v1/receive/latest", q, nil, &out); err != nil {
		return nil, err
	}

	return out, nil
}

// TodayReceived lists today's inbound messages.
// GET /v1/receive/live
func (c *Client) TodayReceived(ctx context.Context, pageSize, pageNumber int) ([]ReceivedMessage, error) {
	q := pageQuery(pageSize, pageNumber)
	var out []ReceivedMessage
	if _, err := c.doJSON(ctx, http.MethodGet, "/v1/receive/live", q, nil, &out); err != nil {
		return nil, err
	}

	return out, nil
}

// ArchiveReceived lists archived inbound messages.
// GET /v1/receive/archive
func (c *Client) ArchiveReceived(ctx context.Context, fromDate, toDate *int64, pageSize, pageNumber int) ([]ReceivedMessage, error) {
	q := pageQuery(pageSize, pageNumber)
	if fromDate != nil {
		q.Set("fromDate", strconv.FormatInt(*fromDate, 10))
	}
	if toDate != nil {
		q.Set("toDate", strconv.FormatInt(*toDate, 10))
	}
	var out []ReceivedMessage
	if _, err := c.doJSON(ctx, http.MethodGet, "/v1/receive/archive", q, nil, &out); err != nil {
		return nil, err
	}

	return out, nil
}

// Credit returns remaining account credit.
// GET /v1/credit
func (c *Client) Credit(ctx context.Context) (float64, error) {
	var out float64
	if _, err := c.doJSON(ctx, http.MethodGet, "/v1/credit", nil, nil, &out); err != nil {
		return 0, err
	}

	return out, nil
}

// Lines returns available sender line numbers.
// GET /v1/line
func (c *Client) Lines(ctx context.Context) ([]int64, error) {
	var out []int64
	if _, err := c.doJSON(ctx, http.MethodGet, "/v1/line", nil, nil, &out); err != nil {
		return nil, err
	}

	return out, nil
}

// Templates lists verify templates (panel fast-send).
// Only templates with [TemplateStatusApproved] can be used with [Client.SendVerify].
// GET /v1/templates
// pageNumber defaults to 1 and pageSize to 10 when 0. pageSize max is 100.
// A nil status lists every review state.
func (c *Client) Templates(ctx context.Context, pageNumber, pageSize int, status *byte) ([]Template, error) {
	q := pageQuery(pageSize, pageNumber)
	if status != nil {
		q.Set("status", strconv.Itoa(int(*status)))
	}
	var out []Template
	if _, err := c.doJSON(ctx, http.MethodGet, "/v1/templates", q, nil, &out); err != nil {
		return nil, err
	}

	return out, nil
}

// Template returns one verify template.
// GET /v1/templates/{templateId}
func (c *Client) Template(ctx context.Context, templateID int) (*Template, error) {
	if err := requireTemplateID(templateID); err != nil {
		return nil, err
	}
	var out Template
	if _, err := c.doJSON(ctx, http.MethodGet, templatePath(templateID), nil, nil, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// AddTemplate registers a verify template. It stays pending until SMS.ir approves it.
// POST /v1/templates
func (c *Client) AddTemplate(ctx context.Context, req TemplateRequest) (int, error) {
	if !c.Configured() {
		return 0, ErrNotConfigured
	}
	body, err := templateBody(req)
	if err != nil {
		return 0, err
	}
	var id int
	if _, err := c.doJSON(ctx, http.MethodPost, "/v1/templates", nil, body, &id); err != nil {
		return 0, err
	}

	return id, nil
}

// ModifyTemplate updates a template that is not pending review.
// The template returns to pending until it is approved again.
// PUT /v1/templates/{templateId}
func (c *Client) ModifyTemplate(ctx context.Context, templateID int, req TemplateRequest) (int, error) {
	if !c.Configured() {
		return 0, ErrNotConfigured
	}
	if err := requireTemplateID(templateID); err != nil {
		return 0, err
	}
	body, err := templateBody(req)
	if err != nil {
		return 0, err
	}
	var id int
	if _, err := c.doJSON(ctx, http.MethodPut, templatePath(templateID), nil, body, &id); err != nil {
		return 0, err
	}

	return id, nil
}

// RemoveTemplate deletes a verify template. Sent-message reports are kept.
// DELETE /v1/templates/{templateId}
func (c *Client) RemoveTemplate(ctx context.Context, templateID int) error {
	if !c.Configured() {
		return ErrNotConfigured
	}
	if err := requireTemplateID(templateID); err != nil {
		return err
	}
	var removed bool
	if _, err := c.doJSON(ctx, http.MethodDelete, templatePath(templateID), nil, nil, &removed); err != nil {
		return err
	}
	if !removed {
		return fmt.Errorf("smsir: template was not removed")
	}

	return nil
}

func requireTemplateID(templateID int) error {
	if templateID <= 0 {
		return fmt.Errorf("smsir: templateId is required")
	}

	return nil
}

func templatePath(templateID int) string {
	return "/v1/templates/" + strconv.Itoa(templateID)
}

func templateBody(req TemplateRequest) (map[string]any, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, fmt.Errorf("smsir: title is required")
	}
	text := req.TemplateText
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("smsir: templateText is required")
	}
	if !templatePlaceholder.MatchString(text) {
		return nil, fmt.Errorf("smsir: templateText must include a #PARAMETER# placeholder")
	}
	body := map[string]any{
		"title":        title,
		"templateText": text,
	}
	if req.Type != 0 {
		body["type"] = req.Type
	}
	if len(req.Parameters) > 0 {
		params := make([]TemplateParameter, len(req.Parameters))
		for i, p := range req.Parameters {
			name := strings.Trim(strings.TrimSpace(p.Name), "#")
			if name == "" {
				return nil, fmt.Errorf("smsir: parameter name is required")
			}
			params[i] = TemplateParameter{Name: name, Description: p.Description}
		}
		body["parameters"] = params
	}

	return body, nil
}

func pageQuery(pageSize, pageNumber int) url.Values {
	q := url.Values{}
	if pageSize > 0 {
		q.Set("pageSize", strconv.Itoa(pageSize))
	}
	if pageNumber > 0 {
		q.Set("pageNumber", strconv.Itoa(pageNumber))
	}

	return q
}

func (c *Client) doJSON(ctx context.Context, method, path string, query url.Values, reqBody any, dataOut any) (httpStatus int, err error) {
	if !c.Configured() {
		return 0, ErrNotConfigured
	}

	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if reqBody != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			return 0, fmt.Errorf("smsir: marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return 0, fmt.Errorf("smsir: build request: %w", err)
	}
	req.Header.Set("X-API-KEY", c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("smsir: http: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	if err != nil {
		return resp.StatusCode, fmt.Errorf("smsir: read body: %w", err)
	}

	var env apiEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return resp.StatusCode, fmt.Errorf("smsir: decode response (http %d): %w; body=%s", resp.StatusCode, err, truncate(string(raw), 200))
	}

	if resp.StatusCode == http.StatusUnauthorized || env.Status == StatusInvalidAPIKey || env.Status == StatusInactiveAPIKey {
		return resp.StatusCode, &APIError{HTTPStatus: resp.StatusCode, Status: env.Status, Message: env.Message}
	}
	if resp.StatusCode == http.StatusTooManyRequests || env.Status == StatusRateLimited {
		return resp.StatusCode, &APIError{HTTPStatus: resp.StatusCode, Status: env.Status, Message: env.Message}
	}
	if resp.StatusCode >= 400 || env.Status != StatusOK {
		return resp.StatusCode, &APIError{HTTPStatus: resp.StatusCode, Status: env.Status, Message: env.Message}
	}

	if dataOut == nil || len(env.Data) == 0 || string(env.Data) == "null" {
		return resp.StatusCode, nil
	}
	if err := json.Unmarshal(env.Data, dataOut); err != nil {
		return resp.StatusCode, fmt.Errorf("smsir: decode data: %w", err)
	}

	return resp.StatusCode, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}

	return s[:n] + "..."
}
