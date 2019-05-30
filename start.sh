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
  echo "running Lego to check if our certificates need to be renewed"
  #/usr/local/bin/lego --tls=":8443" --email="contact@shinobu.ninja" --domains="www.shinobu.ninja" --path="/cert/lego" --filename="shinobu" --accept-tos --exclude="http-01" run
  /usr/local/bin/lego --tls=":8443" --email="dedgarsite@shinobu.ninja" --domains="www.dedgar.com" --path="/cert/lego" --filename="dedgar" --accept-tos --exclude="http-01" run
  sleep 10
fi

echo
echo "running Echo with the certs in /cert"
/usr/local/bin/dedgar
