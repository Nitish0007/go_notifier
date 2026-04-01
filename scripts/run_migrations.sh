#!/bin/sh
set -e

if [ -z "$DB_URL" ]; then
  echo "DB_URL is not set"
  exit 1
fi

run_migrate_up() {
  migrate -path /migrations -database "$DB_URL" up
}

echo "Running migrations..."
if run_migrate_up; then
  echo "Migrations completed successfully"
  exit 0
fi

# Migration failed: capture output and check for dirty state
echo "Migration failed, checking for dirty state..."
UP_OUTPUT=$(run_migrate_up 2>&1) || true

if ! echo "$UP_OUTPUT" | grep -qi "dirty"; then
  echo "Migration failed (not due to dirty state):"
  echo "$UP_OUTPUT"
  exit 1
fi

# Parse dirty version from migrate output (e.g. "Dirty database version 3. Fix and force version." -> 3)
DIRTY_VERSION=$(echo "$UP_OUTPUT" | sed -n 's/.*[Dd]irty.*version \([0-9][0-9]*\).*/\1/p' | head -1)
if [ -z "$DIRTY_VERSION" ]; then
  echo "Could not parse dirty version from output:"
  echo "$UP_OUTPUT"
  exit 1
fi

PREV_VERSION=$((DIRTY_VERSION - 1))
[ "$PREV_VERSION" -lt 0 ] && PREV_VERSION=0
echo "Resolving dirty migration: forcing to version $PREV_VERSION (so migration $DIRTY_VERSION will re-run)..."
migrate -path /migrations -database "$DB_URL" force "$PREV_VERSION"
echo "Re-running migrations (will re-apply migration $DIRTY_VERSION and continue)..."
run_migrate_up

echo "Migrations completed successfully"