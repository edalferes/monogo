#!/bin/bash
set -e

# Espera o banco subir
until pg_isready -U "$POSTGRES_USER"; do
  sleep 1
done

# Habilita a extensão uuid-ossp
psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c 'CREATE EXTENSION IF NOT EXISTS "uuid-ossp";'
