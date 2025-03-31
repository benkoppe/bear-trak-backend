#!/usr/bin/env bash
echo 'THETRANSITCLOCK DOCKER: Create Tables'

java -cp /usr/local/transitclock/Core.jar org.transitclock.applications.SchemaGenerator -p org.transitclock.db.structs -o /usr/local/transitclock/db
java -cp /usr/local/transitclock/Core.jar org.transitclock.applications.SchemaGenerator -p org.transitclock.db.webstructs -o /usr/local/transitclock/db

createdb -h "$POSTGRES_PORT_5432_TCP_ADDR" -p "$POSTGRES_PORT_5432_TCP_PORT" -U postgres $AGENCYNAME
psql \
	-h "$POSTGRES_PORT_5432_TCP_ADDR" \
	-p "$POSTGRES_PORT_5432_TCP_PORT" \
	-U postgres \
	-d $AGENCYNAME \
	-f /usr/local/transitclock/db/ddl_postgres_org_transitclock_db_structs.sql
psql \
	-h "$POSTGRES_PORT_5432_TCP_ADDR" \
	-p "$POSTGRES_PORT_5432_TCP_PORT" \
	-U postgres \
	-d $AGENCYNAME \
	-f /usr/local/transitclock/db/ddl_postgres_org_transitclock_db_webstructs.sql
