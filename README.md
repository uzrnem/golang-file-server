# golang-file-server
Golang server for file upload and download

## Installation with Docker

Attach your `files` Directory to Container's `/app/files` Directory, files directory will act as file server
Sample `docker-compose.yml`
```sh
version: '3.7'

services:
  file-server:
    image: uzrnem/file-server:0.1
    container_name: file_server
    environment:
      SERVER_PORT: 9060
    volumes:
      - $HOME/uzrnem/golang-file-server/files:/app/files
    ports:
      - 9060:9060
```