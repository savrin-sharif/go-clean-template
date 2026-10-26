#!/usr/bin/env bash
# Creates a clean-architecture Go REST API project from this template.

if [[ -z "${BASH_VERSION:-}" ]]; then
  exec /usr/bin/env bash "$0" "$@"
fi

set -euo pipefail

die() { printf 'Error: %s\n' "$*" >&2; exit 1; }
need() { command -v "$1" >/dev/null 2>&1 || die "Missing dependency: $1"; }

TEMPLATE_MODULE="github.com/savrin-sharif/go-clean-template"
TEMPLATE_GIT_URL="${TEMPLATE_GIT_URL:-https://github.com/savrin-sharif/go-clean-template.git}"

CYAN='\033[38;5;51m'
GREEN='\033[38;5;84m'
RED='\033[38;5;203m'
WARNING='\033[38;5;214m'
MUTED='\033[38;5;245m'
BOLD='\033[1m'
RESET='\033[0m'
HIDE_CURSOR='\033[?25l'
SHOW_CURSOR='\033[?25h'
BLINK_CURSOR='\033[5 q'
DEFAULT_CURSOR='\033[0 q'
MODERN_UI=false
[[ -t 0 && -t 1 && "${TERM:-dumb}" != dumb && -z "${NO_COLOR:-}" ]] && MODERN_UI=true

warning() {
  if [[ "$MODERN_UI" == true ]]; then
    printf '%b⚠  %s%b\n' "$WARNING$BOLD" "$*" "$RESET" >&2
  else
    printf 'Warning: %s\n' "$*" >&2
  fi
}

restore_terminal() {
  if [[ "$MODERN_UI" == true ]]; then printf '%b%b' "$SHOW_CURSOR" "$DEFAULT_CURSOR"; fi
}
exit_program() {
  restore_terminal
  printf '\nInitializer closed.\n'
  exit 0
}
interrupt_program() {
  restore_terminal
  printf '\nInitializer cancelled.\n' >&2
  exit 130
}
trap restore_terminal EXIT
trap interrupt_program INT TERM

START_DIR="$(pwd)"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

show_intro() {
  [[ "$MODERN_UI" == true ]] || return 0
  printf '\033[2J\033[H'
  printf '%b╭──────────────────────────────────────────╮%b\n' "$CYAN" "$RESET"
  printf '%b│                                          │%b\n' "$CYAN" "$RESET"
  printf '%b│           G O   C L E A N   A P I        │%b\n' "$CYAN$BOLD" "$RESET"
  printf '%b│                                          │%b\n' "$CYAN" "$RESET"
  printf '%b╰──────────────────────────────────────────╯%b\n' "$CYAN" "$RESET"
  printf '%bGenerate a clean-architecture Go REST API project.%b\n\n' "$MUTED" "$RESET"
  printf '%bEsc or Ctrl+C exits at any time.%b\n\n' "$MUTED" "$RESET"
  printf '%b%b' "$SHOW_CURSOR" "$BLINK_CURSOR"
}

select_one() {
  local options=("$@") index=0 key first_render=true count=$# i
  printf '%b' "$HIDE_CURSOR"
  while true; do
    if [[ "$first_render" != true ]]; then printf '\033[%dA' "$count"; fi
    first_render=false
    for i in "${!options[@]}"; do
      printf '\033[2K'
      if [[ "$i" -eq "$index" ]]; then
        printf '        %b◆%b  %b%s%b\n' "$CYAN$BOLD" "$RESET" "$GREEN$BOLD" "${options[$i]}" "$RESET"
      else
        printf '        %b◇  %s%b\n' "$MUTED" "${options[$i]}" "$RESET"
      fi
    done
    IFS= read -rsn1 key
    [[ -z "$key" ]] && break
    if [[ "$key" == $'\x1b' ]]; then
      if ! IFS= read -rsn1 -t 1 _; then exit_program; fi
      if ! IFS= read -rsn1 -t 1 key; then exit_program; fi
      [[ "$key" == A ]] && index=$(( (index - 1 + count) % count ))
      [[ "$key" == B ]] && index=$(( (index + 1) % count ))
    fi
  done
  SELECTED_INDEX="$index"
  printf '%b%b' "$SHOW_CURSOR" "$BLINK_CURSOR"
}

read_input() {
  local prompt="$1" default="${2:-}"
  INPUT_ACTION=ok
  if [[ -n "$default" ]]; then
    printf '%b%s [%s]: %b' "$CYAN" "$prompt" "$default" "$RESET"
  else
    printf '%b%s: %b' "$CYAN" "$prompt" "$RESET"
  fi
  read -r INPUT_VALUE || exit_program
  if [[ -z "$INPUT_VALUE" && -n "$default" ]]; then
    INPUT_VALUE="$default"
  fi
}

prompt_text() {
  local prompt="$1" default="${2:-}"
  if [[ "$MODERN_UI" == true ]]; then
    read_input "$prompt" "$default"
  else
    if [[ -n "$default" ]]; then
      read -r -p "$prompt [$default]: " INPUT_VALUE
    else
      read -r -p "$prompt: " INPUT_VALUE
    fi
    INPUT_VALUE="${INPUT_VALUE:-$default}"
  fi
}

choose_database() {
  local options=('none — local welcome repository (no database)' 'postgres — wire pgx + DATABASE_URL')
  local values=('none' 'postgres')
  if [[ "$MODERN_UI" == true ]]; then
    printf '\n        %bChoose database support%b\n\n' "$MUTED" "$RESET"
    select_one "${options[@]}"
    DATABASE_CHOICE="${values[$SELECTED_INDEX]}"
  else
    printf '\nDatabase support:\n'
    printf '  1) none (default)\n'
    printf '  2) postgres\n'
    while true; do
      read -r -p 'Selection [1]: ' selection
      selection="${selection:-1}"
      case "$selection" in
        1) DATABASE_CHOICE=none; break ;;
        2) DATABASE_CHOICE=postgres; break ;;
        *) warning 'Invalid selection.' ;;
      esac
    done
  fi
}

print_review() {
  printf '\n%bReview configuration%b\n' "$BOLD" "$RESET"
  printf '  Project folder : %s\n' "$PROJECT_NAME"
  printf '  Module path    : %s\n' "$NEW_MODULE"
  printf '  App name       : %s\n' "$APP_NAME"
  printf '  HTTP port      : %s\n' "$HTTP_PORT"
  printf '  CORS origins   : %s\n' "$HTTP_ALLOWED_ORIGINS"
  printf '  Database       : %s\n' "$DATABASE_CHOICE"
  printf '  Output         : %s\n\n' "$TARGET_DIR"
}

confirm_generate() {
  if [[ "${YES:-}" == "1" || "${CI:-}" == "true" || "${CI:-}" == "1" ]]; then
    return 0
  fi
  if [[ "$MODERN_UI" == true ]]; then
    printf '        %bProceed?%b\n\n' "$MUTED" "$RESET"
    select_one 'Generate project' 'Exit without generating'
    [[ "$SELECTED_INDEX" -eq 0 ]] || exit_program
  else
    read -r -p 'Generate project? [Y/n]: ' answer
    answer="$(printf '%s' "${answer:-Y}" | tr '[:upper:]' '[:lower:]')"
    case "$answer" in
      y|yes|'') ;;
      *) exit_program ;;
    esac
  fi
}

need go
need find
need sed

show_intro

WORK_DIR="$(mktemp -d)"
cleanup_work() { rm -rf "$WORK_DIR"; }
trap 'restore_terminal; cleanup_work' EXIT

if [[ -f "$TEMPLATE_ROOT/go.mod" ]] && grep -Eq "^module[[:space:]]+${TEMPLATE_MODULE//\//\\/}$" "$TEMPLATE_ROOT/go.mod"; then
  if command -v rsync >/dev/null 2>&1; then
    rsync -a --exclude '.git' --exclude '.tmp' --exclude 'bin' "$TEMPLATE_ROOT"/ "$WORK_DIR"/
  else
    cp -R "$TEMPLATE_ROOT"/. "$WORK_DIR"/
    rm -rf "$WORK_DIR/.git" "$WORK_DIR/.tmp" "$WORK_DIR/bin"
  fi
else
  need git
  printf 'Template files not found next to script; cloning %s ...\n' "$TEMPLATE_GIT_URL" >&2
  git clone --depth 1 "$TEMPLATE_GIT_URL" "$WORK_DIR" >&2 || die "Failed to clone template."
fi

cd "$WORK_DIR"

NEW_MODULE="${1:-${NEW_MODULE:-}}"
DEFAULT_NAME="my-service"
if [[ -n "$NEW_MODULE" ]]; then
  DEFAULT_NAME="${NEW_MODULE##*/}"
fi

if [[ -z "${PROJECT_NAME:-}" ]]; then
  prompt_text 'Project folder name' "$DEFAULT_NAME"
  PROJECT_NAME="$INPUT_VALUE"
fi
PROJECT_NAME="${PROJECT_NAME:-$DEFAULT_NAME}"

if [[ -z "$NEW_MODULE" ]]; then
  DEFAULT_MODULE="github.com/yourname/${PROJECT_NAME}"
  prompt_text 'Go module path' "$DEFAULT_MODULE"
  NEW_MODULE="$INPUT_VALUE"
fi

if [[ -z "${APP_NAME:-}" ]]; then
  prompt_text 'App display name' "$PROJECT_NAME"
  APP_NAME="$INPUT_VALUE"
fi
APP_NAME="${APP_NAME:-$PROJECT_NAME}"

if [[ -z "${HTTP_PORT:-}" ]]; then
  prompt_text 'HTTP port' '8080'
  HTTP_PORT="$INPUT_VALUE"
fi
HTTP_PORT="${HTTP_PORT:-8080}"

if [[ -z "${HTTP_ALLOWED_ORIGINS:-}" ]]; then
  prompt_text 'CORS allowed origins (comma-separated)' '*'
  HTTP_ALLOWED_ORIGINS="$INPUT_VALUE"
fi
HTTP_ALLOWED_ORIGINS="${HTTP_ALLOWED_ORIGINS:-*}"

if [[ -z "${DATABASE_CHOICE:-}" ]]; then
  choose_database
fi
DATABASE_CHOICE="${DATABASE_CHOICE:-none}"

if [[ "$NEW_MODULE" != */* ]] || [[ "$NEW_MODULE" == *" "* ]]; then
  die "Invalid module path: '${NEW_MODULE}' (example: github.com/yourname/${PROJECT_NAME})"
fi
if [[ "$PROJECT_NAME" == *"/"* ]] || [[ "$PROJECT_NAME" == *" "* ]]; then
  die "Invalid project folder name: '${PROJECT_NAME}'"
fi
if ! [[ "$HTTP_PORT" =~ ^[0-9]+$ ]] || (( HTTP_PORT < 1 || HTTP_PORT > 65535 )); then
  die "Invalid HTTP port: '${HTTP_PORT}'"
fi

TARGET_DIR="${START_DIR%/}/${PROJECT_NAME}"
print_review
confirm_generate

if [[ -e "$TARGET_DIR" ]]; then
  die "Target directory already exists: ${TARGET_DIR}"
fi

GO_VERSION_FULL="$(go env GOVERSION)"
GO_VERSION_SHORT="${GO_VERSION_FULL#go}"
GO_VERSION_MAJOR_MINOR="$(printf '%s' "$GO_VERSION_SHORT" | cut -d. -f1-2)"

go mod edit -module "$NEW_MODULE"
go mod edit -go "$GO_VERSION_MAJOR_MINOR"

# Rewrite import paths in source and docs.
while IFS= read -r -d '' file; do
  if grep -q "$TEMPLATE_MODULE" "$file" 2>/dev/null; then
    if sed --version >/dev/null 2>&1; then
      sed -i "s|${TEMPLATE_MODULE}|${NEW_MODULE}|g" "$file"
    else
      sed -i '' "s|${TEMPLATE_MODULE}|${NEW_MODULE}|g" "$file"
    fi
  fi
done < <(find . -type f \( -name '*.go' -o -name 'go.mod' -o -name '*.md' -o -name '*.yaml' -o -name '*.yml' -o -name 'Makefile' -o -name 'Dockerfile' \) -print0)

# Stamp AGENTS.md profile.
AGENTS_FILE="AGENTS.md"
if [[ -f "$AGENTS_FILE" ]]; then
  if sed --version >/dev/null 2>&1; then
    SED_INPLACE=(sed -i)
  else
    SED_INPLACE=(sed -i '')
  fi
  "${SED_INPLACE[@]}" "s|{{PROJECT_NAME}}|${PROJECT_NAME}|g" "$AGENTS_FILE"
  "${SED_INPLACE[@]}" "s|{{MODULE_PATH}}|${NEW_MODULE}|g" "$AGENTS_FILE"
  "${SED_INPLACE[@]}" "s|{{APP_NAME}}|${APP_NAME}|g" "$AGENTS_FILE"
  "${SED_INPLACE[@]}" "s|{{HTTP_PORT}}|${HTTP_PORT}|g" "$AGENTS_FILE"
  "${SED_INPLACE[@]}" "s|{{DATABASE}}|${DATABASE_CHOICE}|g" "$AGENTS_FILE"
  "${SED_INPLACE[@]}" "s|{{GO_VERSION}}|${GO_VERSION_MAJOR_MINOR}|g" "$AGENTS_FILE"
fi

# OpenAPI title
if [[ -f api/openapi.yaml ]]; then
  if sed --version >/dev/null 2>&1; then
    sed -i "s|^  title:.*|  title: ${APP_NAME} API|" api/openapi.yaml
  else
    sed -i '' "s|^  title:.*|  title: ${APP_NAME} API|" api/openapi.yaml
  fi
fi

DEFAULT_DB_URL="postgres://postgres:postgres@localhost:5432/${PROJECT_NAME//-/_}?sslmode=disable"
ENV_FILE=".env"
ENV_EXAMPLE=".env.example"

cat > "$ENV_EXAMPLE" <<EOF
APP_NAME=${APP_NAME}
APP_ENV=local
HTTP_PORT=${HTTP_PORT}
HTTP_ALLOWED_ORIGINS=${HTTP_ALLOWED_ORIGINS}
HTTP_SHUTDOWN_TIMEOUT=10s
HTTP_READ_TIMEOUT=15s
HTTP_WRITE_TIMEOUT=15s
HTTP_IDLE_TIMEOUT=60s
LOG_LEVEL=debug
LOG_ADD_SOURCE=false
EOF

if [[ "$DATABASE_CHOICE" == postgres ]]; then
  cat >> "$ENV_EXAMPLE" <<EOF

DATABASE_URL=${DEFAULT_DB_URL}
DATABASE_REQUIRED=false
DATABASE_MAX_OPEN_CONNS=10
DATABASE_MAX_IDLE_CONNS=5
DATABASE_CONN_MAX_LIFETIME=1h
EOF
else
  cat >> "$ENV_EXAMPLE" <<'EOF'

# Optional. Leave empty to use the local welcome repository.
# DATABASE_URL=postgres://postgres:postgres@localhost:5432/app_db?sslmode=disable
# DATABASE_REQUIRED=false
EOF
fi

cp "$ENV_EXAMPLE" "$ENV_FILE"

mkdir -p "$TARGET_DIR"
if command -v rsync >/dev/null 2>&1; then
  rsync -a --exclude '.git' "$WORK_DIR"/ "$TARGET_DIR"/
else
  cp -R "$WORK_DIR"/. "$TARGET_DIR"/
  rm -rf "$TARGET_DIR/.git"
fi

# Keep the generator in the template repo only; generated apps do not need it.
rm -f "$TARGET_DIR/scripts/go-clean-init.sh"

cd "$TARGET_DIR"
go mod tidy
gofmt -w $(find . -name '*.go')
go test ./...

printf '\n%bProject initialized%b\n' "$GREEN$BOLD" "$RESET"
printf '  Module : %s\n' "$NEW_MODULE"
printf '  Path   : %s\n\n' "$TARGET_DIR"
cat <<EOF
Next steps:
  1) cd ${PROJECT_NAME}
  2) Review .env
  3) make run
  4) Visit http://localhost:${HTTP_PORT}/ and /docs
EOF
