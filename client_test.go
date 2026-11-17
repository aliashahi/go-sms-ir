package smsir

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestSendViaURL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/send" {
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("X-API-KEY") != "" {
			t.Fatalf("api key=%q", r.Header.Get("X-API-KEY"))
		}
		q := r.URL.Query()
		if q.Get("username") != "user" || q.Get("password") != "secret" {
			t.Fatalf("query=%s", r.URL.RawQuery)
		}
		if q.Get("line") != "30004505000017" || q.Get("mobile") != "09121234567" || q.Get("text") != "hello" {
			t.Fatalf("query=%s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  1,
			"message": "موفق",
			"data":    map[string]any{"messageId": 89545112, "cost": 1.0},
		})
	}))
	defer srv.Close()

	cli := New("", 0, WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
	res, err := cli.SendViaURL(context.Background(), URLSendRequest{
		Username: "user",
		Password: "secret",
		Line:     30004505000017,
		Mobile:   "+989121234567",
		Text:     "hello",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.MessageID != 89545112 || res.Cost != 1 {
		t.Fatalf("res=%+v", res)
	}
}

func TestCancelScheduled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/v1/send/scheduled/pack-1" {
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  1,
			"message": "موفق",
			"data":    map[string]any{"returnedCreditCount": 10, "smsCount": 5},
		})
	}))
	defer srv.Close()

	cli := New("key", 0, WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
	res, err := cli.CancelScheduled(context.Background(), "pack-1")
	if err != nil {
		t.Fatal(err)
	}
	if res.ReturnedCreditCount != 10 || res.SMSCount != 5 {
		t.Fatalf("res=%+v", res)
	}
}

func TestSendVerifyValueTooLong(t *testing.T) {
	cli := New("key", 0)
	_, err := cli.SendVerify(context.Background(), VerifyRequest{
		Mobile:     "09121234567",
		TemplateID: 123456,
		Parameters: []VerifyParameter{{Name: "#CODE#", Value: "12345678901234567890123456"}},
	})
	if err == nil {
		t.Fatal("expected error")
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

func TestTemplate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/templates/123456" {
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  1,
			"message": "موفق",
			"data": map[string]any{
				"id":               123456,
				"title":            "کد تایید ورود",
				"templateText":     "کد تایید شما #CODE# می‌باشد.",
				"type":             1,
				"parameters":       []map[string]any{{"name": "CODE", "description": "کد یکبار مصرف"}},
				"status":           3,
				"rejectionReason":  "متن قالب با نوع انتخاب شده مطابقت ندارد.",
				"creationDateTime": 1628683626,
			},
		})
	}))
	defer srv.Close()

	cli := New("key", 0, WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
	tpl, err := cli.Template(context.Background(), 123456)
	if err != nil {
		t.Fatal(err)
	}
	if tpl.ID != 123456 || tpl.Status != TemplateStatusRejected || tpl.Type != TemplateTypeOTP {
		t.Fatalf("tpl=%+v", tpl)
	}
	if tpl.RejectionReason == nil || *tpl.RejectionReason == "" {
		t.Fatalf("reason=%v", tpl.RejectionReason)
	}
}

func TestAddAndModifyTemplate(t *testing.T) {
	var gotMethod string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotBody = map[string]any{}
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatal(err)
		}
		id := 123456
		if r.Method == http.MethodPut {
			if r.URL.Path != "/v1/templates/123456" {
				t.Fatalf("path=%s", r.URL.Path)
			}
		} else if r.URL.Path != "/v1/templates" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  1,
			"message": "موفق",
			"data":    id,
		})
	}))
	defer srv.Close()

	cli := New("key", 0, WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
	ctx := context.Background()
	id, err := cli.AddTemplate(ctx, TemplateRequest{
		Title:        " کد تایید ورود ",
		TemplateText: "کد تایید شما #CODE# می‌باشد.",
		Type:         TemplateTypeOTP,
		Parameters:   []TemplateParameter{{Name: "#CODE#", Description: "کد یکبار مصرف"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if id != 123456 || gotMethod != http.MethodPost {
		t.Fatalf("id=%d method=%s", id, gotMethod)
	}
	if gotBody["title"] != "کد تایید ورود" || gotBody["type"] != float64(1) {
		t.Fatalf("body=%v", gotBody)
	}
	params, _ := gotBody["parameters"].([]any)
	if len(params) != 1 {
		t.Fatalf("params=%v", gotBody["parameters"])
	}
	param, _ := params[0].(map[string]any)
	if param["name"] != "CODE" {
		t.Fatalf("param=%v", param)
	}

	id, err = cli.ModifyTemplate(ctx, 123456, TemplateRequest{
		Title:        "کد تایید ورود",
		TemplateText: "کد ورود شما #CODE# می‌باشد.",
	})
	if err != nil {
		t.Fatal(err)
	}
	if id != 123456 || gotMethod != http.MethodPut {
		t.Fatalf("id=%d method=%s", id, gotMethod)
	}
	if _, ok := gotBody["type"]; ok {
		t.Fatalf("type should be omitted: %v", gotBody)
	}
	if _, ok := gotBody["parameters"]; ok {
		t.Fatalf("parameters should be omitted: %v", gotBody)
	}
}

func TestRemoveTemplate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/v1/templates/123456" {
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  1,
			"message": "موفق",
			"data":    true,
		})
	}))
	defer srv.Close()

	cli := New("key", 0, WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
	if err := cli.RemoveTemplate(context.Background(), 123456); err != nil {
		t.Fatal(err)
	}
}

func TestAddTemplateRequiresPlaceholder(t *testing.T) {
	cli := New("key", 0)
	_, err := cli.AddTemplate(context.Background(), TemplateRequest{
		Title:        "عنوان",
		TemplateText: "بدون متغیر",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMessageReport(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/send/89545112" {
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("X-API-KEY") != "key" {
			t.Fatalf("api key=%q", r.Header.Get("X-API-KEY"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  1,
			"message": "موفق",
			"data": map[string]any{
				"messageId":        89545112,
				"mobile":           91200000677,
				"messageText":      "سرویس پیامکی ایده پردازان با بیش از یک دهه سابقه همراه شماست",
				"sendDateTime":     1628683626,
				"lineNumber":       30004505000017,
				"cost":             1.0,
				"deliveryState":    1,
				"deliveryDateTime": 1628683629,
			},
		})
	}))
	defer srv.Close()

	cli := New("key", 0, WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
	report, err := cli.MessageReport(context.Background(), 89545112)
	if err != nil {
		t.Fatal(err)
	}
	if report.MessageID != 89545112 || report.Mobile.String() != "91200000677" {
		t.Fatalf("report=%+v", report)
	}
	if report.SendDateTime != 1628683626 || report.LineNumber != 30004505000017 || report.Cost != 1 {
		t.Fatalf("report=%+v", report)
	}
	if report.DeliveryState != DeliveryReachedPhone || report.DeliveryDateTime == nil || *report.DeliveryDateTime != 1628683629 {
		t.Fatalf("delivery=%d %v", report.DeliveryState, report.DeliveryDateTime)
	}
}

func TestTodayPacks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/send/pack" {
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("pageNumber") != "1" || r.URL.Query().Get("pageSize") != "100" {
			t.Fatalf("query=%s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  1,
			"message": "موفق",
			"data": []map[string]any{
				{
					"packId":           "e7c09e23f0db4834b9bcb79e7b054f4c",
					"recipientCount":   100,
					"creationDateTime": 1628683626,
				},
				{
					"packId":           "0cf3017fd9d84babbb5ed5579104dab2",
					"recipientCount":   200,
					"creationDateTime": 1628683626,
				},
			},
		})
	}))
	defer srv.Close()

	cli := New("key", 0, WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
	packs, err := cli.TodayPacks(context.Background(), 100, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(packs) != 2 || packs[0].PackID != "e7c09e23f0db4834b9bcb79e7b054f4c" || packs[0].RecipientCount != 100 {
		t.Fatalf("packs=%+v", packs)
	}
	if packs[1].PackID != "0cf3017fd9d84babbb5ed5579104dab2" || packs[1].CreationDateTime != 1628683626 {
		t.Fatalf("packs=%+v", packs)
	}
}

func TestPackAndInboundReports(t *testing.T) {
	var gotPath, gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		switch r.URL.Path {
		case "/v1/send/pack/pack-1":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"status": 1, "message": "موفق",
				"data": []map[string]any{{
					"messageId": 89545112, "mobile": 91200000677, "messageText": "hi",
					"sendDateTime": 1628683626, "lineNumber": 30004505000017, "cost": 1.0,
					"deliveryState": 1, "deliveryDateTime": 1628683629,
				}},
			})
		case "/v1/receive/latest":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"status": 1, "message": "موفق",
				"data": []map[string]any{{
					"receiveReturnId": 123456789, "messageText": "hi",
					"number": 30004505000017, "mobile": 91200000002, "receivedDateTime": 1628683625,
				}},
			})
		case "/v1/receive/live":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"status": 1, "message": "موفق",
				"data": []map[string]any{{
					"messageText": "hi", "number": 30004505000017,
					"mobile": 91200000002, "receivedDateTime": 1628683625,
				}},
			})
		case "/v1/receive/archive":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"status": 1, "message": "موفق",
				"data": []map[string]any{{
					"receiveReturnId": 12345678987, "messageText": "hi",
					"number": 30004505000017, "mobile": 91200000002, "receivedDateTime": 1628683625,
				}},
			})
		default:
			t.Fatalf("path=%s", r.URL.Path)
		}
	}))
	defer srv.Close()

	cli := New("key", 0, WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
	ctx := context.Background()

	pack, err := cli.PackReport(ctx, "pack-1")
	if err != nil || len(pack) != 1 || pack[0].MessageID != 89545112 {
		t.Fatalf("pack=%v err=%v", pack, err)
	}

	latest, err := cli.LatestReceived(ctx, 20, 30004505000017)
	if err != nil || len(latest) != 1 || latest[0].ReceiveReturnID != 123456789 || latest[0].Number != 30004505000017 {
		t.Fatalf("latest=%v err=%v", latest, err)
	}
	if gotPath != "/v1/receive/latest" || !strings.Contains(gotQuery, "count=20") || !strings.Contains(gotQuery, "lineNumber=30004505000017") {
		t.Fatalf("latest query %s %s", gotPath, gotQuery)
	}

	newest := true
	live, err := cli.TodayReceived(ctx, 20, 3, &newest, "9120000002")
	if err != nil || len(live) != 1 || live[0].Number != 30004505000017 || live[0].Mobile.String() != "91200000002" {
		t.Fatalf("live=%v err=%v", live, err)
	}
	if !strings.Contains(gotQuery, "sortByNewest=true") || !strings.Contains(gotQuery, "mobile=09120000002") {
		t.Fatalf("live query %s", gotQuery)
	}

	from := int64(1613465574)
	archived, err := cli.ArchiveReceived(ctx, &from, nil, 100, 1, "")
	if err != nil || len(archived) != 1 || archived[0].ReceiveReturnID != 12345678987 {
		t.Fatalf("archive=%v err=%v", archived, err)
	}
	if !strings.Contains(gotQuery, "fromDate=1613465574") || strings.Contains(gotQuery, "mobile=") {
		t.Fatalf("archive query %s", gotQuery)
	}
}

func TestNotConfigured(t *testing.T) {
	cli := New("", 0)
	_, err := cli.SendSMS(context.Background(), "0912", "x")
	if err != ErrNotConfigured {
		t.Fatalf("err=%v", err)
	}
}
