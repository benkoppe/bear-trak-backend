#!/usr/bin/env bash
echo 'THETRANSITCLOCK DOCKER: Update travel times : '+$1+'==>'+$2+'.'
# This is to substitute into config file the env values.
find /usr/local/transitclock/config/ -type f -exec sed -i s#"POSTGRES_PORT_5432_TCP_ADDR"#"$POSTGRES_PORT_5432_TCP_ADDR"#g {} \;
find /usr/local/transitclock/config/ -type f -exec sed -i s#"POSTGRES_PORT_5432_TCP_PORT"#"$POSTGRES_PORT_5432_TCP_PORT"#g {} \;
find /usr/local/transitclock/config/ -type f -exec sed -i s#"PGPASSWORD"#"$PGPASSWORD"#g {} \;
find /usr/local/transitclock/config/ -type f -exec sed -i s#"AGENCYNAME"#"$AGENCYNAME"#g {} \;
find /usr/local/transitclock/config/ -type f -exec sed -i s#"TRAK_URL"#"$TRAK_URL"#g {} \;

java -Xmx2048m -Xss12m -Dtransitclock.configFiles=/usr/local/transitclock/config/transitclock.properties -Dtransitclock.core.agencyId=1 -Dtransitclock.logging.dir=/usr/local/transitclock/logs/ -cp /usr/local/transitclock/Core.jar org.transitclock.applications.UpdateTravelTimes $1 $2
