export PGPASSWORD=transitclock

docker stop transitclock-db
docker stop transitclock-server-instance

docker rm transitclock-db
docker rm transitclock-server-instance

docker rmi transitclock-server

docker build --no-cache -t transitclock-server --platform linux/amd64 \
  --build-arg TRANSITCLOCK_PROPERTIES="config/transitclock.properties" \
  --build-arg AGENCYID="1" \
  --build-arg AGENCYNAME="UMICH" \
  --build-arg GTFS_URL="https://webapps.fo.umich.edu/transit_uploads/google_transit.zip" .

docker run --platform linux/amd64 --name transitclock-db -p 5432:5432 -e POSTGRES_PASSWORD=$PGPASSWORD -d postgres:9.6.3

docker run --platform linux/amd64 --name transitclock-server-instance --rm --link transitclock-db:postgres -e PGPASSWORD=$PGPASSWORD -v transitclock-logs:/usr/local/transitclock/logs/ transitclock-server check_db_up.sh

docker run --platform linux/amd64 --name transitclock-server-instance --rm --link transitclock-db:postgres -e PGPASSWORD=$PGPASSWORD -v transitclock-logs:/usr/local/transitclock/logs/ transitclock-server create_tables.sh

docker run --platform linux/amd64 --name transitclock-server-instance --rm --link transitclock-db:postgres -e PGPASSWORD=$PGPASSWORD -v transitclock-logs:/usr/local/transitclock/logs/ transitclock-server import_gtfs.sh

docker run --platform linux/amd64 --name transitclock-server-instance --rm --link transitclock-db:postgres -e PGPASSWORD=$PGPASSWORD -v transitclock-logs:/usr/local/transitclock/logs/ transitclock-server create_api_key.sh

docker run --platform linux/amd64 --name transitclock-server-instance --rm --link transitclock-db:postgres -e PGPASSWORD=$PGPASSWORD -v transitclock-logs:/usr/local/transitclock/logs/ transitclock-server create_webagency.sh

#docker run --name transitclock-server-instance --rm --link transitclock-db:postgres -e PGPASSWORD=$PGPASSWORD transitclock-server ./import_avl.sh

#docker run --name transitclock-server-instance --rm --link transitclock-db:postgres -e PGPASSWORD=$PGPASSWORD transitclock-server ./process_avl.sh

docker run --platform linux/amd64 --name transitclock-server-instance --rm --link transitclock-db:postgres -e PGPASSWORD=$PGPASSWORD -v transitclock-logs:/usr/local/transitclock/logs/ -v transitclock-cache:/usr/local/transitclock/cache/ -p 8080:8080 transitclock-server start_transitclock.sh
