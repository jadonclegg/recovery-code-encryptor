BINARY_NAME=recovery-code-encryptor

build: frontend
	CGO_ENABLED=0 GOOS=linux go build

run: docker
	podman run -it --rm --name recovery-code-encryptor-test -p 8080:8080 recovery-code-encryptor:latest

docker: build
	podman build -f DockerFile --tag recovery-code-encryptor

runfs: docker
	podman run -it --rm --name recovery-code-encryptor-test -p 8080:8080 -v ./web2/dist:/app/dist recovery-code-encryptor:latest /app/recovery-code-encryptor --fs

debugfs: docker
	podman run -it --rm --name recovery-code-encryptor-test -p 8080:8080 -v ./web2/dist:/app/dist recovery-code-encryptor:latest sh

frontend:
	cd web2 && ng build

frontend-watch:
	cd web2 && ng build --watch -c development

.PSEUDO js: dist/js/main.js