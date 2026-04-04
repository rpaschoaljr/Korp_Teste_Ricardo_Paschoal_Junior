#!/bin/bash
set -e

# Só executa se a variável SEED_DATABASE for 'true'
if [ "$SEED_DATABASE" = "true" ]; then
    echo "SEED_DATABASE is set to true. Running super_seed.sql..."
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/03-super_seed.sql
else
    echo "SEED_DATABASE is not set to true. Skipping super_seed.sql."
fi
