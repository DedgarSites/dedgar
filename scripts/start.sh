#!/bin/bash -e

if [ "$PAUSE_ON_START" = "true" ] ; then
  echo
  echo "This container's startup has been paused indefinitely because PAUSE_ON_START has been set."
  echo
  while true; do
    sleep 10    
  done
fi

if [ "$LEGO_CERT" = "true" ] ; then
  echo "Checking if our LE certificates need to be renewed with Lego"
  /usr/local/bin/lego --tls=true \
                      --tls.port=":$TLS_PORT" \
                      --email="$CERT_EMAIL" \
                      --domains="$APP_DOMAIN" \
                      --path="$CERT_PATH" \
                      --filename="$CERT_FILENAME" \
                      --server="$CERT_SERVER" \
                      --accept-tos run
  sleep 10
fi

echo
echo "Starting web server with the certs in $CERT_PATH"
/go/bin/"$APP_NAME"
