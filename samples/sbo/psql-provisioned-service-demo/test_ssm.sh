#!/usr/bin/env bash

host=$1

echo "Health Check:"
curl http://$host:8080/healthcheck

clear_text="Trabajo Libre"

echo "Wrapping Clear Text: $clear_text"
cipher=$(curl -s http://localhost:8080/wrap -H "Content-Type: application/json" -d "{\"key\": \"$clear_text\"}" | jq -r .cipher)

echo "Wrapped Cipher: $cipher"

clear_text_back=$(curl -s http://localhost:8080/unwrap -H "Content-Type: application/json" -d "{\"cipher\": \"$cipher\"}" | jq -r .key)

echo "Unwrapped Cipher: $clear_text_back"

if [[ "$clear_text" == "$clear_text_back" ]]; then
    echo "Unwrapped cipher matches original clear text. Success!!"
else
    echo "Unwrapped cipher does not match original clear text. Failure :-("
fi
