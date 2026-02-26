package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/testkit"
)

func TestPostParentGateVerify(t *testing.T) {
	store := NewStore()
	passwordHash, err := hashPassword("1234567890x")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	parent, err := store.CreateParent("parent@example.com", "DE", "de", passwordHash)
	if err != nil {
		t.Fatalf("create parent: %v", err)
	}
	challengeHandler := PostParentGateChallenge(store)
	challengeReq := testkit.NewJSONRequest(http.MethodPost, "/v1/parents/gates/challenge", contractsapi.ParentGateChallengeRequest{
		ParentUserID:   parent.ID,
		ChildProfileID: "child-2",
		Method:         "pin",
	})
	challengeRR := httptest.NewRecorder()
	challengeHandler.ServeHTTP(challengeRR, challengeReq)
	if challengeRR.Code != http.StatusCreated {
		t.Fatalf("expected 201 challenge, got %d", challengeRR.Code)
	}
	var challengeResp contractsapi.ParentGateChallengeResponse
	if err := json.Unmarshal(challengeRR.Body.Bytes(), &challengeResp); err != nil {
		t.Fatalf("unmarshal challenge response: %v", err)
	}

	handler := PostParentGateVerify(store)
	req := testkit.NewJSONRequest(http.MethodPost, "/v1/parents/gates/verify", contractsapi.ParentGateVerifyRequest{
		ParentUserID:   parent.ID,
		ChildProfileID: "child-2",
		ChallengeID:    challengeResp.ChallengeID,
		Response:       "ok",
	})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp contractsapi.ParentGateVerifyResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if !resp.Verified || resp.GateToken == "" {
		t.Fatalf("expected verified response with gate token")
	}

	replayReq := testkit.NewJSONRequest(http.MethodPost, "/v1/parents/gates/verify", contractsapi.ParentGateVerifyRequest{
		ParentUserID:   parent.ID,
		ChildProfileID: "child-2",
		ChallengeID:    challengeResp.ChallengeID,
		Response:       "ok",
	})
	replayRR := httptest.NewRecorder()
	handler.ServeHTTP(replayRR, replayReq)
	if replayRR.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for replayed challenge, got %d", replayRR.Code)
	}
}
