IMAGE_REPOSITORY_NAME := tabelog-sort
CONTAINER_NAME := tabelog
ECR_REGISTRY := ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o ./build/tabelogbot ./cmd/
	zip -j ./build/tabelogbot ./build/tabelogbot
	docker build -t ${IMAGE_REPOSITORY_NAME} ./build

.PHONY: build-raspberry
build-raspberry:
	GOOS=linux GOARCH=arm GOARM=6 go build -tags netgo -installsuffix netgo -ldflags '-extldflags "-static"' -o ./build/tabelogbot ./cmd/

.PHONY: run
run:
	docker run -it -d --rm -p 9000:8080 --name "${CONTAINER_NAME}" ${IMAGE_REPOSITORY_NAME}

.PHONY: login
login:
	aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com

.PHONY: push
push:
	docker tag ${IMAGE_REPOSITORY_NAME} ${ECR_REGISTRY}/${IMAGE_REPOSITORY_NAME}
	docker push ${ECR_REGISTRY}/${IMAGE_REPOSITORY_NAME}
