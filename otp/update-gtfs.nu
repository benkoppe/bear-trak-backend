#!/usr/bin/env nu

const SCRIPT_PATH = (path self)

def git_repo_root [] {
  let res = (^git rev-parse --show-toplevel | complete)
  if $res.exit_code != 0 {
    ""
  } else {
    $res.stdout | str trim
  }
}

def prefetch_nix_hash [url: string] {
  let res = (^nix store prefetch-file --hash-type sha256 --json $url | complete)
  if $res.exit_code != 0 {
    return {
      ok: false
      error: ($res.stderr | str trim)
    }
  }

  let parsed = (try { $res.stdout | from json } catch {
    return {
      ok: false
      error: "prefetch output was not valid JSON"
    }
  })

  {
    ok: true
    hash: ($parsed | get hash)
  }
}

def nix_escape [s: string] {
  $s
    | str replace -a '\\' '\\\\'
    | str replace -a '"' '\\"'
}

def render_sources_nix [sources: record] {
  let nl = (char nl)
  mut out = "{" + $nl

  for school in (($sources | columns) | sort) {
    let school_rec = ($sources | get $school)
    let gtfs = ($school_rec | get gtfs)

    $out = $out + $"  ($school) = {" + $nl
    $out = $out + $"    gtfs = {" + $nl

    for feed_name in (($gtfs | columns) | sort) {
      let feed = ($gtfs | get $feed_name)
      let url = (nix_escape $feed.url)
      let sha256 = (nix_escape $feed.sha256)
      $out = $out + $"      \"($feed_name)\" = {" + $nl
      $out = $out + $"        url = \"($url)\";" + $nl
      $out = $out + $"        sha256 = \"($sha256)\";" + $nl
      $out = $out + $"      };" + $nl
    }

    $out = $out + $"    };" + $nl
    $out = $out + $"  };" + $nl + $nl
  }

  $out + "}" + $nl
}

def default_sources_path [] {
  let script_path = $SCRIPT_PATH

  if (($script_path | str trim) != "") {
    let script_dir = ($script_path | path dirname)

    if ($script_dir | str starts-with "/nix/store") {
      let repo_root = (git_repo_root)
      if $repo_root != "" {
        let otp_sources = ($repo_root | path join "otp" "gtfs-sources.nix")
        if ($otp_sources | path exists) {
          return $otp_sources
        }

        let repo_sources = ($repo_root | path join "gtfs-sources.nix")
        if ($repo_sources | path exists) {
          return $repo_sources
        }
      }
    }

    return ($script_dir | path join "gtfs-sources.nix")
  }

  "gtfs-sources.nix"
}

def load_sources [sources_path: string] {
  let res = (^nix eval --json --file ($sources_path | path expand) | complete)
  if $res.exit_code != 0 {
    error make {
      msg: $"failed to read ($sources_path) via nix eval: ($res.stderr | str trim)"
    }
  }

  try {
    $res.stdout | from json
  } catch {
    error make {
      msg: $"failed to parse nix eval output from ($sources_path)"
    }
  }
}

def refresh_hashes [sources: record] {
  mut updated = $sources
  mut failures = []

  for school in (($sources | columns) | sort) {
    let school_rec = ($sources | get $school)
    let gtfs = ($school_rec | get gtfs)
    mut school_out = $gtfs

    for feed_name in (($gtfs | columns) | sort) {
      let feed = ($gtfs | get $feed_name)
      let url = ($feed | get url)
      let previous_hash = ($feed | get sha256)

      # Always re-prefetch: GTFS contents change even when URL doesn't.
      let prefetch = (prefetch_nix_hash $url)
      if ($prefetch | get ok) == false {
        let msg = ($prefetch | get error)
        print $"warning ($school)/($feed_name): prefetch failed: ($msg)"
        $failures = ($failures | append $"($school)/($feed_name)")
        continue
      }

      let new_hash = ($prefetch | get hash)
      if $new_hash == $previous_hash {
        print $"unchanged ($school)/($feed_name): ($new_hash)"
      } else {
        print $"updated ($school)/($feed_name): ($new_hash)"
      }

      $school_out = ($school_out | upsert $feed_name ($feed | upsert sha256 $new_hash))
    }

    $updated = ($updated | upsert $school ($school_rec | upsert gtfs $school_out))
  }

  {
    sources: $updated
    failures: $failures
  }
}

def main [
  --strict
  sources_path?: string
] {
  let sources_path = ($sources_path | default (default_sources_path))

  if not ($sources_path | path exists) {
    error make { msg: $"sources file does not exist: ($sources_path)" }
  }

  let sources = (load_sources $sources_path)
  let refresh = (refresh_hashes $sources)

  let rendered = (render_sources_nix ($refresh | get sources))
  let prev_raw = (open $sources_path --raw)

  if (($prev_raw | str trim --right) == ($rendered | str trim --right)) {
    print $"unchanged sources file: ($sources_path)"
  } else {
    $rendered | save --force $sources_path
    print $"updated sources file: ($sources_path)"
  }

  let failures = ($refresh | get failures)
  if (($failures | length) > 0) {
    print $"warning: failed to refresh hashes for ($failures | length) feeds: (($failures | str join ', '))"
    if $strict {
      exit 1
    }
  }

  let git_res = (^git ls-files --error-unmatch $sources_path | complete)
  if ($git_res.exit_code != 0) {
    print $"warning: ($sources_path) is not tracked by git; flakes will fail until you run: git add ($sources_path)"
  }
}
