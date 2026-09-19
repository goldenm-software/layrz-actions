#!/usr/bin/env bash
#
# dry_run_publish.sh
#
# What it does:
#   Runs "flutter pub publish --dry-run" for the package in the current
#   working directory and decides pass/fail by inspecting the actual output,
#   instead of trusting the dry-run exit code alone. This is the script
#   bundled with the goldenm-software/layrz-actions "dry-run-publish"
#   composite action; the action's Setup Flutter and Install dependencies
#   steps run before this script, and this script runs the dry-run itself
#   (which re-resolves dependencies on its own, same as pub get did). Two
#   independent env-var gates control what gets ignored: IGNORE_ADVISORIES
#   for known-benign advisory warnings, and IGNORE_MISSING_METADATA for
#   missing README.md/CHANGELOG.md files on internal monorepo packages.
#
# Why it exists:
#   Some Layrz Flutter packages are git-only (publish_to: none in
#   pubspec.yaml) and depend on pre-release versions of other Layrz packages
#   on purpose. "flutter pub publish --dry-run" flags that as a warning
#   (a package depending on a pre-release should itself be a pre-release)
#   and returns a non-zero exit code even though nothing is actually wrong.
#   That warning is advisory for those packages and must NOT fail CI for
#   them.
#
#   Separately, the dry-run's own embedded analyzer sometimes reports a
#   spurious issue shaped like a real `dart analyze` failure, ending with
#   "Please report this at dartbug.com." Running the package's own
#   `flutter analyze` (or `dart analyze`) directly at the same time reports
#   "No issues found", so this is a quirk of the dry-run's internal analyzer
#   invocation, not a real lint problem. It must NOT fail CI either.
#
#   Separately, some Layrz packages deliberately depend on other Layrz
#   packages (e.g. layrz_utils) via a "git:" source, since those packages are
#   git-only and never published to pub.dev. "flutter pub publish --dry-run"
#   warns "Don't depend on <pkg> from the git source. Use the hosted source
#   instead." for any such dependency, because it assumes the target is
#   pub.dev, where a git dependency would break consumers. That assumption
#   does not hold for a git-only package, so this warning is also advisory
#   and must NOT fail CI for it.
#
#   Separately, in a MONOREPO some packages are internal and are never
#   published: they legitimately have no README.md and no CHANGELOG.md, so
#   "flutter pub publish --dry-run" reports "Please add a README.md file
#   that describes your package." and "Please add a `CHANGELOG.md` to your
#   package." as blocking issues even though nothing is actually wrong for
#   those packages. That is handled by a separate gate from the advisories
#   above (see IGNORE_MISSING_METADATA below), since it is a distinct
#   decision from whether advisories are ignorable.
#
#   Any other warning or error (a real packaging problem, a missing file, an
#   invalid pubspec, a genuine analyzer finding, etc) must still fail CI.
#
# The IGNORE_ADVISORIES environment variable contract:
#   IGNORE_ADVISORIES=true   Ignore the three known advisory classes above
#                            (pre-release dependency, dartbug.com analyzer
#                            quirk, git-source). Every other issue still
#                            fails closed. Use this for git-only packages
#                            that intentionally depend on pre-release and/or
#                            git-source dependencies.
#   IGNORE_ADVISORIES=false  Strict mode (the default when unset). ANY
#                            reported issue block fails the run, including
#                            the three classes above. Use this for hosted
#                            packages, where none of those advisories should
#                            ever legitimately appear.
#   Any other value is treated the same as "false" (strict).
#
# The IGNORE_MISSING_METADATA environment variable contract:
#   IGNORE_MISSING_METADATA=true   Ignore ONLY the two missing-file issues
#                                  (missing README.md, missing CHANGELOG.md).
#                                  Every other issue, including all advisory
#                                  classes, is still handled per
#                                  IGNORE_ADVISORIES above, independently of
#                                  this gate. Use this for internal monorepo
#                                  packages that are never published and
#                                  intentionally omit a README/CHANGELOG.
#   IGNORE_MISSING_METADATA=false  Strict mode (the default when unset).
#                                  Missing README.md/CHANGELOG.md fail the
#                                  run.
#   Any other value is treated the same as "false" (strict).
#
#   IGNORE_ADVISORIES and IGNORE_MISSING_METADATA are INDEPENDENT gates: a
#   block is ignored if it matches an advisory class and IGNORE_ADVISORIES
#   is enabled, OR if it matches a missing-metadata class and
#   IGNORE_MISSING_METADATA is enabled. Either gate can be on or off without
#   affecting the other.
#
# How to invoke:
#   IGNORE_ADVISORIES=true IGNORE_MISSING_METADATA=true bash dry_run_publish.sh
#   (run from the package root, so pub can find pubspec.yaml; the action
#   sets working-directory for this automatically)
#
# Exit codes:
#   0  the dry-run succeeded, or every reported issue was ignorable under
#      IGNORE_ADVISORIES and/or IGNORE_MISSING_METADATA (whichever gates
#      are enabled).
#   1  the dry-run reported at least one issue that is blocking under the
#      current mode(s), or failed for some other reason.
#
# This script is NOT read-only: it does not mutate the repository itself,
# but "flutter pub publish --dry-run" does touch pub's local cache/state as
# a side effect of resolving dependencies. It never publishes anything.

set -uo pipefail
# Note: intentionally no "set -e" here. We need to capture the dry-run's own
# exit code ourselves (it is expected to be non-zero on the ignorable
# advisories, when ignoring is enabled) and react to it explicitly, rather
# than have the script abort before we get a chance to inspect the output.

# The exact, stable substring that identifies the ignorable pre-release
# advisory class. Do NOT add specific package names or version numbers here,
# they change every time a dependency bumps its pre-release version.
readonly IGNORABLE_PATTERN='Packages dependent on a pre-release of another package should themselves be published as a pre-release version'

# The exact, stable substring that identifies the ignorable dry-run analyzer
# quirk. "Please report this at dartbug.com." only ever appears in this one
# spurious issue, so it is specific enough to not accidentally swallow a real
# analyzer error, which would instead list actual file:line diagnostics and
# never mention dartbug.com. Deliberately NOT matching on "dart analyze"
# alone, since that phrase alone would also match a real analyzer failure.
readonly IGNORABLE_DARTBUG_PATTERN='Please report this at dartbug.com'

# The exact, stable substring that identifies the ignorable "git source"
# advisory. A git-only package (publish_to: none) may depend on other
# packages from git on purpose, so pub's advice to use the hosted source
# instead does not apply. Match on the stable phrase common to this warning
# for ANY git dependency, deliberately NOT hardcoding a package name or
# version, so it keeps working if the git dependency set changes. This
# phrase only appears in this specific advisory, never in a real packaging
# error, so it is specific enough not to swallow a genuine problem.
readonly IGNORABLE_GIT_SOURCE_PATTERN='from the git source. Use the hosted source instead'

# The exact, stable substring that identifies the ignorable "missing
# README.md" issue. Some internal monorepo packages are never published and
# intentionally have no README.md, so this issue is advisory for them under
# IGNORE_MISSING_METADATA. This phrase only appears in this specific issue,
# never in a real packaging error, so it is specific enough not to swallow a
# genuine problem.
readonly IGNORABLE_MISSING_README_PATTERN='Please add a README.md file that describes your package'

# The exact, stable substring that identifies the ignorable "missing
# CHANGELOG.md" issue. Note the backticks around CHANGELOG.md are literal
# characters in pub's own output, matched here exactly as printed. This
# string is single-quoted so the backticks are NOT command-substituted by
# the shell. Some internal monorepo packages are never published and
# intentionally have no CHANGELOG.md, so this issue is advisory for them
# under IGNORE_MISSING_METADATA.
# The single quotes are deliberate (see above), so the backticks are literal
# text, not command substitution; silence the resulting shellcheck info.
# shellcheck disable=SC2016
readonly IGNORABLE_MISSING_CHANGELOG_PATTERN='Please add a `CHANGELOG.md` to your package'

# is_ignore_mode: reads the IGNORE_ADVISORIES env var and reports whether
# advisory ignoring is enabled.
# Takes: nothing (reads the IGNORE_ADVISORIES environment variable, which
#   defaults to "false" when unset).
# Prints: nothing.
# Returns: 0 (true) when IGNORE_ADVISORIES is exactly "true", 1 (false)
#   otherwise, including when it is unset or any other value.
is_ignore_mode() {
  local value="${IGNORE_ADVISORIES:-false}"
  [[ "${value}" == "true" ]]
}

# is_ignore_missing_metadata: reads the IGNORE_MISSING_METADATA env var and
# reports whether ignoring missing-metadata-file issues (README.md,
# CHANGELOG.md) is enabled.
# Takes: nothing (reads the IGNORE_MISSING_METADATA environment variable,
#   which defaults to "false" when unset).
# Prints: nothing.
# Returns: 0 (true) when IGNORE_MISSING_METADATA is exactly "true", 1
#   (false) otherwise, including when it is unset or any other value.
is_ignore_missing_metadata() {
  local value="${IGNORE_MISSING_METADATA:-false}"
  [[ "${value}" == "true" ]]
}

# main: runs the dry-run, parses its output into issue blocks, classifies
# each block as ignorable or blocking depending on IGNORE_ADVISORIES, and
# decides the script's exit status.
# Takes: no arguments.
# Prints: the full dry-run output, the exit status it reported, and a
#   summary of which issue blocks were ignored or found blocking.
# Returns: 0 when nothing blocking was found, 1 otherwise.
main() {
  local output
  local status
  local ignore_mode
  local metadata_ignore_mode

  if is_ignore_mode; then
    ignore_mode="true"
  else
    ignore_mode="false"
  fi
  echo "dry_run_publish: IGNORE_ADVISORIES=${ignore_mode}"

  if is_ignore_missing_metadata; then
    metadata_ignore_mode="true"
  else
    metadata_ignore_mode="false"
  fi
  echo "dry_run_publish: IGNORE_MISSING_METADATA=${metadata_ignore_mode}"

  output="$(flutter pub publish --dry-run 2>&1)"
  status=$?

  # Always print the full output so CI logs show exactly what happened.
  printf '%s\n' "${output}"
  printf '\n--- dry_run_publish: dry-run exited with status %d ---\n' "${status}"

  # Pull out every reported issue as a block: the "* " issue line itself plus
  # every following line up to (not including) the next "* " issue line or a
  # blank separator line. This matters because some issues, notably the
  # dartbug.com analyzer quirk, put their distinctive text on an indented
  # body line rather than on the "* " line itself, so classifying by the "* "
  # line alone would miss it. Each block is flattened to one newline-joined
  # string so it can be matched and stored as a single array element.
  #
  # Grep with a fallback so "no matches" (exit 1 from grep) does not trip
  # set -o pipefail below and does not abort the script.
  local has_issue_lines
  has_issue_lines="$(printf '%s\n' "${output}" | grep -cE '^\* ' || true)"

  # No issue lines at all: nothing to classify, nothing blocking. Decide
  # purely on the dry-run's own exit code.
  if [[ "${has_issue_lines}" -eq 0 ]]; then
    if [[ "${status}" -eq 0 ]]; then
      echo "dry_run_publish: no issues reported, passing"
      return 0
    fi
    echo "dry_run_publish: found blocking issues, failing"
    echo "dry_run_publish: dry-run failed with no issue lines to explain it (status ${status})"
    return 1
  fi

  # Build the list of issue blocks. Use process substitution so the while
  # loop runs in the current shell and can update the array below.
  local -a issue_blocks=()
  local current_block=""
  local line
  while IFS= read -r line; do
    if [[ "${line}" == "* "* ]]; then
      # Start of a new issue: flush the previous block, if any.
      if [[ -n "${current_block}" ]]; then
        issue_blocks+=("${current_block}")
      fi
      current_block="${line}"
    elif [[ -z "${line}" ]]; then
      # Blank line: ends the current block, if one is open.
      if [[ -n "${current_block}" ]]; then
        issue_blocks+=("${current_block}")
        current_block=""
      fi
    elif [[ -n "${current_block}" ]]; then
      # Body/continuation line belonging to the open block.
      current_block="${current_block}"$'\n'"${line}"
    fi
  done <<<"${output}"
  # Flush a trailing block that was not closed by a blank line (the output
  # ended right after its last body line).
  if [[ -n "${current_block}" ]]; then
    issue_blocks+=("${current_block}")
  fi

  # Classify each issue block as ignorable or blocking. A block is ignorable
  # under either of two INDEPENDENT gates:
  #   - the advisory gate: is_ignore_mode is enabled AND the block matches
  #     one of the three known advisory patterns (pre-release dependency,
  #     dartbug.com analyzer quirk, git-source), or
  #   - the metadata gate: is_ignore_missing_metadata is enabled AND the
  #     block matches the missing-README or missing-CHANGELOG pattern.
  # Everything else is blocking. Matching against the whole block, not just
  # its "* " line, is what lets the dartbug.com substring on a body line
  # still mark that issue ignorable. When a gate is disabled, its patterns
  # never make a block ignorable: strict mode fails on any reported issue
  # not covered by the other (enabled) gate.
  local -a blocking_lines=()
  local -a ignored_lines=()
  local block
  local advisory_ignorable
  local metadata_ignorable
  for block in "${issue_blocks[@]}"; do
    advisory_ignorable="false"
    if is_ignore_mode && { [[ "${block}" == *"${IGNORABLE_PATTERN}"* ]] || [[ "${block}" == *"${IGNORABLE_DARTBUG_PATTERN}"* ]] || [[ "${block}" == *"${IGNORABLE_GIT_SOURCE_PATTERN}"* ]]; }; then
      advisory_ignorable="true"
    fi

    metadata_ignorable="false"
    if is_ignore_missing_metadata && { [[ "${block}" == *"${IGNORABLE_MISSING_README_PATTERN}"* ]] || [[ "${block}" == *"${IGNORABLE_MISSING_CHANGELOG_PATTERN}"* ]]; }; then
      metadata_ignorable="true"
    fi

    if [[ "${advisory_ignorable}" == "true" ]] || [[ "${metadata_ignorable}" == "true" ]]; then
      ignored_lines+=("${block}")
    else
      blocking_lines+=("${block}")
    fi
  done

  if [[ "${#blocking_lines[@]}" -eq 0 ]]; then
    # Every issue block is ignorable under one of the two gates: the known
    # advisory classes (pre-release, dartbug.com analyzer quirk,
    # git-source) and/or the missing-metadata-file classes (missing
    # README.md, missing CHANGELOG.md). Expected for a git-only package
    # intentionally pinned to pre-release and/or git dependencies, and/or
    # an internal monorepo package that intentionally omits README/CHANGELOG.
    echo "dry_run_publish: only ignorable issues found (pre-release, dartbug.com analyzer quirk, git-source advisories, and/or missing README.md/CHANGELOG.md), ignoring"
    echo "dry_run_publish: ignored ${#ignored_lines[@]} issue block(s):"
    printf -- '--- ignored issue block ---\n%s\n' "${ignored_lines[@]}"
    return 0
  fi

  # At least one issue block is blocking: either it is not one of the known
  # ignorable classes, or advisory ignoring is disabled (strict mode), in
  # which case every reported issue is blocking. Fail closed regardless of
  # what the dry-run's own exit code was.
  echo "dry_run_publish: found blocking issues, failing"
  echo "dry_run_publish: ${#blocking_lines[@]} blocking issue block(s):"
  printf -- '--- blocking issue block ---\n%s\n' "${blocking_lines[@]}"
  return 1
}

main "$@"
