#!/usr/bin/env bash
echo 'THETRANSITCLOCK DOCKER: Start TheTransitClock.'
# This is to substitute into config file the env values
find "$CONFIG_DIR" -type f -exec sed -i s#"POSTGRES_PORT_5432_TCP_ADDR"#"$POSTGRES_PORT_5432_TCP_ADDR"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"POSTGRES_PORT_5432_TCP_PORT"#"$POSTGRES_PORT_5432_TCP_PORT"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"PGPASSWORD"#"$PGPASSWORD"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"AGENCYNAME"#"$AGENCYNAME"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"TRAK_URL"#"$TRAK_URL"#g {} \;
find "$CONFIG_DIR" -type f -exec sed -i s#"CONFIG_DIR"#"$CONFIG_DIR"#g {} \;

rmiregistry &

#set the API as an environment variable so we can set in JSP of template/includes.jsp in the transitime webapp
export APIKEY=$(get_api_key.sh)

# make it so we can also access as a system property in the JSP
export JAVA_OPTS="$JAVA_OPTS -Dtransitclock.apikey=$(get_api_key.sh)"

export JAVA_OPTS="$JAVA_OPTS -Dtransitclock.configFiles=${CONFIG_DIR}/transitclock.properties"

echo JAVA_OPTS $JAVA_OPTS

"$CATALINA_HOME/bin/startup.sh"

java -Xss512k -Xms512m -Xmx1024m \
  -Duser.timezone=America/New_York \
  -Dtransitclock.configFiles="$CONFIG_DIR/transitclock.properties" \
  -Dtransitclock.core.agencyId=$AGENCYID \
  -Dtransitclock.logging.dir=/usr/local/transitclock/logs/ \
  -cp "${LIB_DIR}/Core.jar:${LIB_DIR}/*" \
  org.transitclock.applications.Core >/usr/local/transitclock/logs/output.txt &

tail -f /dev/null
