# ① 构建阶段：只要编译器和源码
 FROM golang:1.26.5-bookworm AS build
 WORKDIR /src

 # ② 先只拷依赖清单，单独一层 → 代码改了不用重下依赖
 COPY go.mod go.sum ./
 RUN go mod download

 # ③ 再拷源码编译
 COPY . .
 RUN CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o /out/server .

 # ④ 运行阶段：只带二进制
 FROM debian:bookworm-slim
 RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
     && rm -rf /var/lib/apt/lists/* \
     && useradd -m -u 10001 app \
     && mkdir -p /data && chown app:app /data

 WORKDIR /app
 COPY --from=build /out/server .

 ENV DB_PATH=/data/api.db PORT=8090
 VOLUME /data
 EXPOSE 8090
 USER app
 CMD ["/app/server"]