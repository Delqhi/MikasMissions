# MikasMissions Abschlussplan (7. März bis 19. Juni 2026)

## Ziel
1. Public Go-Live am 7. März 2026 mit hartem Scope.
2. Danach Stabilitäts- und Produktisierungsprogramm bis 19. Juni 2026.

## Fixe Regeln
1. Kein Scope-Creep vor stabilem Public-Betrieb.
2. `/v1` bleibt bis nach Launch eingefroren.
3. Kein Public-Go-Live ohne externes Kubernetes-Staging.
4. Runtime in Produktion: `AUTH_MODE=enforce`, `GO_ENV=production`, `PERSISTENCE_MODE=strict`, `NEXT_PUBLIC_USE_API_FALLBACKS=false`.

## Phase 0 (Cluster-Freigabe)
1. Externes Kubeconfig aktivieren.
2. `KUBECONFIG=<external> make kube-validate`
3. `KUBECONFIG=<external> make staging-deploy-dry-run`
4. `make launch-readiness-gate`

Gate:
1. `kube-validate` grün.
2. Dry-run grün.
3. Readiness-Status `READY_FOR_STAGING_DEPLOY`.

## Phase 1 (Staging Hardening)
1. `DATABASE_URL=... AUTH_JWT_SECRET=... make migrate-db`
2. `DATABASE_URL=... AUTH_JWT_SECRET=... make staging-deploy`
3. `make launch-preflight`
4. `SOAK_HOURS=24 SOAK_INTERVAL_MINUTES=60 make staging-soak`
5. `make staging-rollback`
6. `make e2e-smoke`
7. `make e2e-auth-smoke`

Gate:
1. 24h ohne ungeplante SEV-1/SEV-2.
2. Rollback <15 Minuten nachgewiesen.
3. Pflichtchecks grün.

## Phase 2 (Launch Decision)
1. `make launch-decision-packet`
2. Artefakte prüfen:
   1. `docs/runbooks/evidence/launch-*/decision_packet.md`
   2. `docs/runbooks/evidence/launch-*/launch_preflight.log`
   3. `docs/runbooks/evidence/staging-soak-*/summary.md`
3. Go/No-Go mit Owner, Uhrzeit und Metriken freigeben.

Gate:
1. Vollständiges Decision-Packet.
2. On-call primär/sekundär benannt.
3. Keine offenen SEV-1/SEV-2 Risiken.

## Phase 3 (Public Launch Day, 7. März 2026)
1. 09:30 CET: Rollout 1%
2. 10:15 CET: Rollout 10%
3. 11:30 CET: Rollout 50%
4. 13:45 CET: Rollout 100%

Pro Stufe:
1. `make launch-stage ARGS="--stage <1|10|50|100> --read-p95-ms <x> --write-p95-ms <y> --error-5xx-percent <z> --owner <name> --notes <window>"`

Stop/Go-Kriterien:
1. `5xx < 1.0%`
2. `p95 read < 400ms`
3. `p95 write < 700ms`

Verletzung:
1. Sofortiger Rollback via `make staging-rollback`.
2. Incident starten.
3. Nächste Stufe sperren.

## Phase 4 (Post-Launch bis 19. Juni 2026)
1. 30 Tage ohne ungeplante SEV-1.
2. Wöchentliche Reliability- und Performance-Regressionen.
3. Admin/NIM intern weiter produktisieren:
   1. Workflow-Templates standardisieren.
   2. Provider-Profile DB-basiert pflegen.
   3. QC-Fail blockierend halten.
4. Public-Scope erst nach Stabilitätsfenster erweitern.

## Pflichttests (Launch-Blocker)
1. `make guard`
2. `make test`
3. `make build`
4. `make contract-check`
5. `make e2e-smoke`
6. `make e2e-auth-smoke`
7. `make e2e-admin-smoke`
8. `make e2e-generator-smoke`
9. `make a11y-smoke`
10. `cd frontend/web && NEXT_PUBLIC_USE_API_FALLBACKS=false npm run build`

## KPI/SLO Gates
1. Launch: `p95 read < 400ms`, `p95 write < 700ms`, `5xx < 1%`.
2. Post-Launch Ziel: `p95 read < 250ms`, `p95 write < 500ms`.
3. UX: Parent-Onboarding >90% Task Success, Kinder-Navigation <8% Errors.

## Verantwortlichkeiten
1. Montag: KPI-Review.
2. Mittwoch: Risk-Review.
3. Freitag: Release-Readiness.
4. ADR-Pflicht vor Merge für Architektur-/Policy-Änderungen.
5. Incident-Policy: SEV-1..SEV-3, Postmortem <48h.
