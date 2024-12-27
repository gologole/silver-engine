run-db:
	docker run -d \
      --name postgres-library \
      -e POSTGRES_USER=username \
      -e POSTGRES_PASSWORD=password \
      -e POSTGRES_DB=library \
      -p 5432:5432 \
      postgres:latest
run frontend:
	npm run dev
