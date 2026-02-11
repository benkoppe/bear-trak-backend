#!/usr/bin/env nu

def get_nix_hash [url: string] {
  nix store prefetch-file --hash-type sha256 --json $url
    | from json
    | get hash
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

let feeds = {
  cornell: {
    "gtfs.zip": {
      url: "https://realtimetcatbus.availtec.com/InfoPoint/GTFS-zip.ashx"
    }
  }

  harvard: {
    "gtfs.zip": {
      url: "https://passio3.com/harvard/passioTransit/gtfs/google_transit.zip"
    }
    "gtfs-mbta.zip": {
      url: "https://cdn.mbta.com/MBTA_GTFS.zip"
    }
  }

  umich: {
    "gtfs.zip": {
      url: "https://webapps.fo.umich.edu/transit_uploads/google_transit.zip"
    }
  }
}

let script_dir = ($env.FILE_PWD? | default "")
let sources_path = if $script_dir != "" {
    $script_dir | path join "gtfs-sources.nix"
  } else {
    "gtfs-sources.nix"
  }

mut sources = {}

for row in ($feeds | transpose school schoolFeeds) {
  let school = $row.school
  let schoolFeeds = $row.schoolFeeds
  mut schoolOut = {}

  for feed_row in ($schoolFeeds | transpose name f) {
    let name = $feed_row.name
    let f = $feed_row.f
    let url = $f.url

    # Always re-prefetch: GTFS contents change even when URL doesn't.
    let hash = (try { get_nix_hash $url } catch {|err|
      print $"warning ($school)/($name): prefetch failed: ($err.msg | default $err)"
      null
    })
    if $hash == null {
      continue
    }

    print $"prefetched ($school)/($name): ($hash)"
    $schoolOut = ($schoolOut | upsert $name { url: $url sha256: $hash })
  }

  $sources = ($sources | upsert $school { gtfs: $schoolOut })
}

let rendered = (render_sources_nix $sources)
let prev_raw = if ($sources_path | path exists) { open $sources_path --raw } else { "" }
if (($prev_raw | str trim --right) == ($rendered | str trim --right)) {
  print $"unchanged sources file: ($sources_path)"
} else {
  $rendered | save --force $sources_path
  print $"updated sources file: ($sources_path)"
}

let git_res = (do { ^git ls-files --error-unmatch $sources_path } | complete)
if ($git_res.exit_code != 0) {
  print $"warning: ($sources_path) is not tracked by git; flakes will fail until you run: git add ($sources_path)"
}
