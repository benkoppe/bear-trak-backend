#!/usr/bin/env bash

echo 'THETRANSITCLOCK DOCKER: Generate SQL to create tables.'

java -cp "$LIB_DIR/Core.jar" org.transitclock.applications.SchemaGenerator -p org.transitclock.db.structs -o "$DB_DIR"
java -cp "$LIB_DIR/Core.jar" org.transitclock.applications.SchemaGenerator -p org.transitclock.db.webstructs -o "$DB_DIR"
