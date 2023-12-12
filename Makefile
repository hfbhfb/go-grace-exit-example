all: step1 step2

step1: #编译server二进制 和打包镜像
	export GOOS=linux; export GOARCH=amd64 && go build -o server grace-server/grace-server.go 
	docker build -t mygrace:0.1 .

step2: #部署负载
	kubectl delete -f grace-app.yaml
	kubectl apply -f grace-app.yaml

curlexample:
	curl http://192.168.88.246:30579/
	curl http://192.168.88.246:30580/1.txt
