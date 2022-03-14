#!/bin/bash 

function a711() {
    if ! which jq; then
        ( >&2 cat <<EOF
jq isn't installed. You can install it using:
    https://stedolan.github.io/jq/download/
EOF
)
        return;
    fi

    if ! which op; then
        ( >&2 cat <<EOF
op CLI tool isn't installed. You can install it using the following guide:
    https://support.1password.com/command-line-getting-started/
EOF
        )
        return;
    fi

    ## initialze op session if it's not ready
    (op list items >/dev/null 2>&1) || (op signin)

    command=$1
    shift

    case "$command" in
        load)
            a711_load_script "$1"
            ;;
        save)
            true
            ;;
        *)
            ( >&2 echo "unrecognized command $command" )
            return
            ;;
    esac

}

function a711_load_script() {
    echo "loading credential $1"
}