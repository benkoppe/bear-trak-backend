{
  description = "Custom module for Trak AVL integration";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    transitclock.url = "path:../transitclock-flake";
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
          inputs',
          ...
        }:
        let
          src = pkgs.symlinkJoin {
            name = "trak-avl-module-src";
            paths = [
              ./.
              "${inputs'.transitclock.packages.default}/lib"
            ];
          };
          package = pkgs.maven.buildMavenPackage rec {
            pname = "trak-avl-module";
            version = "1.0";
            inherit src;

            pomFile = ./pom.xml;
            mvnHash = "sha256-vOMQne8HtmT8a5vYRI2SvqovOrD+L1SkYDaw2/g8o24=";

            installPhase = ''
              mkdir -p $out/lib
              cp target/trak-avl-module-${version}-SNAPSHOT.jar $out/lib/trak-avl-module.jar
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
