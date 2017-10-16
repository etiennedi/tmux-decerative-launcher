#!/bin/bash

CONFIG_FILE=config.yml

function main() {
  check_dependencies
  check_config
  check_window_is_configured "$1"
}

function check_dependencies() {
  type yq >/dev/null 2>&1 || { echo >&2 "yq not found on path. Please instal first"; exit 1; }
}

function check_config() {
  check_config_exists
  check_config_valid
}

function check_config_exists {
    if [ ! -f "$CONFIG_FILE" ]; then
      echo >&2 "no $CONFIG_FILE exists"; exit 1
    fi
}

function check_config_valid {
  if ! yq -e '.windows' < "$CONFIG_FILE" 2>&1 >/dev/null; then
    echo >&2 "invalid config, needs to have 'windows' on top level"; exit 1
  fi
}

function check_window_is_configured {
  if ! window=$(yq -e ".windows[\"${1}\"]" < "$CONFIG_FILE") ; then
    echo >&2 "window ${1} is not configured"; exit 1
  fi

  echo "$window"
}

main "$@"
