# Vercel Deployment Contract

Last updated: 2026-02-26

## Canonical Project (ONLY project for this repo)
- Team/Scope: `info-zukunftsories-projects`
- Project name: `web`
- Project ID: `prj_lVZswXfsyRAwcfa9AigZ9C9AMFK9`
- Org ID: `team_VTipbYr7L5qhqXdu38e0Z0OL`
- Local app path: `/Users/jeremy/dev/projects/family-projects/MikasMissions/frontend/web`
- Local project link file: `/Users/jeremy/dev/projects/family-projects/MikasMissions/frontend/web/.vercel/project.json`
- Current production URL: `https://web-cyan-three-26.vercel.app`
- Latest preview URL: `https://web-hf87f67ki-info-zukunftsories-projects.vercel.app`

## Hard Rule
- Never create new Vercel projects for this repository.
- Always deploy changes to the existing `web` project only.

## Standard Deploy Commands
From repo root:

```bash
# Preview deploy (default for normal changes)
vercel deploy /Users/jeremy/dev/projects/family-projects/MikasMissions/frontend/web -y

# Production deploy (only when explicitly requested)
vercel deploy /Users/jeremy/dev/projects/family-projects/MikasMissions/frontend/web --prod -y
```

## Recovery / Re-Link
If the local link is ever wrong, re-link to the canonical project:

```bash
cd /Users/jeremy/dev/projects/family-projects/MikasMissions/frontend/web
vercel link --yes --project web
```

## Duplicate Project Cleanup Log
- 2026-02-26: Removed duplicate project `mikasmissions-premium-pass2`.
