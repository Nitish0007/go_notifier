#!/bin/bash

DB_URL="postgres://postgres:12345678@notifier_db:5432/notifier_dev_db?sslmode=disable"

echo "Cleaning migration dirty state...."

psql $DB_URL <<EOF
UPDATE schema_migrations SET dirty = false WHERE dirty = true;
EOF

echo "Done, now you can run migrations again"


# command for rollback

# docker run --rm \
#   --network go_notifier_default \
#   -v $(pwd)/db/migrations:/migrations \
#   migrate/migrate \
#   -path=/migrations \
#   -database "postgres://postgres:12345678@notifier_db:5432/notifier_dev_db?sslmode=disable" \
#   force 2

# TO ROLLBACK LAST SUCCESSFULL MIGRATION
# docker run --rm \
#   --network go_notifier_default \
#   -v $(pwd)/db/migrations:/migrations \
#   migrate/migrate \
#   -path=/migrations \
#   -database "postgres://postgres:12345678@notifier_db:5432/notifier_dev_db?sslmode=disable" \
#   down 1
