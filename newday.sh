#!/bin/bash

set -euo pipefail
shopt -s extglob

THIS_DIR=$(dirname "${BASH_SOURCE[0]}")
[[ -f "${THIS_DIR}/.env" ]] && source "${THIS_DIR}/.env"

last_day=$(basename "$(printf "%s\n" "${THIS_DIR}"/[0-9]?([0-9]) | tail -1)")
next_day=$((last_day + 1))

mkdir "${next_day}"
cp "${THIS_DIR}/template.go" "${THIS_DIR}/${next_day}/solution.go"
curl -f --cookie "session=${SESSION_COOKIE}" \
    "https://adventofcode.com/2024/day/${next_day}/input" \
    >"${THIS_DIR}/${next_day}/input.txt"

git -C "${THIS_DIR}" add "${next_day}"
git -C "${THIS_DIR}" commit -m "Initialize day ${next_day}"
