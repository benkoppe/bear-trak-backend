#!/usr/bin/env bash
echo 'THETRANSITCLOCK DOCKER: Create API key.'
# This is to substitute into config file the env values
find "$CONFIG_DIR" -type f -exec sed -i s#"POSTGRES_PORT_5432_TCP_ADDR"#"$POSTGRES_PORT_5432_TCP_ADDR"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"POSTGRES_PORT_5432_TCP_PORT"#"$POSTGRES_PORT_5432_TCP_PORT"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"PGPASSWORD"#"$PGPASSWORD"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"AGENCYNAME"#"$AGENCYNAME"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"TRAK_URL"#"$TRAK_URL"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"CONFIG_DIR"#"$CONFIG_DIR"#g {} \;

java -cp "$LIB_DIR/Core.jar" org.transitclock.applications.CreateAPIKey -c "$CONFIG_DIR/transitclock.properties" -d "foo" -e "og.crudden@gmail.com" -n "Sean Og Crudden" -p "123456" -u "http://www.transitclock.org"
