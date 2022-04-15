#! /bin/bash

if [ "$DRIFTSERVER_S3_ACCESSKEY" != "" ]; then
    echo "dbs:" > /etc/litestream.yml
    echo "  - path: /app/data/points.db" >> /etc/litestream.yml
    echo "    replicas:" >> /etc/litestream.yml
    echo "      - type: s3" >> /etc/litestream.yml
    echo "        bucket: $DRIFTSERVER_S3_BUCKET" >> /etc/litestream.yml
    echo "        access-key-id: $DRIFTSERVER_S3_ACCESSKEY" >> /etc/litestream.yml
    echo "        secret-access-key: $DRIFTSERVER_S3_SECRETKEY" >> /etc/litestream.yml
    echo "        path: database/points.db" >> /etc/litestream.yml
    echo "        region: $DRIFTSERVER_S3_REGION" >> /etc/litestream.yml
    echo "        retention: 24h" >> /etc/litestream.yml
    echo "        sync-interval: 15m" >> /etc/litestream.yml
    
    litestream restore -if-replica-exists -v /app/data/points.db 
    litestream replicate -exec /app/drift-server
else 
    /app/drift-server
fi

