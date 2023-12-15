# Description

A simple implementation of a concurrent broadcast server that dispatches a message to the connected clients via the websocket protocol.
ALso, the server has HTTP endpoint to receive a message and sends it out to all  (or selected) clients.

# Prerequisites:

One need to have Docker installer.

# Architecture 
![architecture](./docs/architecture.png)

# Startup instructions:

## 1. Build source code to docker image:

```bash
sh 1_build.sh
```

## 2. Run server (image in docker)

```bash
sh 2_run_server.sh
```

## 3. Run 3 clients 

```bash
sh 3_run_client_1.sh
```

```bash
sh 3_run_client_2.sh
```

```bash
sh 3_run_client_3.sh
```

## 4. Test using HTTP

Send message to the first device:
```bash
sh 4_test_http_1.sh
```

Send message to all devices:
```bash
sh 4_test_http_2.sh
```
