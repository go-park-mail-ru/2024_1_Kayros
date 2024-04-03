version: '3'

services:
  minio:
    image: minio/minio:latest
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - minio-data:/data
      - /etc/nginx/sslcert/server.crt:/root/.minio/certs/public.crt:ro
      - /etc/nginx/sslcert/server.key:/root/.minio/certs/private.key:ro
    environment:
      MINIO_ROOT_USER: kayrosteam
      MINIO_ROOT_PASSWORD: resto123
      MINIO_SERVER_URL: "https://resto-go.ru:9000"
      MINIO_BROWSER_REDIRECT_URL: "https://resto-go.ru/minio-admin"
    hostname: resto-go.ru
    command: server --console-address :9001 /data

volumes:
    minio-data: {}