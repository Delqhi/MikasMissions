SHELL := /bin/bash
GOCACHE ?= $(CURDIR)/.cache/go-build
GOENV := GOCACHE=$(GOCACHE)

.PHONY: fmt test lint guard build contract contract-check migrate-db e2e-smoke e2e-auth-smoke e2e-admin-smoke e2e-generator-smoke a11y-smoke compose-smoke launch-preflight launch-go-nogo launch-stage launch-decision-packet launch-readiness-gate today-ready kube-validate staging-soak staging-deploy staging-deploy-dry-run staging-rollback outbox-replay run-gateway run-identity run-profile run-catalog run-playback run-progress run-recommendation run-creator run-admin run-moderation run-billing run-outbox-relay

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './bin/*')

guard:
	./scripts/enforce_micro_files.sh
	./scripts/check_event_contracts.sh
	./scripts/check_api_schema_alignment.sh
	./scripts/check_api_schema_lock.sh

test:
	$(GOENV) go test ./...

lint: guard test

build:
	$(GOENV) go build ./...

contract:
	./scripts/generate_contract_artifacts.sh

contract-check:
	./scripts/check_generated_contract_artifacts.sh
	./scripts/check_openapi_breaking.sh

migrate-db:
	./scripts/apply_sql_migrations.sh

e2e-smoke:
	./scripts/e2e_smoke_with_stack.sh

e2e-auth-smoke:
	AUTH_MODE=enforce SMOKE_SCRIPT=./scripts/smoke_auth_enforce.sh ./scripts/e2e_smoke_with_stack.sh

e2e-admin-smoke:
	AUTH_MODE=enforce SMOKE_SCRIPT=./scripts/smoke_admin.sh ./scripts/e2e_smoke_with_stack.sh

e2e-generator-smoke:
	./scripts/smoke_generator_workers.sh

a11y-smoke:
	./scripts/a11y_smoke.sh

compose-smoke:
	./scripts/compose_smoke.sh

launch-preflight:
	./scripts/launch_preflight.sh

launch-go-nogo:
	./scripts/launch_go_nogo.sh $(ARGS)

launch-stage:
	./scripts/launch_stage_execute.sh $(ARGS)

launch-decision-packet:
	./scripts/launch_decision_packet.sh

launch-readiness-gate:
	./scripts/launch_readiness_gate.sh

today-ready:
	./scripts/today_ready.sh

kube-validate:
	./scripts/check_kube_config.sh

staging-soak:
	./scripts/staging_soak.sh

staging-deploy:
	./scripts/staging_deploy.sh

staging-deploy-dry-run:
	DRY_RUN=true ./scripts/staging_deploy.sh

staging-rollback:
	./scripts/staging_rollback.sh

outbox-replay:
	$(GOENV) go run ./tools/outbox-replay/cmd $(ARGS)

run-identity:
	$(GOENV) go run ./apps/identity-service/cmd

run-gateway:
	$(GOENV) go run ./apps/api-gateway-service/cmd

run-profile:
	$(GOENV) go run ./apps/profile-service/cmd

run-catalog:
	$(GOENV) go run ./apps/catalog-service/cmd

run-playback:
	$(GOENV) go run ./apps/playback-service/cmd

run-progress:
	$(GOENV) go run ./apps/progress-service/cmd

run-recommendation:
	$(GOENV) go run ./apps/recommendation-service/cmd

run-creator:
	$(GOENV) go run ./apps/creator-studio-service/cmd

run-admin:
	$(GOENV) go run ./apps/admin-studio-service/cmd

run-moderation:
	$(GOENV) go run ./apps/moderation-service/cmd

run-billing:
	$(GOENV) go run ./apps/billing-service/cmd

run-outbox-relay:
	$(GOENV) go run ./workers/worker-outbox-relay/cmd
