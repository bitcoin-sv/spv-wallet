#!/bin/bash

# Define color codes
color_success="\033[0;32m"  # green
color_danger="\033[0;31m"   # red
color_warning="\033[0;33m"  # yellow
color_debug="\033[0;34m"    # blue
color_user="\033[0;35m"  # purple

# Reset color code
color_reset="\033[0m"
choice=''

function print_debug() {
  if [ "$debug" == "true" ]; then
      echo -e "${color_debug}$1${color_reset}"
  fi
}

function print_success() {
  echo -e "${color_success}$1${color_reset}"
}

function print_warning() {
  echo -e "${color_warning}$1${color_reset}"
}

function print_error() {
  echo -e "${color_danger}$1${color_reset}"
}

function printPrompt() {
  echo -e "${color_user}$1${color_reset}"
}

function ask_for_value() {
    local prompt="$1"
    local prefix="$2"

    printPrompt "$prompt"
    read -p ">" choice

    if [[ -z "$choice" ]]; then
        return
    fi

    while ! [[ "$choice" =~ ^"$prefix".*$ ]]; do
      echo -e "${color_danger}Invalid value!${color_reset}"
      read -p "> " choice
      if [[ -z "$choice" ]]; then
          return
      fi
    done
}

function ask_for_choice() {
    local prompt="$1"
    local options=("${@:2}")

    printPrompt "$prompt"
    for (( i = 0; i < ${#options[@]}; i++ )); do
        echo "$((i+1)). ${options[i]}"
    done


    read -p ">" choice

    while ! [[ "$choice" =~ ^[0-9]+$ ]] || (( choice < 1 || choice > ${#options[@]} )); do
        echo -e "${color_danger}Invalid choice! Please select a valid option.${color_reset}"
        read -p "> " choice
    done
}

function ask_for_yes_or_no() {
    local prompt="$1"
    local default_value="${2:-true}"

    local default_prompt="[Y/n]"
    local inverse_default_prompt="[y/N]"

    if [[ "$default_value" == "true" ]]; then
        prompt="$prompt $default_prompt"
    elif [[ "$default_value" == "false" ]]; then
        prompt="$prompt $inverse_default_prompt"
    fi

    printPrompt "$prompt"

    local response
    read -p ">" response

    if [[ -z "$response" ]]; then
        choice="$default_value"
        return
    fi

    while ! [[ "$response" =~ ^(yes|no|y|n)$ ]]; do
        echo -e "${color_danger}Invalid response! Please enter 'yes' or 'no'.${color_reset}"
        read -p "> " response
        if [[ -z "$response" ]]; then
            choice="$default_value"
            return
        fi
    done

    if [[ "$response" =~ ^(yes|y)$ ]]; then
        choice="true"
    else
        choice="false"
    fi
}

function print_state() {
    print_debug "State:"
    print_debug "  database=${database}"
    print_debug "  cache=${cache}"
    print_debug "  bux_server=${bux_server}"
    print_debug "  bux_wallet_frontend=${bux_wallet_frontend}"
    print_debug "  bux_wallet_backend=${bux_wallet_backend}"
    print_debug "  background=${background}"
    print_debug "  default_xpub: $default_xpub"
    print_debug "  admin_xpub=${admin_xpub}"
    print_debug "  admin_xpriv=${admin_xpriv}"
    print_debug "  load_config=${load_config}"
    print_debug "  no_rebuild=${no_rebuild}"
    print_debug "  debug=${debug}"
    print_debug ""
}

function load_from() {
    local key="$1"
    local variable="$2"

    if [[ -n "${!variable}" ]]; then
        return
    fi

    if [[ "$line" =~ ^("$key"=) ]]; then
        print_debug "Loading $key to variable $variable"
        value="${line#*=}"
        value="${value//\"}"
        print_debug "Value for $variable is '$value'"
        eval "$variable=\"$value\""
    fi
}

function save_to() {
  local key="$1"
  local variable="$2"

  if [[ "${!variable}" == "" ]]; then
      return
  fi

  save_value "$key" "${!variable}"
}

function save_value() {
  local key="$1"
  local value="$2"
  echo "$key=\"${value}\"" >> .env.config
}

function docker_compose_up() {
  local compose_plugin=false
  if command docker compose version &> /dev/null; then
  	compose_plugin=true
  fi

  local run="true"

  if [ "$debug" == 'true' ]; then
    echo ""
    if [ $compose_plugin == true ]; then
      print_debug "docker compose up $1"
    else
      print_debug "docker-compose up $1"
    fi
    echo ""
    ask_for_yes_or_no "You use debug mode. Do you want to run docker compose now?"
    run=$choice
  fi
  if [ "$run" != "true" ]; then
    return
  fi

  if [ $compose_plugin == true ]; then
    docker compose up $1
  else
    docker-compose up $1
  fi
}

function prefix_each() {
    local delimiter="$1"
    local result=""
    shift
    for element in "$@"; do
       result+="$delimiter $element "
    done
    echo "$result"
}

# === LOAD FROM CLI ===

while [[ $# -gt 0 ]]; do
    key="$1"

    case $key in
        -db|--database)
        database="$2"
        shift
        ;;
        -c|--cache)
        cache="$2"
        shift
        ;;
        -bs|--bux-server)
        bux_server="$2"
        shift
        ;;
        -pl|--pulse)
        pulse="$2"
        shift
        ;;
        -bwf|--bux-wallet-frontend)
        bux_wallet_frontend="$2"
        shift
        ;;
        -bwb|--bux-wallet-backend)
        bux_wallet_backend="$2"
        shift
        ;;
        --xpub)
        admin_xpub="$2"
        shift
        ;;
        --xprv)
        admin_xpriv="$2"
        shift
        ;;
        -pm|--paymail)
        paymail_domain="$2"
        shift
        ;;
        -l|--load)
        load_config="true"
        # no additional arguments so now `shift` command
        ;;
        -nrb|--no-rebuild)
        no_rebuild="true"
        # no additional arguments so now `shift` command
        ;;
        -b|--background)
        background="true"
        # no additional arguments so now `shift` command
        ;;
        -d|--debug)
        debug="true"
        # no additional arguments so now `shift` command
        ;;
        -h|--help)
        echo -e "Usage: ./start.sh [OPTIONS]"
        echo ""
        echo "This script helps you to run SPV Wallet with your preferred database and cache storage."
        echo ""
        echo -e "Options:"
        echo -e "  -pm,  --paymail\t\t PayMail domain for which to run all applications"
        echo -e "  -l,   --load\t\t\t Load previously stored config from .env.config file"
        echo -e "  -nrb, --no-rebuild\t\t Prevent rebuild of docker images before running"
        echo -e "  -b,   --background\t\t Whether the applications should be run in background"
        echo -e "  -d,   --debug\t\t\t Run in debug mode"
        echo -e "  -h,   --help\t\t\t Show this message"
        echo -e ""
        echo -e "<----------   BUX SERVER SECTION"
        echo -e "  -bs,  --bux-server\t\t Whether the bux-server should be run - true/false"
        echo -e "  -db,  --database\t\t Define database - postgresql, mongodb, sqlite"
        echo -e "  -c,   --cache\t\t\t Define cache storage - freecache(in-memory), redis"
        echo -e "  --xpub\t\t\t Define admin xPub"
        echo ""
        echo -e "<----------   BUX WALLET SECTION"
        echo -e "  -bwf,  --bux-wallet-frontend\t Whether the bux-wallet-frontend should be run - true/false"
        echo -e "  -bwb,  --bux-wallet-backend\t Whether the bux-wallet-backend should be run - true/false"
        echo -e "  --xprv\t\t\t Define admin xPriv"
        exit 1;
        shift
        ;;
        *)
        ;;
    esac
    shift
done

# Welcome message
echo -e "${color_user}Welcome in SPV Wallet!${color_reset}"

print_debug "Loaded config from CLI:"
print_state


# === LOAD FROM FILE ===
if [ "$load_config" == "true" ]; then
    if [ -f .env.config ]; then
        print_debug "Loading config from .env.config"

        while IFS= read -r line; do
            print_debug "Checking line '$line'"
            load_from 'BUX_DB_DATASTORE_ENGINE' database
            load_from 'BUX_CACHE_ENGINE' cache
            load_from 'RUN_BUX_SERVER' bux_server
            load_from 'RUN_PULSE' pulse
            load_from 'RUN_BUX_WALLET_FRONTEND' bux_wallet_frontend
            load_from 'RUN_BUX_WALLET_BACKEND' bux_wallet_backend
            load_from 'RUN_PAYMAIL_DOMAIN' paymail_domain
            load_from 'RUN_IN_BACKGROUND' background
            load_from 'RUN_WITH_DEFAULT_XPUB' default_xpub
            load_from 'BUX_AUTH_ADMIN_KEY' admin_xpub
            load_from 'BUX_ADMIN_XPRIV' admin_xpriv
            load_from 'RUN_WITHOUT_REBUILD' no_rebuild
        done < ".env.config"
        print_success "Config loaded from .env.config file"
        print_debug "Config after loading .env.config:"
        print_state
    else
        print_warning "File .env.config does not exist, but you choose to load from it."
    fi
fi

# === COLLECT CONFIG FROM USER IF NEEDED ===

# <----------   BUX SERVER SECTION
if [ "$database" == "" ]; then
    database_options=("postgresql" "mongodb" "sqlite")
    ask_for_choice "Select your database:" "${database_options[@]}"

    case $choice in
        1) database="postgresql";;
        2) database="mongodb";;
        3) database="sqlite";;
    esac
    print_debug "database: $database"
fi

if [ "$cache" == "" ]; then
    cache_options=("freecache" "redis")
    ask_for_choice "Select your cache storage:" "${cache_options[@]}"

    case $choice in
        1) cache="freecache";;
        2) cache="redis";;
    esac
    print_debug "cache: $cache"
fi

if [ "$bux_server" == "" ]; then
    ask_for_yes_or_no "Do you want to run Bux-server?"
    bux_server="$choice"
    print_debug "bux_server: $bux_server"
fi

if [ "$pulse" == "" ]; then
    ask_for_yes_or_no "Do you want to run Pulse?"
    pulse="$choice"
    print_debug "pulse: $pulse"
fi

if [ "$bux_wallet_frontend" == "" ]; then
    ask_for_yes_or_no "Do you want to run bux-wallet-frontend?"
    bux_wallet_frontend="$choice"
    print_debug "bux_wallet_frontend: $bux_wallet_frontend"
fi

if [ "$bux_wallet_backend" == "" ]; then
    ask_for_yes_or_no "Do you want to run bux-wallet-backend?"
    bux_wallet_backend="$choice"
    print_debug "bux_wallet_backend: $bux_wallet_backend"
fi

if [ "$bux_server" == "true" ] && [ "$admin_xpub" == "" ] && [ "$default_xpub" != "true" ]; then
    ask_for_value "Define admin xPub (Leave empty to use the default one)" 'xpub'

    if [[ -n "$choice" ]]; then
        admin_xpub=$choice
    else
        default_xpub="true"
    fi
    print_debug "admin_xpub: $admin_xpub"
    print_debug "default_xpub: $default_xpub"
fi

if [ "$bux_wallet_backend" == "true" ] && [ "$admin_xpriv" == "" ] && [ "$default_xpub" != "true" ]; then
  ask_for_value "Define admin xPriv (Leave empty to use the default one)" 'xprv'
  admin_xpriv=$choice
  print_debug "admin_xpriv: $admin_xpriv"
fi

if [ "$paymail_domain" == "" ] && { [ "$bux_wallet_backend" == "true" ] || [ "$bux_wallet_frontend" == "true" ] || [ "$bux_server" == "true" ]; }; then
    ask_for_value "What PayMail domain should be configured in applications?"
    paymail_domain=$choice
    print_debug "paymail_domain: $paymail_domain"
fi

if [ "$background" == "" ]; then
    ask_for_yes_or_no "Do you want to run everything in the background?" "false"
    background="$choice"
    print_debug "background: $background"
fi

# === SAVE CONFIG ===
print_debug "Config before storing:"
print_state

# Create the .env.config file
print_debug "Creating/Cleaning .env.config file."
echo "# Used by start.sh. All unknown variables will be removed after running the script" > ".env.config"
save_to 'BUX_DB_DATASTORE_ENGINE' database
save_to 'BUX_CACHE_ENGINE' cache
save_to 'RUN_BUX_SERVER' bux_server
save_to 'RUN_PULSE' pulse
save_to 'RUN_PAYMAIL_DOMAIN' paymail_domain
save_to 'RUN_BUX_WALLET_FRONTEND' bux_wallet_frontend
save_to 'RUN_BUX_WALLET_BACKEND' bux_wallet_backend
save_to 'RUN_IN_BACKGROUND' background
save_to 'RUN_WITH_DEFAULT_XPUB' default_xpub
save_to 'BUX_AUTH_ADMIN_KEY' admin_xpub
save_to 'BUX_ADMIN_XPRIV' admin_xpriv
save_to 'RUN_WITHOUT_REBUILD' no_rebuild
case $database in
  postgresql)
    save_value 'BUX_DB_SQL_HOST' "bux-postgresql"
    save_value 'BUX_DB_SQL_NAME' "postgres"
    save_value 'BUX_DB_SQL_USER' "postgres"
    save_value 'BUX_DB_SQL_PASSWORD' "postgres"
  ;;
  mongodb)
    save_value 'BUX_DB_MONGODB_URI' "mongodb://mongo:mongo@bux-mongodb:27017/"
  ;;
esac

if [ "$cache" == "redis" ]; then
  save_value 'BUX_CACHE_REDIS_URL' "redis://redis:6379"
fi

if [ "$bux_server" == "true" ]; then
  save_value 'BUX_SERVER_URL' "http://bux-server:3003/v1"
else
  save_value 'BUX_SERVER_URL' "http://host.docker.internal:3003/v1"
fi
if [ "$bux_wallet_backend" == "true" ]; then
  save_value 'DB_HOST' "bux-postgresql"
fi

if [ "$pulse" == "true" ]; then
  save_value 'BUX_PAYMAIL_BEEF_PULSE_URL' "http://pulse:8080/api/v1/chain/merkleroot/verify"
else
  save_value 'BUX_PAYMAIL_BEEF_PULSE_URL' "http://host.docker.internal:8080/api/v1/chain/merkleroot/verify"
fi
print_debug "Exporting RUN_PAYMAIL_DOMAIN environment variable"
export RUN_PAYMAIL_DOMAIN="$paymail_domain"

print_success "File .env.config updated!"
print_debug "$(cat .env.config)"

# === RUN WHAT IS NEEDED ===
servicesToRun=()
servicesToHideLogs=()
additionalFlags=()

case $database in
  postgresql)
    servicesToRun+=("bux-postgresql")
    servicesToHideLogs+=("bux-postgresql")
  ;;
  mongodb)
    servicesToRun+=("bux-mongodb")
    servicesToHideLogs+=("bux-mongodb")
  ;;
esac

if [ "$cache" == "redis" ]; then
  servicesToRun+=("bux-redis")
  servicesToHideLogs+=("bux-redis")
fi

if [ "$bux_server" == "true" ]; then
  servicesToRun+=("bux-server")
fi

if [ "$pulse" == "true" ]; then
  servicesToRun+=("pulse")
fi

if [ "$bux_wallet_backend" == "true" ]; then
  servicesToRun+=("bux-wallet-backend")
  servicesToRun+=("bux-postgresql")
  servicesToHideLogs+=("bux-postgresql")
fi

if [ "$bux_wallet_frontend" == "true" ]; then
  servicesToRun+=("bux-wallet-frontend")
  servicesToHideLogs+=("bux-wallet-frontend")
fi

if [ "$no_rebuild" != "true" ]; then
  additionalFlags+=("--build")
fi

if [ "$background" == "true" ]; then
  additionalFlags+=("-d")
fi

docker_compose_up  "${servicesToRun[*]} ${additionalFlags[*]} $(prefix_each '--no-attach ' ${servicesToHideLogs[*]})"
