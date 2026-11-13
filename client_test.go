package smsir

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendBulkSuccess(t *testing.T) {
	var gotBody map[string]any
	var gotKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("X-API-KEY")
		if r.URL.Path != "/v1/send/bulk" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatal(err)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  1,
			"message": "موفق",
			"data": map[string]any{
				"packId":     "pack-1",
				"messageIds": []int64{111},
				"cost":       1.0,
			},
		})
	}))
	defer srv.Close()

	cli := New("test-key", 3000, WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
	res, err := cli.SendSMS(context.Background(), "+989121234567", "hello")
	if err != nil {
		t.Fatal(err)
	}
	if gotKey != "test-key" {
		t.Fatalf("api key=%q", gotKey)
	}
	if gotBody["messageText"] != "hello" {
		t.Fatalf("body=%v", gotBody)
	}
	mobiles, _ := gotBody["mobiles"].([]any)
	if len(mobiles) != 1 || mobiles[0] != "09121234567" {
		t.Fatalf("mobiles=%v", gotBody["mobiles"])
	}
	if res.PackID != "pack-1" || len(res.MessageIDs) != 1 || res.MessageIDs[0] != 111 {
		t.Fatalf("res=%+v", res)
	}
}

func TestSendVerifySuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/send/verify" || r.Method != http.MethodPost {
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["mobile"] != "09121234567" {
			t.Fatalf("mobile=%v", body["mobile"])
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  1,
			"message": "موفق",
			"data":    map[string]any{"messageId": 89545112, "cost": 1.0},
		})
	}))
	defer srv.Close()

	cli := New("key", 0, WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
	res, err := cli.SendVerify(context.Background(), VerifyRequest{
		Mobile:     "9121234567",
		TemplateID: 123456,
		Parameters: []VerifyParameter{{Name: "CODE", Value: "1234"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.MessageID != 89545112 {
		t.Fatalf("res=%+v", res)
	}
	if res.HTTPStatus != http.StatusOK {
		t.Fatalf("http=%d", res.HTTPStatus)
	}
}

func TestAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  102,
			"message": "اعتبار کافی نمیباشد",
			"data":    nil,
		})
	}))
	defer srv.Close()

	cli := New("key", 3000, WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
	_, err := cli.SendSMS(context.Background(), "09121234567", "x")
	apiErr, ok := AsAPIError(err)
	if !ok {
		t.Fatalf("err type=%T %v", err, err)
	}
	if apiErr.Status != 102 || apiErr.HTTPStatus != http.StatusBadRequest {
		t.Fatalf("apiErr=%+v", apiErr)
	}
}

func TestTemplates(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/templates" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		if r.URL.Query().Get("status") != "2" {
			t.Fatalf("query=%s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  1,
			"message": "موفق",
			"data": []map[string]any{{
				"id":               123456,
				"title":            "کد تایید ورود",
				"templateText":     "کد تایید شما #CODE# می‌باشد.",
				"type":             1,
				"parameters":       []map[string]any{{"name": "CODE", "description": "کد"}},
				"status":           2,
				"rejectionReason":  nil,
				"creationDateTime": 1628683626,
			}},
		})
	}))
	defer srv.Close()

	cli := New("key", 0, WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
	st := TemplateStatusApproved
	list, err := cli.Templates(context.Background(), 1, 10, &st)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != 123456 || list[0].Status != TemplateStatusApproved {
		t.Fatalf("list=%+v", list)
	}
}

func TestNotConfigured(t *testing.T) {
	cli := New("", 0)
	_, err := cli.SendSMS(context.Background(), "0912", "x")
	if err != ErrNotConfigured {
		t.Fatalf("err=%v", err)
	}
}
