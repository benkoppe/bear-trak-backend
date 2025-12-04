#!/usr/bin/env bash
echo 'THETRANSITCLOCK DOCKER: Set Config.'

find "$CONFIG_DIR" -type f -exec sed -i s#"POSTGRES_PORT_5432_TCP_ADDR"#"$POSTGRES_PORT_5432_TCP_ADDR"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"POSTGRES_PORT_5432_TCP_PORT"#"$POSTGRES_PORT_5432_TCP_PORT"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"PGPASSWORD"#"$PGPASSWORD"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"AGENCYNAME"#"$AGENCYNAME"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"TRAK_URL"#"$TRAK_URL"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"CONFIG_DIR"#"$CONFIG_DIR"#g {} \;
