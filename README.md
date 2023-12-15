# Prerequisites:

One need to have Docker installer.

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
