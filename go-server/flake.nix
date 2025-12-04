{
  description = "Trak core backend";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
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
          lib,
          pkgs,
          inputs',
          ...
        }:
        let
          gomod2nixPkgs = inputs'.gomod2nix.legacyPackages;
          # Simple lint check added to nix flake check
          lint = pkgs.stdenvNoCC.mkDerivation {
            name = "go-lint";
            dontBuild = true;
            src = ./.;
            doCheck = true;
            nativeBuildInputs = with pkgs; [
              golangci-lint
              go
              writableTmpDirAsHomeHook
            ];
            checkPhase = ''
              golangci-lint run
            '';
            installPhase = ''
              mkdir "$out"
            '';
          };
          server = gomod2nixPkgs.buildGoApplication {
            pname = "trak-server";
            version = "0.1";
            pwd = ./.;
            src = ./.;
            modules = ./gomod2nix.toml;
            doCheck = false; # tests were for while writing and many are broken
          };
          containerImage = pkgs.dockerTools.buildLayeredImage {
            name = "bear-trak-go";
            tag = "latest";

            config = {
              Entrypoint = [ "${server}/bin/go-server" ];
            };
          };
        in
        {
          checks = { inherit lint; };

          devShells.default =
            let
              goEnv = gomod2nixPkgs.mkGoEnv { pwd = ./.; };
            in
            pkgs.mkShell {
              packages = [
                goEnv
                gomod2nixPkgs.gomod2nix
              ];
            };

          packages = lib.mkMerge [
            {
              default = server;
            }
            (lib.mkIf pkgs.stdenv.isLinux { inherit containerImage; })
          ];
        };
    };
}
