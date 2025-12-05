{
  description = "Description for the project";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    transitclock.url = "path:./transitclock-flake";
    trak-avl-module.url = "path:./trak-avl-module";
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
          runtimeEnv = pkgs.stdenv.mkDerivation {
            pname = "transitclock-runtime-env";
            version = "1.0";

            buildInputs = with pkgs; [
              jq
              postgresql
              openjdk8
              coreutils
            ];

            src = ./.;
            phases = [
              "installPhase"
              "fixupPhase"
            ];

            installPhase =
              let
                transit = inputs'.transitclock.packages.default;
                trak = inputs'.trak-avl-module.packages.default;
              in
              ''
                mkdir -p $out/lib
                mkdir -p $out/etc/{config,logs,cache,data,test,bin}
                mkdir -p $out/var/tomcat
                mkdir -p $out/bin

                cp ${transit}/lib/Core.jar $out/lib/
                cp ${trak}/lib/trak-avl-module.jar $out/lib/

                cp -r $src/config/* $out/etc/config

                cp -r $src/bin/* $out/bin/
                chmod +x $out/bin/*.sh

                mkdir -p $out/var/tomcat/webapps
                cp ${transit}/lib/api.war $out/var/tomcat/webapps/
                cp ${transit}/lib/web.war $out/var/tomcat/webapps/
              '';
          };
          entrypoint = pkgs.writeShellApplication {
            name = "start-transitclock";
            runtimeInputs = [
              pkgs.openjdk8
              pkgs.tomcat9
              pkgs.jq
              pkgs.postgresql
              pkgs.coreutils
              pkgs.findutils
              pkgs.gnused
            ];
            runtimeEnv = {
              LIB_DIR = "${runtimeEnv}/lib";
              DB_DIR = "/usr/local/transitclock/db";

              READONLY_CONFIG_DIR = "${runtimeEnv}/etc/config";
              CONFIG_DIR = "/tmp/transitclock-config";

              WEBAPPS_DIR = "${runtimeEnv}/var/tomcat/webapps";
              CATALINA_BASE = "/tmp/tomcat";

              # GTFS_URL = "file://${umich-gtfs}";
              SOURCE_DIR = "${inputs'.transitclock.packages.source}";
            };
            text = ''
              ln -sf ${runtimeEnv}/bin ./bin

              mkdir -p $CONFIG_DIR
              cp -r $READONLY_CONFIG_DIR/* $CONFIG_DIR/

              mkdir -p $CATALINA_BASE
              mkdir -p $CATALINA_BASE/logs
              cp -r ${pkgs.tomcat9}/conf $CATALINA_BASE/conf
              cp -r $WEBAPPS_DIR $CATALINA_BASE/webapps
            ''
            + builtins.readFile ./entrypoint.sh;
          };
          containerImage = pkgs.dockerTools.buildLayeredImage {
            name = "trak-transitclock";
            tag = "latest";

            config = {
              Cmd = [ "${entrypoint}/bin/start-transitclock" ];
              ExposedPorts = {
                "8080/tcp" = { };
              };
              Env = [
                "CATALINA_HOME=${pkgs.tomcat9}"
                "JAVA_HOME=${pkgs.openjdk8}"
              ];
            };
          };
        in
        {
          packages = {
            default = containerImage;
            inherit runtimeEnv containerImage entrypoint;
          };
          apps = {
            entrypoint.program = "${entrypoint}/bin/start-transitclock";
          };
        };
    };
}
