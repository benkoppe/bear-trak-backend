{
  description = "OpenTripPlanner";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    devenv.url = "github:cachix/devenv";
    nix2container = {
      url = "github:nlewo/nix2container";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.devenv.flakeModule
      ];
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "aarch64-darwin"
        "x86_64-darwin"
      ];
      perSystem =
        {
          pkgs,
          lib,
          ...
        }:
        let
          otpVersion = "2.8.1";
          otpShaded = pkgs.fetchurl {
            url = "https://github.com/opentripplanner/OpenTripPlanner/releases/download/v${otpVersion}/otp-shaded-${otpVersion}.jar";
            sha256 = "sha256-60Ikyw9Z7ZE+4kDReDgJGg5totgdghbDHAfFFQKDFpc=";
          };
          ramFlag = "-XX:MaxRAMPercentage=75";
          schools =
            lib.mapAttrs
              (
                name: attrs:
                attrs
                // rec {
                  gtfsPackages = lib.mapAttrs (gtfsName: g: pkgs.fetchurl { inherit (g) url sha256; }) attrs.gtfs;
                  dataDir = ./data/${name};
                  otpRoot = pkgs.runCommand "otp-root-${name}" { } ''
                    mkdir -p $out
                    cp ${dataDir}/* $out

                    ${lib.concatStringsSep "\n" (
                      lib.mapAttrsToList (gtfsName: pkg: "cp ${pkg} $out/${gtfsName}") gtfsPackages
                    )}
                  '';
                }
              )
              {
                cornell.gtfs."gtfs.zip" = {
                  url = "https://realtimetcatbus.availtec.com/InfoPoint/GTFS-zip.ashx";
                  sha256 = "sha256-oAvmvfqJxd3o2bcTFuJp5M+u9unIjkxKypPCJ7KahIA=";
                };
                harvard.gtfs = {
                  "gtfs.zip".url = "https://passio3.com/harvard/passioTransit/gtfs/google_transit.zip";
                  "gtfs.zip".sha256 = "sha256-EpMUc81e4ESTWzZeTNwjBzuOzGJlXILKq0pbPorBe6Q=";
                  "gtfs-mbta.zip".url = "https://cdn.mbta.com/MBTA_GTFS.zip";
                  "gtfs-mbta.zip".sha256 = "sha256-UYTv2cHiB4GbX+0vJCFOJb5oATwrnkgxoT8knI14o7Y=";
                };
                umich.gtfs."gtfs.zip" = {
                  url = "https://webapps.fo.umich.edu/transit_uploads/google_transit.zip";
                  sha256 = "sha256-9oT89xjuneNMXZbwA08RI/d7e7BsfI7AS+wuk0S/DKU=";
                };
              };
        in
        {
          devenv.shells.default =
            { config, ... }:
            {
              packages = with pkgs; [
                osmium-tool
              ];
              languages.java = {
                enable = true;
                jdk.package = pkgs.jdk21;
                maven.enable = true;
              };
              scripts = lib.mkMerge [
                {
                  otp.exec = "${config.languages.java.jdk.package}/bin/java ${ramFlag} -jar ${otpShaded} $@";
                  otp-build.exec = "otp --build $@";
                }
                (lib.mapAttrs' (
                  school: schoolAttrs:
                  lib.nameValuePair "otp-${school}" {
                    exec = "otp-build ${schoolAttrs.otpRoot} $@";
                  }
                ) schools)
              ];

            };
        }
        // (
          let
            mkImage =
              school: schoolAttrs:
              let
                otpGraph = pkgs.runCommand "otp-graph-${school}" { } ''
                  mkdir work
                  cp ${schoolAttrs.otpRoot}/* work/
                  cd work
                  ${pkgs.jdk21_headless}/bin/java ${ramFlag} -jar ${otpShaded} --build --save .
                  mkdir -p $out
                  cp ./*-config.json ./graph.obj $out
                '';
              in
              pkgs.dockerTools.buildLayeredImage {
                name = "bear-trak-otp-${school}";
                tag = "latest";

                contents = [
                  pkgs.jdk21_headless
                  otpGraph
                ];
                config = {
                  Entrypoint = [
                    "${pkgs.jdk21_headless}/bin/java"
                    "${ramFlag}"
                    "-jar"
                    "${otpShaded}"
                  ];
                  Cmd = [
                    "--load"
                    "--serve"
                    "${otpGraph}"
                  ];
                  ExposedPorts = {
                    "8080/tcp" = { };
                  };
                };
              };
          in
          {
            packages = lib.mkMerge [
              (lib.mkIf pkgs.stdenv.isLinux (lib.mapAttrs mkImage schools))

              (lib.concatMapAttrs (
                school: schoolAttrs:
                lib.mapAttrs' (
                  gtfsName: drv: lib.nameValuePair ("gtfs-" + school + "-" + gtfsName) drv
                ) schoolAttrs.gtfsPackages
              ) schools)
            ];
          }
        );
    };
}
