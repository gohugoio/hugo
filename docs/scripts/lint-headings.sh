#!/usr/bin/env bash

#------------------------------------------------------------------------------
# @file
# Verifies heading structure across the site.
#
# Every Markdown file under content/en must not contain a level 1, level 5,
# or level 6 heading. In addition, every Markdown file within the
# directories listed in the sections array in the main function below,
# excluding _index.md files and the paths listed in the exclusions array,
# must contain exactly two level 2 headings, in this order: "Usage" and
# "Examples". Reports one error message per file per check that does not
# comply.
#------------------------------------------------------------------------------

#------------------------------------------------------------------------------
# @function
# Displays usage message.
#
# @return string $msg
#   The usage message.
#------------------------------------------------------------------------------
usage() {
  declare msg
  msg=$(cat <<EOT
Verify heading structure across the site: no level 1, level 5, or level 6
headings anywhere, and standard level 2 headings, Usage followed by
Examples, within the configured sections.

Run with no arguments to discover and check every Markdown file under
content/en. Alternatively, pass one or more file paths to check specific
files.

Usage:    $(basename "$0") [file ...]

Example:  $(basename "$0") content/en/functions/compare/Eq.md
EOT
  )
  printf "%s\\n" "${msg}"
  exit 1
}

#------------------------------------------------------------------------------
# @function
# Reports whether the given path matches any of the given excluded paths.
#
# @param string $1
#   The path to check.
# @param string $2...
#   Zero or more excluded paths.
#
# @return bool
#   True if excluded.
#------------------------------------------------------------------------------
is_excluded() {
  declare -r path="$1"
  shift
  declare excluded
  for excluded in "$@"; do
    [[ "${path}" == "${excluded}" ]] && return 0
  done
  return 1
}

#------------------------------------------------------------------------------
# @function
# Reports whether the given path is within any of the given section
# directories.
#
# @param string $1
#   The path to check.
# @param string $2...
#   Zero or more section directories.
#
# @return bool
#   True if within a section.
#------------------------------------------------------------------------------
in_section() {
  declare -r path="$1"
  shift
  declare section
  for section in "$@"; do
    [[ "${path}" == "${section}"/* ]] && return 0
  done
  return 1
}

#------------------------------------------------------------------------------
# @function
# Main function.
#
# With no arguments, discovers every Markdown file under content/en, then
# re-invokes this script with that file list. With one or more arguments,
# checks each given file for a disallowed heading level, and, for files
# within the sections array that are not _index.md files or listed in the
# exclusions array, for the expected level 2 headings. Reports a separate
# error message for each check that fails.
#
# @param string $@
#   Zero or more file paths to check.
#------------------------------------------------------------------------------
main() {
  # If run with no arguments, discover files and re-run this script with the file list.
  if (( $# == 0 )); then
    find content/en -type f -name "*.md" -print0 | xargs -0 "$0"
    exit $?
  fi

  # Directories, relative to the project root, whose Markdown files must
  # contain exactly two level 2 headings: Usage and Examples. Add to this
  # list as needed.
  declare -ra sections=(
    "content/en/functions"
    "content/en/methods"
  )

  # Paths, relative to the project root, to exclude from the Usage and
  # Examples heading check. Add to this list as needed.
  declare -ra exclusions=(
    "content/en/functions/js/Batch.md"
  )

  declare -i errors=0
  declare -r expected=$'Usage\nExamples'
  declare actual file

  for file in "$@"; do
    if grep -qE '^# |^##### |^###### ' "${file}"; then
      echo "Invalid heading structure in ${file}. Found a level 1, level 5, or level 6 heading."
      errors+=1
    fi

    in_section "${file}" "${sections[@]}" || continue
    [[ "$(basename "${file}")" == "_index.md" ]] && continue
    is_excluded "${file}" "${exclusions[@]}" && continue

    actual=$(grep -E '^## ' "${file}" | sed 's/^## *//; s/ *$//' || true)
    if [[ "${actual}" != "${expected}" ]]; then
      echo "Invalid heading structure in ${file}. Expected exactly two level 2 headings, Usage and Examples."
      errors+=1
    fi
  done

  (( errors == 0 ))
}

set -euo pipefail
main "$@"
