all: step1 step2

step1: #编译server二进制 和打包镜像
	export GOOS=linux; export GOARCH=amd64 && go build -o server grace-server/grace-server.go 
	docker build -t mygrace:0.1 .

step2: #部署负载
	-kubectl delete -f grace-app.yaml
	kubectl apply -f grace-app.yaml


step2envmove: #触发滚动升级
	kubectl apply -f grace-app-envmore.yaml

step2prestop: #preStop 和 SIGTERM信号 的顺序关系？
	kubectl apply -f grace-app-prestop.yaml


curlexample:
	curl http://192.168.88.246:30579/

clean:
	kubectl delete -f grace-app.yaml
