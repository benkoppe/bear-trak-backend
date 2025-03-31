#!/usr/bin/env bash
echo 'THETRANSITCLOCK DOCKER: Start TheTransitClock.'
# This is to substitute into config file the env values
find /usr/local/transitclock/config/ -type f -exec sed -i s#"POSTGRES_PORT_5432_TCP_ADDR"#"$POSTGRES_PORT_5432_TCP_ADDR"#g {} \;
find /usr/local/transitclock/config/ -type f -exec sed -i s#"POSTGRES_PORT_5432_TCP_PORT"#"$POSTGRES_PORT_5432_TCP_PORT"#g {} \;
find /usr/local/transitclock/config/ -type f -exec sed -i s#"PGPASSWORD"#"$PGPASSWORD"#g {} \;
find /usr/local/transitclock/config/ -type f -exec sed -i s#"AGENCYNAME"#"$AGENCYNAME"#g {} \;
find /usr/local/transitclock/config/ -type f -exec sed -i s#"TRAK_URL"#"$TRAK_URL"#g {} \;

rmiregistry &

#set the API as an environment variable so we can set in JSP of template/includes.jsp in the transitime webapp
export APIKEY=$(get_api_key.sh)

# make it so we can also access as a system property in the JSP
export JAVA_OPTS="$JAVA_OPTS -Dtransitclock.apikey=$(get_api_key.sh)"

export JAVA_OPTS="$JAVA_OPTS -Dtransitclock.configFiles=/usr/local/transitclock/config/transitclock.properties"

echo JAVA_OPTS $JAVA_OPTS

/usr/local/tomcat/bin/startup.sh

java -Xss512k -Xms512m -Xmx1024m \
  -Duser.timezone=EST \
  -Dtransitclock.configFiles=/usr/local/transitclock/config/transitclock.properties \
  -Dtransitclock.core.agencyId=$AGENCYID \
  -Dtransitclock.logging.dir=/usr/local/transitclock/logs/ \
  -cp "/usr/local/transitclock/Core.jar:/usr/local/transitclock/lib/*" \
  org.transitclock.applications.Core >/usr/local/transitclock/logs/output.txt &

tail -f /dev/null
