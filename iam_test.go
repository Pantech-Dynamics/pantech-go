package pantech

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
)

const policyJSON = `{"id":"pol_1","name":"staging-operators","description":"","managed":false,
"statements":[{"sid":"RestartStaging","effect":"allow","actions":["instances:Read","instances:Reboot"],"resources":["name:staging-*"],
"conditions":{"source_ip":["102.89.0.0/16"],"mfa":false}}],"attachment_count":3,
"created_at":"2026-11-10T09:00:00Z","updated_at":"2026-11-10T10:00:00Z","version":5}`

func TestAllowAndDenyBuildStatements(t *testing.T) {
	s := Deny([]string{"instances:Delete"}, []string{"name:prod-*"}).WithSid("NoProdDeletes").WithMFA()
	data, _ := json.Marshal(s)
	if string(data) != `{"sid":"NoProdDeletes","effect":"deny","actions":["instances:Delete"],"resources":["name:prod-*"],"conditions":{"mfa":true}}` {
		t.Fatalf("statement %s", data)
	}
	a := Allow([]string{"*:Read"}, nil)
	if a.Effect != EffectAllow || !reflect.DeepEqual(a.Resources, []string{"*"}) || a.Conditions != nil {
		t.Fatalf("statement %+v", a)
	}
	// The With methods never change the statement they were called on.
	base := Allow([]string{"instances:Read"}, nil).WithSourceIP("10.0.0.0/8")
	office := base.WithSourceIP("102.89.0.0/16")
	if len(base.Conditions.SourceIP) != 1 || len(office.Conditions.SourceIP) != 2 || base.Conditions.MFA != nil {
		t.Fatalf("base %+v office %+v", base.Conditions, office.Conditions)
	}
	if mfa := base.WithMFA(); base.Conditions.MFA != nil || !*mfa.Conditions.MFA {
		t.Fatal("WithMFA changed the original")
	}
	// Allow copies its arguments.
	actions := []string{"dns:Read"}
	s2 := Allow(actions, nil)
	actions[0] = "*"
	if s2.Actions[0] != "dns:Read" {
		t.Fatal("Allow kept the caller's slice")
	}
}

func TestIAMCreatePolicySendsTheStatements(t *testing.T) {
	rec := newRecorder(t, status(201, policyJSON))
	c := newTestClientAt(t, rec.srv.URL)
	p, err := c.IAM.CreatePolicy(context.Background(), CreateIAMPolicyRequest{
		Name:       "staging-operators",
		Statements: []PolicyStatement{Allow([]string{"instances:Read", "instances:Reboot"}, []string{"name:staging-*"}).WithSid("RestartStaging").WithSourceIP("102.89.0.0/16")},
	})
	if err != nil {
		t.Fatal(err)
	}
	r := rec.request(0)
	if r.Method != "POST" || r.URL.Path != "/public/v1/iam/policies" || r.Header.Get("Idempotency-Key") == "" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
	want := `{"name":"staging-operators","statements":[{"sid":"RestartStaging","effect":"allow","actions":["instances:Read","instances:Reboot"],"resources":["name:staging-*"],"conditions":{"source_ip":["102.89.0.0/16"]}}]}`
	if rec.bodies[0] != want {
		t.Fatalf("body %s", rec.bodies[0])
	}
	if p.Version != 5 || p.Statements[0].Conditions == nil || *p.Statements[0].Conditions.MFA || p.AttachmentCount != 3 {
		t.Fatalf("policy %+v", p)
	}
}

func TestIAMUpdatePolicySendsIfMatch(t *testing.T) {
	rec := newRecorder(t, ok(policyJSON))
	c := newTestClientAt(t, rec.srv.URL)
	if _, err := c.IAM.UpdatePolicy(context.Background(), "pol_1", UpdateIAMPolicyRequest{Name: "x", Statements: []PolicyStatement{Allow([]string{"*:Read"}, nil)}}, WithIfMatch(5)); err != nil {
		t.Fatal(err)
	}
	r := rec.request(0)
	if r.Method != "PUT" || r.URL.Path != "/public/v1/iam/policies/pol_1" || r.Header.Get("If-Match") != "5" {
		t.Fatalf("sent %s %s If-Match %q", r.Method, r.URL.Path, r.Header.Get("If-Match"))
	}
}

func TestIAMInvalidPolicyNamesTheField(t *testing.T) {
	rec := newRecorder(t, status(422, `{"title":"Invalid policy","code":"INVALID_POLICY","detail":"statements[0].actions[0]: unknown action.","errors":[{"field":"statements[0].actions[0]","code":"unknown_action","message":"unknown action"}]}`))
	c := newTestClientAt(t, rec.srv.URL)
	_, err := c.IAM.CreatePolicy(context.Background(), CreateIAMPolicyRequest{Name: "x", Statements: []PolicyStatement{Allow([]string{"instances:Explode"}, nil)}})
	e, ok := AsError(err)
	if !ok || e.Code != CodeInvalidPolicy || !e.HasField("statements[0].actions[0]", "unknown_action") {
		t.Fatalf("err %v", err)
	}
}

func TestIAMListPoliciesManagedFilter(t *testing.T) {
	rec := newRecorder(t, ok(`{"data":[`+policyJSON+`],"next_cursor":null}`), ok(`{"data":[],"next_cursor":null}`))
	c := newTestClientAt(t, rec.srv.URL)
	if _, err := c.IAM.ListPolicies(context.Background(), &ListIAMPoliciesParams{Managed: Ptr(false)}); err != nil {
		t.Fatal(err)
	}
	if q := rec.request(0).URL.RawQuery; q != "managed=false" {
		t.Fatalf("query %q", q)
	}
	if _, err := c.IAM.ListPolicies(context.Background(), nil); err != nil || rec.request(1).URL.RawQuery != "" {
		t.Fatalf("query %q err %v", rec.request(1).URL.RawQuery, err)
	}
}

func TestIAMTeamMembers(t *testing.T) {
	rec := newRecorder(t,
		ok(`{"user_id":"usr_1","email":"ada@example.com","name":"Ada Obi","role":"member","added_at":"2026-11-10T09:00:00Z"}`),
		status(204, ""),
	)
	c := newTestClientAt(t, rec.srv.URL)
	m, err := c.IAM.AddTeamMember(context.Background(), "team_1", "usr_1")
	if err != nil || m.Email != "ada@example.com" {
		t.Fatalf("member %+v err %v", m, err)
	}
	if r := rec.request(0); r.Method != "PUT" || r.URL.Path != "/public/v1/iam/teams/team_1/members/usr_1" || r.Header.Get("Idempotency-Key") == "" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
	if err := c.IAM.RemoveTeamMember(context.Background(), "team_1", "usr_1"); err != nil {
		t.Fatal(err)
	}
	if r := rec.request(1); r.Method != "DELETE" || r.URL.Path != "/public/v1/iam/teams/team_1/members/usr_1" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
}

func TestIAMAttachAndDetach(t *testing.T) {
	rec := newRecorder(t,
		status(201, `{"policy_id":"pol_1","principal_type":"team","principal_id":"team_1","attached_at":"2026-11-10T09:00:00Z"}`),
		status(204, ""),
	)
	c := newTestClientAt(t, rec.srv.URL)
	att, err := c.IAM.Attach(context.Background(), "pol_1", AttachIAMPolicyRequest{PrincipalType: PrincipalTeam, PrincipalID: "team_1"})
	if err != nil || att.PrincipalID != "team_1" {
		t.Fatalf("attachment %+v err %v", att, err)
	}
	if rec.bodies[0] != `{"principal_type":"team","principal_id":"team_1"}` {
		t.Fatalf("body %s", rec.bodies[0])
	}
	if err := c.IAM.Detach(context.Background(), "pol_1", PrincipalMember, "usr_1"); err != nil {
		t.Fatal(err)
	}
	if r := rec.request(1); r.Method != "DELETE" || r.URL.Path != "/public/v1/iam/policies/pol_1/attachments/member/usr_1" {
		t.Fatalf("sent %s %s", r.Method, r.URL.Path)
	}
}

func TestIAMAdminTargetIsACode(t *testing.T) {
	rec := newRecorder(t, status(403, `{"title":"Forbidden","code":"IAM_ADMIN_TARGET","detail":"Only owners can attach policies to admins."}`))
	c := newTestClientAt(t, rec.srv.URL)
	_, err := c.IAM.Attach(context.Background(), "pol_1", AttachIAMPolicyRequest{PrincipalType: PrincipalMember, PrincipalID: "usr_admin"})
	if !IsCode(err, CodeIAMAdminTarget) || rec.count() != 1 {
		t.Fatalf("err %v after %d requests", err, rec.count())
	}
}

func TestIAMSimulate(t *testing.T) {
	rec := newRecorder(t, ok(`{"allowed":false,"decision":"explicit_deny","policy_id":"pol_1","sid":"NoProdDeletes"}`))
	c := newTestClientAt(t, rec.srv.URL)
	sim, err := c.IAM.Simulate(context.Background(), SimulateIAMAccessRequest{UserID: "usr_1", Action: "instances:Delete", ResourceName: Ptr("prod-api"), MFA: Ptr(true)})
	if err != nil {
		t.Fatal(err)
	}
	if sim.Allowed || sim.Decision != DecisionExplicitDeny || *sim.Sid != "NoProdDeletes" {
		t.Fatalf("simulation %+v", sim)
	}
	if rec.bodies[0] != `{"user_id":"usr_1","action":"instances:Delete","resource_name":"prod-api","mfa":true}` {
		t.Fatalf("body %s", rec.bodies[0])
	}
}

func TestIAMEffectiveAccessAndSettings(t *testing.T) {
	rec := newRecorder(t,
		ok(`{"user_id":"usr_1","role":"member","teams":["team_1"],"policies":[{"policy_id":"pol_1","team_id":"team_1"},{"policy_id":"pol_pantech_read_only","team_id":null}],"member_baseline":"none"}`),
		ok(`{"member_baseline":"none","version":2,"updated_at":"2026-11-10T09:00:00Z"}`),
	)
	c := newTestClientAt(t, rec.srv.URL)
	ea, err := c.IAM.GetEffectiveAccess(context.Background(), "usr_1")
	if err != nil || rec.request(0).URL.Path != "/public/v1/iam/members/usr_1/effective-access" {
		t.Fatalf("err %v", err)
	}
	if len(ea.Policies) != 2 || *ea.Policies[0].TeamID != "team_1" || ea.Policies[1].TeamID != nil || ea.MemberBaseline != MemberBaselineNone {
		t.Fatalf("effective access %+v", ea)
	}
	if _, err := c.IAM.UpdateSettings(context.Background(), UpdateIAMSettingsRequest{MemberBaseline: MemberBaselineNone}); err != nil {
		t.Fatal(err)
	}
	if r := rec.request(1); r.Method != "PUT" || r.URL.Path != "/public/v1/iam/settings" || rec.bodies[1] != `{"member_baseline":"none"}` {
		t.Fatalf("sent %s %s %s", r.Method, r.URL.Path, rec.bodies[1])
	}
}

func TestIAMAllMembersSendsQAndFollowsTheCursor(t *testing.T) {
	m := `{"user_id":"usr_1","email":"ada@example.com","name":"Ada Obi","role":"member","team_ids":["team_1"],"policy_count":2}`
	rec := newRecorder(t, ok(`{"data":[`+m+`],"next_cursor":"c2"}`), ok(`{"data":[`+m+`],"next_cursor":null}`))
	c := newTestClientAt(t, rec.srv.URL)
	n := 0
	for member, err := range c.IAM.AllMembers(context.Background(), &ListIAMMembersParams{Q: "ada"}) {
		if err != nil {
			t.Fatal(err)
		}
		if member.Email != "ada@example.com" || member.PolicyCount != 2 || member.TeamIDs[0] != "team_1" {
			t.Fatalf("member %+v", member)
		}
		n++
	}
	r := rec.request(1)
	if n != 2 || r.URL.Path != "/public/v1/iam/members" || r.URL.Query().Get("q") != "ada" || r.URL.Query().Get("cursor") != "c2" {
		t.Fatalf("%d members, second request %s?%s", n, r.URL.Path, r.URL.RawQuery)
	}
}
