#!/usr/bin/env bash
echo 'THETRANSITCLOCK DOCKER: Import GTFS file.'
# This is to substitute into config file the env values.
find /usr/local/transitclock/config/ -type f -exec sed -i s#"POSTGRES_PORT_5432_TCP_ADDR"#"$POSTGRES_PORT_5432_TCP_ADDR"#g {} \;
find /usr/local/transitclock/config/ -type f -exec sed -i s#"POSTGRES_PORT_5432_TCP_PORT"#"$POSTGRES_PORT_5432_TCP_PORT"#g {} \;
find /usr/local/transitclock/config/ -type f -exec sed -i s#"PGPASSWORD"#"$PGPASSWORD"#g {} \;
find /usr/local/transitclock/config/ -type f -exec sed -i s#"AGENCYNAME"#"$AGENCYNAME"#g {} \;
find /usr/local/transitclock/config/ -type f -exec sed -i s#"TRAK_URL"#"$TRAK_URL"#g {} \;

java -Xmx512m -Dtransitclock.core.agencyId=$AGENCYID -Dtransitclock.configFiles=/usr/local/transitclock/config/transitclock.properties -Dtransitclock.logging.dir=/usr/local/transitclock/logs/ -Dlogback.configurationFile=$TRANSITCLOCK_CORE/transitclock/src/main/resouces/logbackGtfs.xml -cp /usr/local/transitclock/Core.jar org.transitclock.applications.GtfsFileProcessor -gtfsUrl $GTFS_URL -maxTravelTimeSegmentLength 100

psql \
  -h "$POSTGRES_PORT_5432_TCP_ADDR" \
  -p "$POSTGRES_PORT_5432_TCP_PORT" \
  -U postgres \
  -d $AGENCYNAME \
  -c "update activerevisions set configrev=0 where configrev = -1; update activerevisions set traveltimesrev=0 where traveltimesrev = -1;"
