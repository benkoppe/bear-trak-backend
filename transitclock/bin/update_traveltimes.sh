#!/usr/bin/env bash
echo 'THETRANSITCLOCK DOCKER: Update travel times : '+$1+'==>'+$2+'.'
# This is to substitute into config file the env values.
find "$CONFIG_DIR" -type f -exec sed -i s#"POSTGRES_PORT_5432_TCP_ADDR"#"$POSTGRES_PORT_5432_TCP_ADDR"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"POSTGRES_PORT_5432_TCP_PORT"#"$POSTGRES_PORT_5432_TCP_PORT"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"PGPASSWORD"#"$PGPASSWORD"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"AGENCYNAME"#"$AGENCYNAME"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"TRAK_URL"#"$TRAK_URL"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"CONFIG_DIR"#"$CONFIG_DIR"#g {} \;

java -Xmx512m -Xss512k -Dtransitclock.configFiles="$CONFIG_DIR/transitclock.properties" -Dtransitclock.core.agencyId=1 -Dtransitclock.logging.dir=/usr/local/transitclock/logs/ -cp "$LIB_DIR/Core.jar" org.transitclock.applications.UpdateTravelTimes $1 $2
