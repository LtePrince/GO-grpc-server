jaeger-up:
	sudo docker-compose -f ./docker/docker-compose.yml up -d

jaeger-down:
	sudo docker-compose -f ./docker/docker-compose.yml down