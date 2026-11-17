package smsir

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSendSMSNilClient(t *testing.T) {
	var client *Client
	_, err := client.SendSMS(context.Background(), "09121234567", "hi")
	if err != ErrNotConfigured {
		t.Fatalf("err=%v", err)
	}
}

func TestWithTimeoutKeepsCustomClient(t *testing.T) {
	custom := &http.Client{Timeout: 5 * time.Second}
	after := New("key", 0, WithHTTPClient(custom), WithTimeout(time.Second))
	if after.httpClient != custom {
		t.Fatal("timeout replaced the custom client")
	}
	before := New("key", 0, WithTimeout(time.Second), WithHTTPClient(custom))
	if before.httpClient != custom {
		t.Fatal("custom client did not replace the timeout client")
	}
	onlyTimeout := New("key", 0, WithTimeout(2*time.Second))
	if onlyTimeout.httpClient.Timeout != 2*time.Second {
		t.Fatalf("timeout=%s", onlyTimeout.httpClient.Timeout)
	}
}

func TestSendLikeToLike(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/send/likeToLike" {
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("User-Agent") != "app/1" {
			t.Fatalf("ua=%q", r.Header.Get("User-Agent"))
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": 1, "message": "موفق",
			"data": map[string]any{"packId": "pack-1", "messageIds": []int64{1, 0}, "cost": 2},
		})
	}))
	defer srv.Close()

	when := int64(1700000000)
	cli := New("key", 3000, WithBaseURL(srv.URL), WithHTTPClient(srv.Client()), WithUserAgent("app/1"))
	res, err := cli.SendLikeToLike(context.Background(), LikeToLikeRequest{
		MessageTexts: []string{"one", "two"},
		Mobiles:      []string{"+989121234567", "09120000002"},
		SendDateTime: &when,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.PackID != "pack-1" || len(res.MessageIDs) != 2 || res.MessageIDs[1] != 0 {
		t.Fatalf("res=%+v", res)
	}
	if got["lineNumber"] != float64(3000) {
		t.Fatalf("line=%v", got["lineNumber"])
	}
	if got["sendDateTime"] != float64(when) {
		t.Fatalf("when=%v", got["sendDateTime"])
	}
}

func TestSendLikeToLikeLengthMismatch(t *testing.T) {
	cli := New("key", 3000)
	_, err := cli.SendLikeToLike(context.Background(), LikeToLikeRequest{
		MessageTexts: []string{"only one"},
		Mobiles:      []string{"09120000001", "09120000002"},
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreditAndLines(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/credit":
			_ = json.NewEncoder(w).Encode(map[string]any{"status": 1, "message": "موفق", "data": 12.5})
		case "/v1/line":
			_ = json.NewEncoder(w).Encode(map[string]any{"status": 1, "message": "موفق", "data": []int64{10002155613464, 30004505000017}})
		default:
			t.Fatalf("path=%s", r.URL.Path)
		}
	}))
	defer srv.Close()

	cli := New("key", 0, WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
	credit, err := cli.Credit(context.Background())
	if err != nil || credit != 12.5 {
		t.Fatalf("credit=%v err=%v", credit, err)
	}
	lines, err := cli.Lines(context.Background())
	if err != nil || len(lines) != 2 || lines[1] != 30004505000017 {
		t.Fatalf("lines=%v err=%v", lines, err)
	}
}

func TestArchiveAndTodayReportQuery(t *testing.T) {
	var gotPath, gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		_ = json.NewEncoder(w).Encode(map[string]any{"status": 1, "message": "موفق", "data": []any{}})
	}))
	defer srv.Close()

	cli := New("key", 0, WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
	if _, err := cli.TodayReport(context.Background(), 0, 0); err != nil {
		t.Fatal(err)
	}
	if gotPath != "/v1/send/live" || gotQuery != "" {
		t.Fatalf("today %s ?%s", gotPath, gotQuery)
	}
	from := int64(1613465574)
	to := int64(1623805200)
	if _, err := cli.ArchiveReport(context.Background(), &from, &to, 25, 2); err != nil {
		t.Fatal(err)
	}
	if gotPath != "/v1/send/archive" ||
		!strings.Contains(gotQuery, "fromDate=1613465574") ||
		!strings.Contains(gotQuery, "toDate=1623805200") ||
		!strings.Contains(gotQuery, "pageSize=25") ||
		!strings.Contains(gotQuery, "pageNumber=2") {
		t.Fatalf("archive %s", gotQuery)
	}
}

func TestAPIErrorShapes(t *testing.T) {
	bare := &APIError{Status: 104, HTTPStatus: http.StatusBadRequest}
	if bare.Error() == "" {
		t.Fatal("empty error")
	}
	wrapped := fmt.Errorf("call: %w", &APIError{Status: 20, HTTPStatus: http.StatusOK, Message: "limited"})
	apiErr, ok := AsAPIError(wrapped)
	if !ok || apiErr.Status != StatusRateLimited {
		t.Fatalf("unwrap=%v ok=%v", apiErr, ok)
	}
	if _, ok := AsAPIError(fmt.Errorf("plain")); ok {
		t.Fatal("plain error unwrapped")
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"status": 20, "message": "limited", "data": nil})
	}))
	defer srv.Close()
	cli := New("key", 1, WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
	_, err := cli.SendSMS(context.Background(), "09121234567", "x")
	apiErr, ok = AsAPIError(err)
	if !ok || apiErr.Status != StatusRateLimited || apiErr.HTTPStatus != http.StatusOK {
		t.Fatalf("err=%v", err)
	}
}

func TestRemoveTemplateNotRemoved(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"status": 1, "message": "موفق", "data": false})
	}))
	defer srv.Close()
	cli := New("key", 0, WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
	if err := cli.RemoveTemplate(context.Background(), 123456); err == nil {
		t.Fatal("expected error")
	}
}
