build:
	@echo "Building binary file" && \
	go build -o logwatcher main.go && \
	echo "Done"