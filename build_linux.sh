CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
scp ./Journal root@140.82.56.114:/root
scp -r ./conf root@140.82.56.114:/root
# 建立ssh信任
