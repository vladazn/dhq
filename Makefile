SERVICE_TAG = vladazn/dhq/service
SWAGGER_TAG = vladazn/dhq/swagger
API_TAG = vladazn/dhq/api
NGINX_TAG = vladazn/dhq/nginx
VERSION = test


all: generate swagger

generate:
	rm -rf proto/gen
	cd proto && buf generate

.PHONY: swagger
swagger:
	rm -f swagger/ui/swagger.json
	cp proto/gen/openapiv2/proto/storage/storage.swagger.json swagger/ui/swagger.json

docker_service:
	docker build -t $(SERVICE_TAG):$(VERSION) -f service/docker/Dockerfile .

docker_api:
	docker build -t $(API_TAG):$(VERSION) -f api/docker/Dockerfile .

docker_swagger:
	docker build -t $(SWAGGER_TAG):$(VERSION) -f swagger/docker/Dockerfile .

docker: docker_service docker_swagger docker_api
