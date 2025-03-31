#!/usr/bin/env bash
echo 'THETRANSITCLOCK DOCKER: Create API key.'
# This is to substitute into config file the env values
find /usr/local/transitclock/config/ -type f -exec sed -i s#"POSTGRES_PORT_5432_TCP_ADDR"#"$POSTGRES_PORT_5432_TCP_ADDR"#g {} \;
find /usr/local/transitclock/config/ -type f -exec sed -i s#"POSTGRES_PORT_5432_TCP_PORT"#"$POSTGRES_PORT_5432_TCP_PORT"#g {} \;
find /usr/local/transitclock/config/ -type f -exec sed -i s#"PGPASSWORD"#"$PGPASSWORD"#g {} \;
find /usr/local/transitclock/config/ -type f -exec sed -i s#"AGENCYNAME"#"$AGENCYNAME"#g {} \;
find /usr/local/transitclock/config/ -type f -exec sed -i s#"TRAK_URL"#"$TRAK_URL"#g {} \;

java -cp /usr/local/transitclock/Core.jar org.transitclock.applications.CreateAPIKey -c "/usr/local/transitclock/config/transitclock.properties" -d "foo" -e "og.crudden@gmail.com" -n "Sean Og Crudden" -p "123456" -u "http://www.transitclock.org"
