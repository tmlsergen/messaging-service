#!/bin/bash
# Declare variables.
API_URL=http://localhost:80

# Run the api tests.
echo
echo "Running end-to-end testing..."
echo
echo "Testing GET route '/api/v1/messages'..."
curl $API_URL/api/v1/messages; echo

echo
echo "Testing POST route '/api/v1/messages/cron' start action..."
curl -X POST -H 'Content-Type: application/json' -d '{"action": "stop"}' $API_URL/api/v1/messages/cron; echo
echo "Testing POST route '/api/v1/messages/cron' stop action..."
curl -X POST -H 'Content-Type: application/json' -d '{"action": "start"}' $API_URL/api/v1/messages/cron; echo

# Finish the testing.
echo
echo "Finished testing the application!"