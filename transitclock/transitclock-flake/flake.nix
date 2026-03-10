{
  description = "Description for the project";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    transitime = {
      url = "github:benkoppe/transitime";
      flake = false;
    };
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
          package = pkgs.maven.buildMavenPackage {
            pname = "transitclock";
            version = "2.0.0";
            src = inputs.transitime;

            mvnHash = "sha256-5hVY7takDYYywMPddheZ5HViAOVX5iyU19kgGlYVB0Y=";

            installPhase = ''
              mkdir -p $out/lib
              cp transitclock/target/Core.jar $out/lib
              cp transitclockApi/target/api.war $out/lib
              cp transitclockWebapp/target/web.war $out/lib
            '';
            doCheck = false;
          };
        in
        {
          packages = {
            default = package;
            source = pkgs.stdenv.mkDerivation {
              name = "transitime-src";
              src = inputs.transitime;
              installPhase = ''
                mkdir -p $out
                cp -r . $out
              '';
            };
          };
        };
    };
}
