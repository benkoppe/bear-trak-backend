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
          system,
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

          chrome-headless-amd64 = pkgs.dockerTools.pullImage {
            imageName = "docker.io/chromedp/headless-shell";
            imageDigest = "sha256:de5b057849de96955de7662023420a46355abfbb24d57aa01282ec7c811aacab";
            finalImageTag = "latest";
            sha256 = "sha256-SYPr3y9n79RRJEv0Tms9/SlVFzfYkzU9EDIRNXfwq6s=";
            os = "linux";
            arch = "amd64";
          };
          chrome-headless-arm64 = pkgs.dockerTools.pullImage {
            imageName = "docker.io/chromedp/headless-shell";
            imageDigest = "sha256:de5b057849de96955de7662023420a46355abfbb24d57aa01282ec7c811aacab";
            finalImageTag = "latest";
            sha256 = "sha256-0ND5n5Q+yZ4ACCISxzvgFpe+K4nISITXyf/NYWkpncc=";
            os = "linux";
            arch = "arm64";
          };
          containerImage = pkgs.dockerTools.buildLayeredImage {
            name = "bear-trak-go";
            tag = "latest";

            fromImage = if pkgs.stdenv.isAarch64 then chrome-headless-arm64 else chrome-headless-amd64;
            architecture = system;
            # Default is 100, so this ensures this image gets its own layer(s)
            # after being merged with the base image.
            maxLayers = 120;
            contents = [
              pkgs.cacert
            ];
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

          packages = {
            default = server;
            inherit containerImage;
          };
        };
    };
}
