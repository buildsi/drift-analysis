#! /bin/bash

if [ "$DRIFTSERVER_ACCESS_KEY_ID" != "" ]; then
    echo "dbs:" > /etc/litestream.yml
    echo "  - path: /app/points.db" >> /etc/litestream.yml
    echo "    replicas:" >> /etc/litestream.yml
    echo "      - url: $DRIFTSERVER_REPLICA_URL" >> /etc/litestream.yml
    echo "        retention: 24h" >> /etc/litestream.yml
    echo "        sync-interval: 15m" >> /etc/litestream.yml
    export LITESTREAM_SECRET_ACCESS_KEY=$DRIFTSERVER_SECRET_ACCESS_KEY
    litestream restore -if-replica-exists -v -o points.db "$DRIFTSERVER_REPLICA_URL"
    litestream replicate -exec /app/drift-server
else 
    /app/drift-server
fi

