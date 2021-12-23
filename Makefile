run-server:
	docker run -d -p 32801:32801 \
		-e STOMPER_HOSTNAME=0.0.0.0 \
		-e STOMPER_TOPICS=/channel/main \
		-e STOMPER_TCPDEADLINE=0 \
	ghcr.io/tydar/stomper:main
