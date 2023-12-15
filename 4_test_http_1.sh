docker run -it --net=host curlimages/curl curl --location '127.0.0.1:9999/send' \
--header 'Content-Type: application/json' \
--data '{
"device_id":"00000000-0000-1111-2222-334455667781",
"id":"ffffffff-0000-1111-2222-334455667788",
"kind":1,
"message":"Text message"
}'
