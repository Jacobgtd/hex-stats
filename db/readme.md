docker-compose up -d
Get-Content .\schema.sql -Raw | docker exec -i postgres-local psql -U username -d mydb
docker compose down -v