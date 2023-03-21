
docker.start.components:
  docker-compose up -d --remove-orphans postgres;

docker.stop:
  docker-compose down;
  