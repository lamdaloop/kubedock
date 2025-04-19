#!/bin/bash
docker run --name kubedock-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=kubedock \
  -p 5432:5432 \
  -d postgres:15

