#!/bin/bash

set -euo pipefail
shopt -s extglob

spin_chrs=(
    '\u28b9'
    '\u28f8'
    '\u28f4'
    '\u28e6'
    '\u28c7'
    '\u284f'
    '\u281f'
    '\u283b'
)

THIS_DIR=$(dirname "${BASH_SOURCE[0]}")
[[ -f "${THIS_DIR}/.env" ]] && source "${THIS_DIR}/.env"

next_day="$1"

starts_at=$(date +%s -d "2024-12-${next_day}T00:00:00EST")

while true; do
    seconds_left=$((starts_at - $(date +%s)))
    [[ $((starts_at - $(date +%s))) -lt 0 ]] && break
    for i in {0..7}; do
        echo -ne "\r${spin_chrs[$i]} ${seconds_left}\t\t"
        sleep 0.125
    done
done
echo -e "\rBegin day ${next_day}!"

mkdir "${THIS_DIR}/${next_day}"
cp "${THIS_DIR}/template.go" "${THIS_DIR}/${next_day}/solution.go"
curl -f --cookie "session=${SESSION_COOKIE}" \
    "https://adventofcode.com/2024/day/${next_day}/input" \
    >"${THIS_DIR}/${next_day}/input.txt"

git -C "${THIS_DIR}" add "${next_day}"
git -C "${THIS_DIR}" commit -m "Initialize day ${next_day}"
code "${THIS_DIR}/${next_day}/solution.go"
