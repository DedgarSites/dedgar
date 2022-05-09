#!/bin/bash -e

if [ "$PAUSE_ON_START" = "true" ] ; then
  echo
  echo "This container's startup has been paused indefinitely because PAUSE_ON_START has been set."
  echo
  while true; do
    sleep 10    
  done
fi

while true; do
  echo
  echo "Starting web server:"
  "$SITE_PATH"/"$APP_NAME"
  echo "Sleeping for a 10 seconds before starting the web server again with new certs."
  sleep 10
done
