#!/usr/bin/env bats

load_lib() {
    # Mock gh command
    function gh() {
        if [[ "$1" == "auth" && "$2" == "status" ]]; then
            return 0
        fi
        if [[ "$1" == "extension" && "$2" == "install" ]]; then
            echo "Installed extension"
            return 0
        fi
        if [[ "$1" == "copilot" ]]; then
            echo "Copilot output"
            return 0
        fi
    }
    export -f gh
}

setup() {
    export SRC_DIR="$(dirname "$BATS_TEST_DIRNAME")"
    export PATH="$SRC_DIR:$PATH"
}

@test "install_gh.sh: skips if gh is installed" {
    # Mock command -v gh to return true
    function command() {
        if [[ "$1" == "-v" && "$2" == "gh" ]]; then
            return 0
        fi
        builtin command "$@"
    }
    export -f command

    run bash "$SRC_DIR/install_gh.sh"
    [ "$status" -eq 0 ]
    [[ "$output" == *"gh version"* ]]
}

@test "configure_gh_auth.sh: fails without GH_TOKEN" {
    run bash "$SRC_DIR/configure_gh_auth.sh"
    [ "$status" -eq 1 ]
    [[ "$output" == *"::error::GH_TOKEN is not set"* ]]
}

@test "resolve_mission.sh: fails if no input" {
    run bash "$SRC_DIR/resolve_mission.sh"
    [ "$status" -eq 1 ]
    [[ "$output" == *"::error::Either 'mission' or 'template' input must be provided."* ]]
}

@test "resolve_mission.sh: resolves mission input" {
    export INPUT_MISSION="Analyze code"
    # Create a temp file for GITHUB_ENV to ensure it works in any environment
    export GITHUB_ENV_FILE="$(mktemp)"
    # Temporarily set GITHUB_ENV to our temp file
    OLD_GITHUB_ENV="${GITHUB_ENV:-}"
    export GITHUB_ENV="$GITHUB_ENV_FILE"

    run bash "$SRC_DIR/resolve_mission.sh"
    [ "$status" -eq 0 ]
    
    # Check if the variable was written to GITHUB_ENV
    grep -q "RESOLVED_MISSION<<EOF" "$GITHUB_ENV_FILE"
    grep -q "Analyze code" "$GITHUB_ENV_FILE"

    # Cleanup
    rm "$GITHUB_ENV_FILE"
    if [ -n "$OLD_GITHUB_ENV" ]; then
        export GITHUB_ENV="$OLD_GITHUB_ENV"
    else
        unset GITHUB_ENV
    fi
    unset INPUT_MISSION
}

@test "resolve_mission.sh: fails if template not found" {
    export INPUT_TEMPLATE="non-existent"
    run bash "$SRC_DIR/resolve_mission.sh"
    [ "$status" -eq 1 ]
    [[ "$output" == *"::error::Template file not found"* ]]
    unset INPUT_TEMPLATE
}
