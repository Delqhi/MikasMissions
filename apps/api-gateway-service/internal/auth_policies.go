package internal

func requiredRolesForPattern(pattern string) []string {
	switch pattern {
	case "GET /healthz":
		return nil
	case "POST /v1/parents/signup":
		return nil
	case "POST /v1/parents/login":
		return nil
	case "POST /v1/admin/login":
		return nil
	case "POST /v1/parents/consent/verify":
		return nil
	case "GET /v1/parents/dashboard":
		return []string{"parent", "service"}
	case "GET /v1/parents/controls/{child_profile_id}":
		return []string{"parent", "service"}
	case "PUT /v1/parents/controls/{child_profile_id}":
		return []string{"parent", "service"}
	case "POST /v1/parents/gates/challenge":
		return []string{"parent", "service"}
	case "POST /v1/parents/gates/verify":
		return []string{"parent", "service"}
	case "POST /v1/children/profiles":
		return []string{"parent", "service"}
	case "GET /v1/children/profiles":
		return []string{"parent", "service"}
	case "GET /v1/home/rails":
		return []string{"parent", "child", "service"}
	case "GET /v1/kids/home":
		return []string{"parent", "child", "service"}
	case "GET /v1/kids/progress/{child_profile_id}":
		return []string{"parent", "child", "service"}
	case "GET /v1/catalog/episodes/{id}":
		return []string{"parent", "child", "service"}
	case "POST /v1/playback/sessions":
		return []string{"parent", "child", "service"}
	case "POST /v1/progress/watch-events":
		return []string{"parent", "child", "service"}
	case "GET /v1/billing/entitlements":
		return []string{"parent", "service"}
	case "POST /v1/creator/assets/upload":
		return []string{"service"}
	case "GET /v1/admin/workflows":
		return []string{"admin", "service"}
	case "POST /v1/admin/workflows":
		return []string{"admin", "service"}
	case "PUT /v1/admin/workflows/{workflow_id}":
		return []string{"admin", "service"}
	case "DELETE /v1/admin/workflows/{workflow_id}":
		return []string{"admin", "service"}
	case "POST /v1/admin/workflows/{workflow_id}/runs":
		return []string{"admin", "service"}
	case "GET /v1/admin/runs/{run_id}":
		return []string{"admin", "service"}
	case "GET /v1/admin/runs/{run_id}/logs":
		return []string{"admin", "service"}
	case "POST /v1/admin/runs/{run_id}/retry":
		return []string{"admin", "service"}
	case "POST /v1/admin/runs/{run_id}/cancel":
		return []string{"admin", "service"}
	case "GET /v1/admin/model-profiles/{id}":
		return []string{"admin", "service"}
	case "PUT /v1/admin/model-profiles/{id}":
		return []string{"admin", "service"}
	default:
		return nil
	}
}

func roleAllowed(role string, allowed []string) bool {
	if len(allowed) == 0 {
		return true
	}
	for _, candidate := range allowed {
		if role == candidate {
			return true
		}
	}
	return false
}
