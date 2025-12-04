#!/usr/bin/env bash
echo 'THETRANSITCLOCK DOCKER: Create WebAgency.'
# This is to substitute into config file the env values
find "$CONFIG_DIR" -type f -exec sed -i s#"POSTGRES_PORT_5432_TCP_ADDR"#"$POSTGRES_PORT_5432_TCP_ADDR"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"POSTGRES_PORT_5432_TCP_PORT"#"$POSTGRES_PORT_5432_TCP_PORT"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"PGPASSWORD"#"$PGPASSWORD"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"AGENCYNAME"#"$AGENCYNAME"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"TRAK_URL"#"$TRAK_URL"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"CONFIG_DIR"#"$CONFIG_DIR"#g {} \;

echo "HELLO!!!!"
echo java -Dtransitclock.db.dbName=$AGENCYNAME -Dtransitclock.hibernate.configFile="$CONFIG_DIR/hibernate.cfg.xml" -Dtransitclock.db.dbHost=$POSTGRES_PORT_5432_TCP_ADDR:$POSTGRES_PORT_5432_TCP_PORT -Dtransitclock.db.dbUserName=postgres -Dtransitclock.db.dbPassword=$PGPASSWORD -Dtransitclock.db.dbType=postgresql -cp "$LIB_DIR/Core.jar" org.transitclock.db.webstructs.WebAgency $AGENCYID 0.0.0.0 $AGENCYNAME postgresql $POSTGRES_PORT_5432_TCP_ADDR postgres $PGPASSWORD
java -Dtransitclock.db.dbName=$AGENCYNAME -Dtransitclock.hibernate.configFile="$CONFIG_DIR/hibernate.cfg.xml" -Dtransitclock.db.dbHost=$POSTGRES_PORT_5432_TCP_ADDR:$POSTGRES_PORT_5432_TCP_PORT -Dtransitclock.db.dbUserName=postgres -Dtransitclock.db.dbPassword=$PGPASSWORD -Dtransitclock.db.dbType=postgresql -cp "$LIB_DIR/Core.jar" org.transitclock.db.webstructs.WebAgency $AGENCYID 0.0.0.0 $AGENCYNAME postgresql $POSTGRES_PORT_5432_TCP_ADDR postgres $PGPASSWORD
