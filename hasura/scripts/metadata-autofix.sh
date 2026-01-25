#!/usr/bin/env bash
set -euo pipefail

# Configurable env vars
HASURA_ENDPOINT="${HASURA_ENDPOINT:-http://localhost:8082}"
HASURA_ADMIN_SECRET="${HASURA_ADMIN_SECRET:-admin-secret}"
# ACTION_NAMES: comma-separated list. Falls back to ACTION_NAME (single). If still empty, auto-detect from actions.yaml.
ACTION_NAMES="${ACTION_NAMES:-${ACTION_NAME:-}}"

# Paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
HASURA_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
REPO_ROOT="$(cd "${HASURA_DIR}/.." && pwd)"
# Default to script directory for logs, override with LOG_DIR if set
LOG_DIR="${LOG_DIR:-${SCRIPT_DIR}}"

if [ -z "${ACTION_NAMES}" ]; then
  ACTION_NAMES=$(python3 <<'PY'
import os
import pathlib
import yaml

base = pathlib.Path(os.environ.get("HASURA_DIR", ""))
actions_path = base / "metadata" / "actions.yaml"
fallback = "getCustomerInfo"

if not actions_path.exists():
    print(fallback)
    raise SystemExit

try:
    data = yaml.safe_load(actions_path.read_text())
except Exception:
    print(fallback)
    raise SystemExit

actions = data.get("actions") or []
names = [a.get("name") for a in actions if isinstance(a, dict) and a.get("name")]
print(",".join(names) if names else fallback)
PY
  )
fi

require_pyyaml() {
  python3 - <<'PY'
try:
    import yaml  # noqa: F401
except Exception as e:
    import sys
    sys.stderr.write(f"PyYAML is required. Install with: pip install pyyaml (error: {e})\n")
    sys.exit(1)
PY
}

export_metadata() {
  curl -sS -f -X POST "${HASURA_ENDPOINT}/v1/metadata" \
    -H "Content-Type: application/json" \
    -H "x-hasura-admin-secret: ${HASURA_ADMIN_SECRET}" \
    -d '{"type":"export_metadata","args":{}}'
}

action_exists() {
  local action_name="$1"
  python3 - "$action_name" <<'PY'
import json, os, sys
action = sys.argv[1]
try:
    raw = os.environ.get("ACTION_JSON", "")
    data = json.loads(raw)
except Exception as e:
    sys.stderr.write(f"Invalid JSON input to action_exists: {e}\n")
    sys.exit(2)
meta = data.get("metadata") if isinstance(data, dict) else None
if meta is None:
    meta = data if isinstance(data, dict) else {}
actions = meta.get("actions", [])
names = [a.get("name") for a in actions if isinstance(a, dict)]
sys.exit(0 if action in names else 1)
PY
}

actions_missing() {
  local missing=()
  local parse_error=0
  for action in ${ACTION_NAMES//,/ }; do
    if ACTION_JSON="${export_json}" action_exists "${action}"; then
      continue
    else
      rc=$?
      if [ "$rc" -eq 2 ]; then
        parse_error=1
        break
      fi
      missing+=("${action}")
    fi
  done

  if [ "$parse_error" -eq 1 ]; then
    return 2
  fi

  if [ ${#missing[@]} -gt 0 ]; then
    printf '%s ' "${missing[@]}"
  fi
  return 0
}

build_replace_payload() {
  python3 - "$HASURA_DIR" <<'PY'
import json, pathlib, sys, yaml
base = pathlib.Path(sys.argv[1])
meta_dir = base / "metadata"

def load_yaml(rel_path, default):
    path = meta_dir / rel_path
    if not path.exists():
        return default
    data = yaml.safe_load(path.read_text())
    return default if data is None else data

actions_yaml = load_yaml("actions.yaml", {})
allow_list = load_yaml("allow_list.yaml", [])
query_collections = load_yaml("query_collections.yaml", [])
remote_schemas = load_yaml("remote_schemas.yaml", [])
sources = load_yaml("databases/databases.yaml", [])

def resolve_include_tables(tables_value):
    if isinstance(tables_value, str) and tables_value.startswith("!include "):
        rel = tables_value.split(" ", 1)[1]
        path = meta_dir / rel
        if path.exists():
            loaded = yaml.safe_load(path.read_text())
            return [] if loaded is None else loaded
        return []
    return tables_value

for src in sources:
    if isinstance(src, dict) and "tables" in src:
        src["tables"] = resolve_include_tables(src.get("tables"))

metadata = {
    "version": 3,
    "actions": actions_yaml.get("actions", []),
    "custom_types": actions_yaml.get("custom_types", {}),
    "allowlist": allow_list,
    "query_collections": query_collections,
    "remote_schemas": remote_schemas,
    "sources": sources,
}

payload = {
    "type": "replace_metadata",
    "args": {
        "metadata": metadata,
        "allow_inconsistent_metadata": False,
    },
}
print(json.dumps(payload))
PY
}

apply_replace_metadata() {
  tmp_payload="${LOG_DIR}/hasura-replace-payload-$$.json"
  tmp_response="${LOG_DIR}/hasura-replace-response-$$.txt"
  build_replace_payload >"${tmp_payload}"

  http_status=$(curl -sS -w "\nHTTP_STATUS:%{http_code}" -X POST "${HASURA_ENDPOINT}/v1/metadata" \
    -H "Content-Type: application/json" \
    -H "x-hasura-admin-secret: ${HASURA_ADMIN_SECRET}" \
    -d @"${tmp_payload}" \
    -o "${tmp_response}")

  status_code=$(printf '%s' "${http_status}" | awk -F: '/HTTP_STATUS/ {print $2}')
  if [ "${status_code}" != "200" ]; then
    echo "replace_metadata failed with status ${status_code}" >&2
    echo "Payload: ${tmp_payload}" >&2
    echo "Response body:" >&2
    cat "${tmp_response}" >&2
    return 1
  fi
  return 0
}

main() {
  echo "[1/4] Running export_metadata..." >&2
  if ! export_json=$(export_metadata); then
    echo "Failed to call export_metadata. Check HASURA_ENDPOINT/admin secret or container health." >&2
    exit 1
  fi
  if [ -z "$(printf '%s' "${export_json}" | tr -d ' \r\n\t')" ]; then
    echo "export_metadata returned empty response. Is Hasura running?" >&2
    exit 1
  fi

  echo "[2/4] Checking action existence: ${ACTION_NAMES}" >&2
  missing_actions=$(actions_missing)
  rc=$?
  if [ "$rc" -eq 2 ]; then
    echo "Failed to parse export_metadata response as JSON." >&2
    echo "export_metadata sample (first 200 chars):" >&2
    printf '%s' "${export_json}" | head -c 200 >&2
    echo >&2
    tmp_dump="${LOG_DIR}/hasura-export-metadata-$$.json"
    printf '%s' "${export_json}" >"${tmp_dump}" 2>/dev/null || true
    echo "Full export_metadata written to ${tmp_dump}" >&2
    exit 1
  fi

  if [ -z "${missing_actions}" ]; then
    echo "All actions already exist. replace_metadata not needed." >&2
    exit 0
  fi

  echo "[3/4] Actions missing (${missing_actions}) → running replace_metadata" >&2
  require_pyyaml
  if ! apply_replace_metadata; then
    echo "replace_metadata call failed." >&2
    exit 1
  fi

  echo "[4/4] Re-checking after replace_metadata..." >&2
  if ! export_json=$(export_metadata); then
    echo "Failed to call export_metadata after replace_metadata." >&2
    exit 1
  fi
  missing_actions=$(actions_missing)
  rc=$?
  if [ "$rc" -eq 2 ]; then
    echo "Failed to parse export_metadata response as JSON after replace_metadata." >&2
    echo "export_metadata sample (first 200 chars):" >&2
    printf '%s' "${export_json}" | head -c 200 >&2
    echo >&2
    tmp_dump="${LOG_DIR}/hasura-export-metadata-$$.json"
    printf '%s' "${export_json}" >"${tmp_dump}" 2>/dev/null || true
    echo "Full export_metadata written to ${tmp_dump}" >&2
    exit 1
  fi

  if [ -z "${missing_actions}" ]; then
    echo "All actions restored." >&2
    exit 0
  fi

  echo "Actions still missing after replace_metadata: ${missing_actions}" >&2
  exit 1
}

main "$@"
