
## 配置优雅退出实例

make step1

make step2

## k8s删除pod时说明1

1. 信号
   - SIGTERM：用于终止程序，也称转终止，因为接收SIGTERM信号的进程可以选择忽略它
   - SIGKILL：用于立即终止，这个信号不能被忽略或者阻止。这是杀死进程的野蛮方式，只能作为最后的手段

2. 当pod进入Terminating逻辑时（删除pod，获取缩容deployment）
   - endpoint资源会立即将对应的pod从svc中摘除
   - 给进程发送SIGTERM信号（默认进程会立即死终止，但优雅退出会hold这个信号）
     - 等待进程结束。
     - 或者等 terminationGracePeriodSeconds 到达上限后，给进程发 SIGKILL 信号量。
   - preStop逻辑：可提供额外优雅退逻辑


## k8s删除pod时说明2

Pod 的生命周期： https://kubernetes.io/zh-cn/docs/concepts/workloads/pods/pod-lifecycle/

1. Pod 生命期
   - Pending
   - Running
   - Succeeded
   - Failed
   - Unknown
   - 说明：当一个 Pod 被删除时，执行一些 kubectl 命令会展示这个 Pod 的状态为 Terminating（终止）。 这个 Terminating 状态并不是 Pod 阶段之一。 Pod 被赋予一个可以体面终止的期限，默认为 30 秒。 你可以使用 --force 参数来强制终止 Pod。

2. 容器状态
   - Waiting
   - Running
   - Terminated ： 如果容器配置了 preStop 回调，则该回调会在容器进入 Terminated 状态之前执行。

3. 容器探针
   - livenessProbe：指示容器是否正在运行。如果存活态探测失败，则 kubelet 会杀死容器
   - readinessProbe：指示容器是否准备好为请求提供服务。如果就绪态探测失败， 端点控制器将从与 Pod 匹配的所有服务的端点列表中删除该 Pod 的 IP 地址
   - startupProbe：指示容器中的应用是否已经启动。如果提供了启动探针，则所有其他探针都会被 禁用，直到此探针成功为止。如果启动探测失败，kubelet 将杀死容器

为容器的生命周期事件设置处理函数：https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/attach-handler-lifecycle-event/
1. postStart：Kubernetes 在容器创建后立即发送 postStart 事件。 然而，postStart 处理函数的调用不保证早于容器的入口点（entrypoint） 的执行。

2. preStop：在容器被终结之前， Kubernetes 将发送一个 preStop 事件。容器可以为每个事件指定一个处理程序。




