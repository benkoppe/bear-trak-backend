#!/usr/bin/env bash
echo 'THETRANSITCLOCK DOCKER: Create Tables'

java -cp "$LIB_DIR/Core.jar" org.transitclock.applications.SchemaGenerator -p org.transitclock.db.structs -o "$DB_DIR"
java -cp "$LIB_DIR/Core.jar" org.transitclock.applications.SchemaGenerator -p org.transitclock.db.webstructs -o "$DB_DIR"

createdb -h "$POSTGRES_PORT_5432_TCP_ADDR" -p "$POSTGRES_PORT_5432_TCP_PORT" -U postgres $AGENCYNAME
psql \
  -h "$POSTGRES_PORT_5432_TCP_ADDR" \
  -p "$POSTGRES_PORT_5432_TCP_PORT" \
  -U postgres \
  -d $AGENCYNAME \
  -f "$DB_DIR/ddl_postgres_org_transitclock_db_structs.sql"
psql \
  -h "$POSTGRES_PORT_5432_TCP_ADDR" \
  -p "$POSTGRES_PORT_5432_TCP_PORT" \
  -U postgres \
  -d $AGENCYNAME \
  -f "$DB_DIR/ddl_postgres_org_transitclock_db_webstructs.sql"
