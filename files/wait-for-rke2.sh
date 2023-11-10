#!/usr/bin/bash
while true; do
  status=$(curl --write-out '%{http_code}' -sk --output /dev/null https://localhost:9345/ping)
  if [ "$status" -eq 200 ]; then
    sleep 2
    break
  fi
  sleep 5
done
