#!/bin/bash
# =============================================================================
# Sub2API Docker Deployment Preparation Script
# =============================================================================
# This script prepares deployment files for Sub2API:
#   - Downloads .env.example
#   - Generates secure secrets (JWT_SECRET, TOTP_ENCRYPTION_KEY, POSTGRES_PASSWORD)
#   - Generates a random admin login email (ADMIN_EMAIL)
#   - Creates necessary data directories
#
# After running this script, you can start services with:
#   docker-compose up -d
#
# IMPORTANT — re-running this script is safe:
#   PostgreSQL runs initdb only when the data directory is empty, so
#   POSTGRES_PASSWORD is consumed exactly once. Re-generating it on an existing
#   deployment produces a .env whose password the database has never seen, and
#   the app then fails with "password authentication failed for user ...".
#   This script therefore never re-generates secrets for an existing deployment;
#   it reuses the values already present in .env.
# =============================================================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# GitHub raw content base URL (override with: GITHUB_RAW_URL=... ./docker-deploy.sh)
GITHUB_RAW_URL="${GITHUB_RAW_URL:-https://raw.githubusercontent.com/Wei-Shaw/sub2api/main/deploy}"

# PostgreSQL data directory — must match the bind mount in docker-compose*.yml
PGDATA_DIR="${PGDATA_DIR:-./postgres_data}"

# Comma-separated container name hints used when a password sync is required
PG_CONTAINER_HINT="${PG_CONTAINER_HINT:-sub2api-postgres-dev,sub2api-postgres,tokenpilot-postgres}"

FORCE_SECRETS=false
for arg in "$@"; do
    case "$arg" in
        --force-secrets) FORCE_SECRETS=true ;;
        -h|--help)
            echo "Usage: bash docker-deploy.sh [--force-secrets]"
            echo ""
            echo "  --force-secrets  仅在没有数据库数据目录时可用；显式要求重新生成密钥。"
            echo "                   已存在 ${PGDATA_DIR} 时该参数会被拒绝。"
            exit 0
            ;;
        *) ;;
    esac
done

# Print colored message
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Generate random secret
generate_secret() {
    openssl rand -hex 32
}

# Generate a random admin login email so new installs never use a guessable default
generate_admin_email() {
    local suffix
    suffix=$(openssl rand -hex 6) || return 1
    [ -n "$suffix" ] || return 1
    echo "admin-${suffix}@sub2api.local"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Read a KEY=value from the existing .env (last occurrence wins, like Compose)
env_value() {
    local key="$1"
    [ -f .env ] || return 0
    grep -E "^[[:space:]]*${key}=" .env 2>/dev/null | tail -n 1 | cut -d= -f2-
}

# Show a secret without leaking it entirely
mask_secret() {
    local v="$1"
    if [ -z "$v" ]; then
        echo "(empty)"
    elif [ "${#v}" -le 12 ]; then
        echo "********"
    else
        echo "${v:0:8}...${v: -4}"
    fi
}

# True when the PostgreSQL cluster has already been initialised
has_cluster() {
    [ -s "${PGDATA_DIR}/PG_VERSION" ]
}

# Append keys from .env.example that .env does not define yet. Existing values
# are never touched — that is the whole point of this script.
# Prints the number of appended keys on stdout; per-key notes go to stderr.
append_missing_keys() {
    local added=0 key line
    [ -f .env.example ] || { echo 0; return 0; }
    while IFS= read -r line; do
        case "$line" in
            ''|'#'*) continue ;;
        esac
        key="${line%%=*}"
        key="${key#export }"
        [ -n "$key" ] || continue
        if ! grep -qE "^[[:space:]]*${key}=" .env 2>/dev/null; then
            printf '%s\n' "$line" >> .env
            echo "  新增配置项: ${key}" >&2
            added=$((added + 1))
        fi
    done < .env.example
    echo "$added"
}

# Download .env.example, falling back to the local copy when the network fails
download_env_example() {
    print_info "Downloading .env.example..."
    local ok=false
    if command_exists curl; then
        if curl -fsSL "${GITHUB_RAW_URL}/.env.example" -o .env.example.part 2>/dev/null; then
            ok=true
        fi
    elif command_exists wget; then
        if wget -q "${GITHUB_RAW_URL}/.env.example" -O .env.example.part; then
            ok=true
        fi
    fi

    if [ "$ok" = true ] && [ -s .env.example.part ]; then
        mv .env.example.part .env.example
        print_success "Downloaded .env.example"
    else
        rm -f .env.example.part
        if [ -s .env.example ]; then
            print_warning "下载失败，继续使用本地已有的 .env.example"
        else
            print_error "下载 .env.example 失败，且本地不存在该文件。"
            print_error "请检查网络，或先手动准备 .env.example 后重试。"
            exit 1
        fi
    fi
}

# Find a running PostgreSQL container that belongs to this deployment
find_pg_container() {
    command_exists docker || return 0
    docker ps --format '{{.Names}} {{.Image}}' 2>/dev/null \
        | awk '$2 ~ /postgres/ {print $1}' \
        | grep -E 'sub2api|tokenpilot' \
        | head -n 1
}

# Force the live database to accept the password stored in .env.
# Inside the container the unix socket matches `local all all trust`, so no
# password is required to run ALTER USER.
sync_db_password() {
    local user="$1" dbname="$2" password="$3"
    local container pw_esc out

    container="$(find_pg_container)"
    if [ -z "$container" ]; then
        print_warning "未找到运行中的 PostgreSQL 容器，无法自动同步密码。"
        echo ""
        echo "  请先启动容器，然后执行（容器内 local socket 为 trust，无需密码）："
        echo ""
        echo "  for c in $(echo "${PG_CONTAINER_HINT}" | tr ',' ' '); do"
        echo "    docker ps --format '{{.Names}}' | grep -qx \"\$c\" || continue"
        echo "    docker exec -it \"\$c\" psql -U ${user} -d ${dbname} \\"
        echo "      -c \"ALTER USER ${user} WITH PASSWORD '${password}';\""
        echo "  done"
        echo ""
        return 1
    fi

    pw_esc="${password//\'/\'\'}"
    if out=$(docker exec "$container" psql -U "$user" -d "$dbname" -v ON_ERROR_STOP=1 \
             -c "ALTER USER ${user} WITH PASSWORD '${pw_esc}';" 2>&1); then
        print_success "已将数据库密码同步为 .env 中的值（容器: ${container}）"
        return 0
    fi
    # The application database may not exist yet; retry against the maintenance DB
    if out=$(docker exec "$container" psql -U "$user" -d postgres -v ON_ERROR_STOP=1 \
             -c "ALTER USER ${user} WITH PASSWORD '${pw_esc}';" 2>&1); then
        print_success "已将数据库密码同步为 .env 中的值（容器: ${container}, 库: postgres）"
        return 0
    fi

    print_warning "自动同步数据库密码失败：${out}"
    echo ""
    echo "  请手动执行："
    echo "  docker exec -it ${container} psql -U ${user} -d ${dbname} \\"
    echo "    -c \"ALTER USER ${user} WITH PASSWORD '${password}';\""
    echo ""
    return 1
}

# Main installation function
main() {
    echo ""
    echo "=========================================="
    echo "  Sub2API Deployment Preparation"
    echo "=========================================="
    echo ""

    # Check if openssl is available
    if ! command_exists openssl; then
        print_error "openssl is not installed. Please install openssl first."
        exit 1
    fi

    local cluster_exists=false
    local env_exists=false
    local reuse_mode=false

    has_cluster && cluster_exists=true
    [ -f .env ] && env_exists=true

    # ---------------------------------------------------------------------
    # Existing deployment: never re-generate secrets
    # ---------------------------------------------------------------------
    if [ "$cluster_exists" = true ] || [ "$env_exists" = true ]; then
        reuse_mode=true

        if [ "$cluster_exists" = true ]; then
            print_warning "检测到已有数据库：${PGDATA_DIR}/PG_VERSION 存在"
        else
            print_warning "检测到已有 .env"
        fi
        print_info "本次运行不会重新生成 POSTGRES_PASSWORD / JWT_SECRET / TOTP_ENCRYPTION_KEY。"
        echo ""

        if [ "$FORCE_SECRETS" = true ] && [ "$cluster_exists" = true ]; then
            print_error "拒绝执行 --force-secrets：${PGDATA_DIR} 中已存在数据库。"
            print_error "重新生成 POSTGRES_PASSWORD 会导致应用永远无法连接（initdb 不会重跑）。"
            echo ""
            echo "  · 只需换数据库密码：先 ALTER USER 再改 .env，或直接重新运行本脚本"
            echo "  · 想彻底重装：停止容器并删除 ${PGDATA_DIR}（会丢失全部数据）后重跑"
            echo ""
            exit 1
        fi

        local old_pg old_jwt old_totp old_admin
        old_pg="$(env_value POSTGRES_PASSWORD)"
        old_jwt="$(env_value JWT_SECRET)"
        old_totp="$(env_value TOTP_ENCRYPTION_KEY)"
        old_admin="$(env_value ADMIN_EMAIL)"

        if [ -n "$old_pg" ] && [ "$FORCE_SECRETS" = false ]; then
            POSTGRES_PASSWORD="$old_pg"
        else
            POSTGRES_PASSWORD="$(generate_secret)"
        fi
        JWT_SECRET="${old_jwt:-$(generate_secret)}"
        TOTP_ENCRYPTION_KEY="${old_totp:-$(generate_secret)}"
        ADMIN_EMAIL="${old_admin:-$(generate_admin_email)}"

        # A fresh password is only possible when the old one is unknown; the
        # database must then be told about it explicitly.
        if [ "$cluster_exists" = true ] && [ "$POSTGRES_PASSWORD" != "$old_pg" ]; then
            NEED_DB_SYNC=true
        fi

        download_env_example

        if [ "$env_exists" = true ]; then
            local backup=".env.bak.$(date +%Y%m%d-%H%M%S)"
            cp -p .env "$backup"
            print_info "已备份原 .env → ${backup}"

            local added
            added="$(append_missing_keys)"
            print_info "从新模板补齐 ${added} 个缺失配置项（已有值一律保持不变）"

            # .env existed but carried no POSTGRES_PASSWORD: persist the new one
            if [ -z "$old_pg" ]; then
                printf 'POSTGRES_PASSWORD=%s\n' "$POSTGRES_PASSWORD" >> .env
                print_warning ".env 中原本没有 POSTGRES_PASSWORD，已写入新生成的值"
            fi
        else
            # Database exists but .env is gone: the original plaintext password
            # is unrecoverable (only a SCRAM hash lives in pg_authid), so create
            # a fresh .env and force the database to accept it.
            print_warning ".env 已丢失，无法还原原密码（pg_authid 只存单向哈希）。"
            print_warning "将生成新密码，并尝试把数据库侧密码同步为同一个值。"
            cp .env.example .env
            sed -i "s/^JWT_SECRET=.*/JWT_SECRET=${JWT_SECRET}/" .env
            sed -i "s/^TOTP_ENCRYPTION_KEY=.*/TOTP_ENCRYPTION_KEY=${TOTP_ENCRYPTION_KEY}/" .env
            sed -i "s/^POSTGRES_PASSWORD=.*/POSTGRES_PASSWORD=${POSTGRES_PASSWORD}/" .env
            sed -i "s/^ADMIN_EMAIL=.*/ADMIN_EMAIL=${ADMIN_EMAIL}/" .env
            NEED_DB_SYNC=true
        fi
    else
        # -----------------------------------------------------------------
        # Fresh installation
        # -----------------------------------------------------------------
        download_env_example

        print_info "Generating secure secrets..."
        echo ""

        JWT_SECRET=$(generate_secret)
        TOTP_ENCRYPTION_KEY=$(generate_secret)
        POSTGRES_PASSWORD=$(generate_secret)
        ADMIN_EMAIL=$(generate_admin_email)

        cp .env.example .env

        # Update .env with generated secrets (cross-platform compatible)
        if sed --version >/dev/null 2>&1; then
            # GNU sed (Linux)
            sed -i "s/^JWT_SECRET=.*/JWT_SECRET=${JWT_SECRET}/" .env
            sed -i "s/^TOTP_ENCRYPTION_KEY=.*/TOTP_ENCRYPTION_KEY=${TOTP_ENCRYPTION_KEY}/" .env
            sed -i "s/^POSTGRES_PASSWORD=.*/POSTGRES_PASSWORD=${POSTGRES_PASSWORD}/" .env
            sed -i "s/^ADMIN_EMAIL=.*/ADMIN_EMAIL=${ADMIN_EMAIL}/" .env
        else
            # BSD sed (macOS)
            sed -i '' "s/^JWT_SECRET=.*/JWT_SECRET=${JWT_SECRET}/" .env
            sed -i '' "s/^TOTP_ENCRYPTION_KEY=.*/TOTP_ENCRYPTION_KEY=${TOTP_ENCRYPTION_KEY}/" .env
            sed -i '' "s/^POSTGRES_PASSWORD=.*/POSTGRES_PASSWORD=${POSTGRES_PASSWORD}/" .env
            sed -i '' "s/^ADMIN_EMAIL=.*/ADMIN_EMAIL=${ADMIN_EMAIL}/" .env
        fi
    fi

    # Create data directories
    print_info "Creating data directories..."
    mkdir -p data postgres_data redis_data
    print_success "Created data directories"

    # Set secure permissions for .env file (readable/writable only by owner)
    chmod 600 .env
    echo ""

    # Sync the database password when the plaintext is no longer known
    if [ "${NEED_DB_SYNC:-false}" = true ]; then
        local db_user db_name
        db_user="$(env_value POSTGRES_USER)"
        db_user="${db_user:-sub2api}"
        db_name="$(env_value POSTGRES_DB)"
        db_name="${db_name:-sub2api}"
        print_info "尝试同步数据库密码..."
        sync_db_password "$db_user" "$db_name" "$POSTGRES_PASSWORD" || true
        echo ""
    fi

    # Display completion message
    echo "=========================================="
    echo "  Preparation Complete!"
    echo "=========================================="
    echo ""

    if [ "$reuse_mode" = true ]; then
        echo "Existing deployment detected — secrets preserved:"
    else
        echo "Generated secure credentials:"
    fi
    echo "  POSTGRES_PASSWORD:     ${POSTGRES_PASSWORD}"
    echo "  JWT_SECRET:            ${JWT_SECRET}"
    echo "  TOTP_ENCRYPTION_KEY:   ${TOTP_ENCRYPTION_KEY}"
    echo "  ADMIN_EMAIL:           ${ADMIN_EMAIL}  (admin login username)"
    echo ""
    print_warning "These credentials have been saved to .env file."
    print_warning "Please keep them secure and do not share publicly!"
    echo ""

    if [ "$reuse_mode" = true ]; then
        print_info "POSTGRES_PASSWORD 摘要: $(mask_secret "$POSTGRES_PASSWORD")（未改动，数据库仍可连接）"
        echo ""
        print_warning "重要：不要再手动修改 .env 里的 POSTGRES_PASSWORD。"
        print_warning "PostgreSQL 只在数据目录为空时执行 initdb，改 .env 不会改数据库。"
        echo ""
        echo "  验证密码是否与数据库一致（必须 -h 127.0.0.1 强制走 TCP）："
        echo "    PW=\$(grep '^POSTGRES_PASSWORD=' .env | cut -d= -f2-)"
        echo "    docker exec -e PGPASSWORD=\"\$PW\" \$(docker ps --format '{{.Names}}' | grep -E 'postgres' | head -1) \\"
        echo "      psql -h 127.0.0.1 -U sub2api -d sub2api -tAc \"select 'PASSWORD OK'\""
        echo ""
        echo "  确实需要换密码时，两步一起做："
        echo "    1) ALTER USER sub2api WITH PASSWORD '<新密码>';   （在容器内执行）"
        echo "    2) 再把 .env 里的 POSTGRES_PASSWORD 改成同一个值"
        echo ""
    fi

    echo "Directory structure:"
    echo "  docker-compose.yml        - Docker Compose configuration"
    echo "  .env                      - Environment variables (generated secrets)"
    echo "  .env.example              - Example template (for reference)"
    echo "  data/                     - Application data (will be created on first run)"
    echo "  postgres_data/            - PostgreSQL data"
    echo "  redis_data/               - Redis data"
    echo ""
    echo "Next steps:"
    echo "  1. (Optional) Edit .env to customize configuration"
    echo "  2. Start services:"
    echo "     docker-compose up -d"
    echo ""
    echo "  3. View logs:"
    echo "     docker-compose logs -f sub2api"
    echo ""
    echo "  4. Access Web UI:"
    echo "     http://localhost:8080"
    echo ""
    print_info "If admin password is not set in .env, it will be auto-generated."
    print_info "Check logs for the generated admin password on first startup."
    echo ""
}

# Run main function
main "$@"