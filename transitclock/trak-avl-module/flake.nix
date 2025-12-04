{
  description = "Custom module for Trak AVL integration";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
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
          ...
        }:
        let
          coreJar = pkgs.fetchurl {
            url = "https://github.com/TheTransitClock/transitime/releases/latest/download/Core.jar";
            sha256 = "sha256-nLUVRraZPcvVr187TPNIFopnrH5vLmxi7R7rBoJfGDo=";
          };
          src = pkgs.runCommand "trak-avl-module-src" { } ''
            mkdir -p $out
            cp ${coreJar} $out/Core.jar
            cp -r ${./.}/* $out/
          '';
          package = pkgs.maven.buildMavenPackage {
            pname = "trak-avl-module";
            version = "1.0";
            inherit src;

            pomFile = ./pom.xml;
            mvnHash = "sha256-vOMQne8HtmT8a5vYRI2SvqovOrD+L1SkYDaw2/g8o24=";

            installPhase = ''
              mkdir -p $out/lib
              cp target/trak-avl-module-1.0-SNAPSHOT.jar $out/lib/
            '';
          };
        in
        {
          packages = {
            default = package;
          };
        };
    };
}
