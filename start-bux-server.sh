#!/bin/bash

reset='\033[0m'
choice=''

function ask_for_choice() {
    local prompt="$1"
    local options=("${@:2}")

    echo -e "$prompt"
    for (( i = 0; i < ${#options[@]}; i++ )); do
        echo "$((i+1)). ${options[i]}"
    done


    read -p "> " choice

    while ! [[ "$choice" =~ ^[0-9]+$ ]] || (( choice < 1 || choice > ${#options[@]} )); do
        echo -e "\033[0;31mInvalid choice! Please select a valid option.$reset"
        read -p "> " choice
    done
}

function ask_for_yes_or_no() {
    local prompt="$1"
    echo -e "$prompt"

    local response
    read -p "(y/n)> " response

    while ! [[ "$response" =~ ^(yes|no|y|n)$ ]]; do
        echo -e "\033[0;31mInvalid response! Please enter 'yes' or 'no'.$reset"
        read -p "> " response
    done

    if [[ "$response" =~ ^(yes|y)$ ]]; then
        choice="true"
    else
        choice="false"
    fi
}

compose_plugin=false
if command docker compose version &> /dev/null; then
	compose_plugin=true
fi

function docker_compose_up() {
    if [ compose_plugin == true ]; then
        docker compose up $1
    else
        docker-compose up $1
    fi
}

# Welcome message
echo -e "\033[0;33m\033[1mWelcome in Bux Server!$reset"

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
        -b|--background)
        background="$2"
        shift
        ;;
        -x|--xpub)
        admin_xpub="$2"
        shift
        ;;
        -l|--load)
        load_config="true"
        shift
        ;;
        -rb|--rebuild)
        rebuild="true"
        load_config="true"
        ;;
        -h|--help)
        echo -e "\033[1mUsage: ./start-bux-server.sh [OPTIONS]$reset"
        echo ""
        echo "This script helps you to run Bux server with your preferred database and cache storage."
        echo ""
        echo -e "Options:$reset"
        echo -e "  -db,  --database\t Define database - postgresql, mongodb, sqlite$reset"
        echo -e "  -c,   --cache\t\t Define cache storage - freecache(in-memory), redis$reset"
        echo -e "  -bs,  --bux-server\t Whether the bux-server should be run - true/false$reset"
        echo -e "  -b,   --background\t Whether the bux-server should be run in background - true/false$reset"
        echo -e "  -x,   --xpub\t\t Define admin xPub$reset"
        echo -e "  -l,   --load\t\t Load .env.config file and run bux-server with its settings$reset"
        echo -e "  -rb,  --rebuild\t Rebuild docker images before running$reset"
        exit 1;
        shift
        ;;
        *)
        ;;
    esac
    shift
done

if [ "$load_config" == "true" ]; then
    if [ -f .env.config ]; then
        echo "File .env.config exists."

            while IFS= read -r line; do
                if [[ "$line" =~ ^(BUX_DB_DATASTORE_ENGINE=) ]]; then
                    value="${line#*=}"
                    database="${value//\"}"
                fi
                if [[ "$line" =~ ^(BUX_CACHE_ENGINE=) ]]; then
                    value="${line#*=}"
                    cache="${value//\"}"
                fi
                if [[ "$line" =~ ^(RUN_BUX_SERVER=) ]]; then
                    value="${line#*=}"
                    bux_server="${value//\"}"
                fi
                if [[ "$line" =~ ^(RUN_BUX_SERVER_BACKGROUND=) ]]; then
                    value="${line#*=}"
                    background="${value//\"}"
                fi
                if [[ "$line" =~ ^(BUX_AUTH_ADMIN_KEY=) ]]; then
                    value="${line#*=}"
                    admin_xpub="${value//\"}"
                fi
            done < ".env.config"
        else
            echo "File .env.config does not exist."
        fi
fi

if [ "$database" == "" ]; then
    database_options=("postgresql" "mongodb" "sqlite")
    ask_for_choice "\033[1mSelect your database: $reset" "${database_options[@]}"

    case $choice in
        1) database="postgresql";;
        2) database="mongodb";;
        3) database="sqlite";;
    esac
fi

if [ "$cache" == "" ]; then
    cache_options=("freecache" "redis")
    ask_for_choice "\033[1mSelect your cache storage:$reset" "${cache_options[@]}"

    case $choice in
        1) cache="freecache";;
        2) cache="redis";;
    esac
fi

if [ "$load_config" != "true" ]; then
    # Create the .env.config file
    echo -e "\033[0;32mCreating .env.config file...$reset"
    cat << EOF > .env.config
BUX_CACHE_ENGINE="$cache"
BUX_DB_DATASTORE_ENGINE="$database"
EOF

    # Add additional settings to .env.config file based on the selected database
    if [ "$database" == "postgresql" ]; then
        echo 'BUX_DB_SQL_HOST="bux-postgresql"' >> .env.config
        echo 'BUX_DB_SQL_NAME="postgres"' >> .env.config
        echo 'BUX_DB_SQL_USER="postgres"' >> .env.config
        echo 'BUX_DB_SQL_PASSWORD="postgres"' >> .env.config
    fi

    # Add additional settings to .env.config file based on the selected database
    if [ "$database" == "mongodb" ]; then
        echo 'BUX_DB_MONGODB_URI="mongodb://mongo:mongo@bux-mongodb:27017/"' >> .env.config
    fi

    # Add additional settings to .env.config file based on the selected cache storage
    if [ "$cache" == "redis" ]; then
        echo 'BUX_CACHE_REDIS_URL="redis://redis:6379"' >> .env.config
    fi
fi

echo -e "\033[0;32mStarting additional services with docker-compose...$reset"
if [ "$cache" == "redis" ]; then
    echo -e "\033[0;37mdocker compose up -d bux-redis$reset"
    docker_compose_up "-d bux-redis"
fi

if [ "$database" != "sqlite" ]; then
    echo -e "\033[0;37mdocker compose up -d bux-$database$reset"
    docker_compose_up "-d bux-$database"
fi

if [ "$bux_server" == "" ]; then
    ask_for_yes_or_no "\033[1mDo you want to run Bux-server?$reset"
    bux_server=$choice
fi

if [ "$load_config" != "true" ]; then
    echo "RUN_BUX_SERVER=\"$bux_server\"" >> .env.config
fi

if [ "$bux_server" == "true" ]; then
    if [ "$load_config" != "true" ]; then
        if [ "$admin_xpub" == "" ]; then
            # Ask for admin xPub choice
            echo -e "\033[1mDefine admin xPub $reset"
            echo -e "\033[4mLeave empty to use the one from selected environment config file $reset"
            read -p "> " admin_input

            if [[ -n "$admin_input" ]]; then
                admin_xpub=$admin_input
            fi
        fi
    fi

    if [ "$rebuild" == "" ]; then
        ask_for_yes_or_no "\033[1mDo you want to rebuild docker images before running?$reset"
        rebuild=$choice
    fi

    if [ "$rebuild" == "true" ]; then
        docker_command_attributes="--build"
    fi

    if [ "$background" == "" ]; then
        ask_for_yes_or_no "\033[1mDo you want to run Bux-server in background?$reset"
        background=$choice
    fi

    if [ "$load_config" != "true" ]; then
        echo "RUN_BUX_SERVER_BACKGROUND=\"$background\"" >> .env.config

        if [ "$admin_xpub" != "" ]; then
            echo "BUX_AUTH_ADMIN_KEY=\"$admin_xpub\"" >> .env.config
        fi
    fi

    echo -e "\033[0;32mRunning Bux-server...$reset"
    if [ "$background" == "true" ]; then
        echo -e "\033[0;37mdocker compose up $docker_command_attributes -d bux-server$reset"
        docker_compose_up "$docker_command_attributes -d bux-server"
    else
        echo -e "\033[0;37mdocker compose up $docker_command_attributes bux-server$reset"
        docker_compose_up "$docker_command_attributes bux-server"

        function cleanup {
            echo -e "\033[0;31mStopping additional services...$reset"
            if [ compose_plugin == true ]; then
                docker compose stop
            else
                docker-compose stop
            fi
            echo -e "\033[0;31mExiting program...$reset"
        }

        trap cleanup EXIT
    fi
else
    echo -e "\033[0;33m\033[1mThanks for using Bux configurator!"
    echo -e "Additional services are working, remember to start Bux-server manually!$reset"
    exit 1
fi
