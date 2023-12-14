docker build --tag 'websocket_funout' .
docker network create mynetwork
docker network connect mynetwork myserver
docker network connect mynetwork myclient
docker run --detach 'websocket_funout'