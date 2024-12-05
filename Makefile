IMAGE=image.cestc.cn/ies/rmq-tools:v0.0.1

.PHONY: build dev r s clean

build:
	docker build -t ${IMAGE} -f Dockerfile .

dev:
	docker run -it --rm --name rmq-tools ${IMAGE} bash
	
r:
	-docker rm -f rmq-receiver
	docker run --name rmq-receiver ${IMAGE} rmq-tools --role=receiver --topic=test

s:
	-docker rm -f rmq-sender
	docker run --name rmq-sender ${IMAGE} rmq-tools --role=sender --topic=test

clean:
	rm -f rmq-tools	