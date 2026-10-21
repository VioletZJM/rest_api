package routes

import (
	"testing"
)

func TestRegisterForEvent(t *testing.T) {
	e := newTestServer(t)

	token := signupAndLogin(t, e, "test@example.com")
	eventID := createTestEvent(t, e, token)

	rec := doTest(t, e, "POST", eventPath(eventID, "/register"), token, "")
	assertStatusCode(t, rec, 201)

	rec = doTest(t, e, "POST", eventPath(eventID, "/register"), token, "")
	assertStatusCode(t, rec, 409)

	ids := myRegistrationIDs(t, e, token)
	if len(ids) != 1 {
		t.Fatalf("Expected 1 registration, but got %d", len(ids))
	}
	if ids[0] != eventID {
		t.Fatalf("Expected event_id %d, but got %d", eventID, ids[0])
	}
}

func TestCancelRegistration(t *testing.T) {
	e := newTestServer(t)

	token := signupAndLogin(t, e, "test@example.com")
	eventID := createTestEvent(t, e, token)

	rec := doTest(t, e, "POST", eventPath(eventID, "/register"), token, "")
	assertStatusCode(t, rec, 201)

	rec = doTest(t, e, "DELETE", eventPath(eventID, "/register"), token, "")
	assertStatusCode(t, rec, 200)

	rec = doTest(t, e, "DELETE", eventPath(eventID, "/register"), token, "")
	assertStatusCode(t, rec, 404)

	ids := myRegistrationIDs(t, e, token)
	if len(ids) != 0 {
		t.Fatalf("Expected 0 registrations, but got %d (%v)", len(ids), ids)
	}
}

func TestGetMyRegistrations(t *testing.T) {
	e := newTestServer(t)

	tokenA := signupAndLogin(t, e, "testA@example.com")
	tokenB := signupAndLogin(t, e, "testB@example.com")

	eventID := createTestEvent(t, e, tokenA)

	rec := doTest(t, e, "POST", eventPath(eventID, "/register"), tokenB, "")
	assertStatusCode(t, rec, 201)

	idsB := myRegistrationIDs(t, e, tokenB)
	if len(idsB) != 1 {
		t.Fatalf("Expected 1 registration for user B, but got %d", len(idsB))
	}
	if idsB[0] != eventID {
		t.Fatalf("Expected event_id %d for user B, but got %d", eventID, idsB[0])
	}

	idsA := myRegistrationIDs(t, e, tokenA)
	if len(idsA) != 0 {
		t.Fatalf("Expected 0 registrations for user A, but got %d (%v)", len(idsA), idsA)
	}
}

func TestUpdateEventNotOwner(t *testing.T) {
	e := newTestServer(t)

	tokenA := signupAndLogin(t, e, "testA@example.com")
	tokenB := signupAndLogin(t, e, "testB@example.com")

	eventID := createTestEvent(t, e, tokenA)

	rec := doTest(t, e, "PUT", eventPath(eventID, ""), tokenB, `{"name":"Updated Title"}`)
	assertStatusCode(t, rec, 401)
}
