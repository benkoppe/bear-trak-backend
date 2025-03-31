#!/bin/sh
# Exit immediately if a command exits with a non-zero status.
set -e

# --- Database Initialization ---
# Wait for the database service ('db') to be ready.
# Using the provided check_db_up.sh script.
echo "Waiting for database connection..."
./bin/check_db_up.sh # This script must use 'db' as the hostname and PGPASSWORD env var

echo "Creating database tables..."
./bin/create_tables.sh

echo "Importing GTFS data..."
./bin/import_gtfs.sh

echo "Creating API key..."
./bin/create_api_key.sh

echo "Creating web agency..."
./bin/create_webagency.sh

# --- Optional Steps (from original script comments) ---
# Uncomment if needed
# echo "Importing AVL data..."
# ./import_avl.sh
# echo "Processing AVL data..."
# ./process_avl.sh

# --- Start Main Application ---
echo "Starting TransitClock server..."
# Use exec to replace this script process with the main application process
exec ./bin/start_transitclock.sh
