#!/usr/bin/env bash
echo 'THETRANSITCLOCK DOCKER: Check if database is runnng.'
RET=1
SUCCESS=0
until [ "$RET" -eq "$SUCCESS" ]; do

  psql -h "$POSTGRES_PORT_5432_TCP_ADDR" -p "$POSTGRES_PORT_5432_TCP_PORT" -U postgres -c "SELECT EXTRACT(DAY FROM TIMESTAMP '2001-02-16 20:38:40');"
  RET="$?"

  if [ "$RET" -ne "$SUCCESS" ]; then
    echo 'Database is not running.'
    sleep 10
  fi
done
echo 'THETRANSITCLOCK DOCKER: Database is now running.'
