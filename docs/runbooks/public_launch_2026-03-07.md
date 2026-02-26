# Public Launch Runbook (7. März 2026)

## Scope
1. Public Launch als Free-Beta mit Web + PWA.
2. Kernflüsse: Parent Signup/Login/Consent, Child-Profile, Kids-Home, Playback-Enforcement, Progress, Parent Controls.
3. Kein RN-Store-Release, kein Paid-Checkout, keine Growth-Experimente als Launch-Blocker.

## Hard Runtime Defaults
1. `AUTH_MODE=enforce`
2. `GO_ENV=production`
3. `PERSISTENCE_MODE=strict`
4. `NEXT_PUBLIC_USE_API_FALLBACKS=false`

## T-1 Preflight (6. März 2026)
Vom Repo-Root:

```bash
make launch-preflight
```

Evidence packet:

```bash
make launch-decision-packet
```

Zusätzlich:
1. Staging 24h Soak ohne ungeplante SEV-1/SEV-2.
2. Rollback-Drill <15 Minuten dokumentiert.
3. On-Call primär/sekundär + Incident-Kanal bestätigt.
4. Launch readiness Gate ausführen:

```bash
make launch-readiness-gate
```

## Staging Deploy Sequence
Erforderliche Umgebungsvariablen:
1. `DATABASE_URL`
2. `AUTH_JWT_SECRET`
3. Optional: `AUTH_JWKS_URL`

Kubernetes validation:

```bash
make kube-validate
```

Launch readiness must pass with external-cluster policy:

```bash
make launch-readiness-gate
```

`kind|minikube|k3d` are valid for local verification but are blocked for public go-live readiness.

Deploy:

```bash
make staging-deploy
```

Ohne aktiven Kubernetes-Kontext (nur lokale Validierung):

```bash
make staging-deploy-dry-run
```

Rollback (wenn Stop-Kriterium verletzt):

```bash
make staging-rollback
```

## Launch Day Timeline (7. März 2026, CET)
1. 09:00 - Freeze bestätigen, letzte grüne Pipeline bestätigen.
2. 09:30 - Rollout 1%, 30 Minuten beobachten.
3. 10:15 - Rollout 10%, 60 Minuten beobachten.
4. 11:30 - Rollout 50%, 120 Minuten beobachten.
5. 13:45 - Rollout 100% nur bei grünen Stop/Go-Metriken.

## Stop/Go Gates je Stufe
1. `5xx < 1.0%`
2. `p95 read < 400ms`
3. `p95 write < 700ms`

Beispielbewertung:

```bash
make launch-go-nogo ARGS="--stage 10 --read-p95-ms 285 --write-p95-ms 430 --error-5xx-percent 0.22"
```

Empfohlene Stufen-Ausführung inkl. Evidence + Auto-Rollback:

```bash
make launch-stage ARGS="--stage 10 --read-p95-ms 285 --write-p95-ms 430 --error-5xx-percent 0.22 --owner sre-oncall --notes rollout-window-2"
```

Wenn der Befehl mit Exit-Code `1` endet:
1. Sofortiger Rollback auf letzten stabilen Tag.
2. Incident-Prozess starten.
3. Nächste Rollout-Stufe sperren.

## Launch-Blocker Tests (müssen grün sein)
1. `make guard`
2. `make test`
3. `make build`
4. `make contract-check`
5. `make e2e-smoke`
6. `make e2e-auth-smoke`
7. `cd frontend/web && NEXT_PUBLIC_USE_API_FALLBACKS=false npm run build`

## Security-Spezifische Go/No-Go Checks
1. Cross-Parent Zugriff blockiert (`child_profile_forbidden`).
2. Parent-Gate Token-Reuse blockiert (`parent_gate_required`).
3. Entitlement-Bypass durch manipuliertes Payload blockiert (`entitlement_required`).
4. Protected Endpoints in `AUTH_MODE=enforce` ohne Token -> `401`.

## Rollback Procedure
1. Stopp der weiteren Traffic-Erhöhung.
2. Deployment-Rollback ausführen (`make staging-rollback`).
3. Health + Smoke erneut prüfen (`make e2e-smoke`).
4. Incident-Timeline, Ursache, next action dokumentieren.

## Evidence (Pflicht)
1. Screenshot/Log jeder Rollout-Stufe.
2. Gemessene Werte für Read-p95, Write-p95, 5xx.
3. Ergebnis `launch-go-nogo` pro Stufe.
4. Finaler Go/No-Go Entscheid mit Owner und Uhrzeit.
