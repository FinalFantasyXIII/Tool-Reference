## RPCX 文档

+ 快速起步
    + 服务端
    ```
    import (
	"context"
	"github.com/smallnest/rpcx/server")

    //服务参数
    type Args struct {
	    A int
	    B int
    }
    type Reply struct {
    	C int
    }
    type Arith int

    //服务方法定义实现
    func (t *Arith) Mul(ctx context.Context, args *Args, reply *Reply) error{
    	reply.C = args.A * args.B
    	return nil
    }

    func (t * Arith) Add(ctx context.Context, args *Args, reply *Reply) error{
    	reply.C = args.A + args.B
    	return nil
    }

    func (t * Arith) Sub(ctx context.Context, args *Args, reply *Reply) error{
    	reply.C = args.A - args.B
    	return nil
    }

    func (t * Arith) Div(ctx context.Context, args *Args, reply *Reply) error{
    	reply.C = args.A / args.B
    	return nil
    }

    func main(){
    	s := server.NewServer()
    	s.RegisterName("Arith",new(Arith),"")   //服务注册
    	s.Serve("tcp",":12345")                 //服务绑定
    }
    ```
    + 客户端
    ```
    import (
	"context"
	"github.com/smallnest/rpcx/client"
	"log")

    type Args struct {
    	A int
    	B int
    }
    type Reply struct {
    	C int
    }

    func main()  {
        //点对点绑定对应的服务
    	d := client.NewPeer2PeerDiscovery("tcp@"+"127.0.0.1:12345", "")

        //指定需要访问的服务
    	xclient := client.NewXClient("Arith", client.Failtry,   client.RandomSelect, d, client.DefaultOption)
    	defer xclient.Close()

    	//客户端入参要传地址
    	args := &Args{
    		A: 10,
    		B: 20,
    	}
    	reply := &Reply{}

        //同步调用
    	err := xclient.Call(context.Background(), "Add", args, reply)
    	if err != nil {
    		log.Fatalf("failed to call: %v", err)
    	}
    	log.Printf("%d * %d = %d", args.A, args.B, reply.C)

        //异步调用
    	result,err := xclient.Go(context.Background(),"Mul",args,reply,nil)
    	if err != nil{
    		log.Fatalf("failed to call: %v", err)
    	}
    	r := <-result.Done
    	if r.Error != nil{
    		log.Fatalf("failed to call: %v", r.Error)
    	}else{
    		log.Printf("%d * %d = %d", args.A, args.B, reply.C)
    	}
    }
    ```

+ 服务端控制选项
    ```
    func NewServer(options ...OptionFn) *Server
    func (s *Server) Close() error
    func (s *Server) RegisterOnShutdown(f func())
    func (s *Server) Serve(network, address string) (err error)
    func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request)
    ```

+ 客户端控制选项
    ```
    Client 基本函数:
        func (client *Client) Call(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}) error
        func (client *Client) Close() error
        func (c *Client) Connect(network, address string) error
        func (client *Client) Go(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call
        func (client *Client) IsClosing() bool
        func (client *Client) IsShutdown() bool

    XClient:
        type XClient interface {
            SetPlugins(plugins PluginContainer)
            ConfigGeoSelector(latitude, longitude float64)
            Auth(auth string)
            Go(ctx context.Context, serviceMethod string, args interface{}, reply interface{}, done chan *Call) (*Call, error)
            Call(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
            Broadcast(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
            Fork(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
            Close() error
        }
    ```
    + 服务发现 
        + Peer to Peer: 客户端直连每个服务节点。 the client connects the single service directly. It acts like the client type.
        + Peer to Multiple: 客户端可以连接多个服务。服务可以被编程式配置。
        + Zookeeper: 通过 zookeeper 寻找服务。
        + Etcd: 通过 etcd 寻找服务。
        + Consul: 通过 consul 寻找服务。
        + mDNS: 通过 mDNS 寻找服务（支持本地服务发现）。
        + In process: 在同一进程寻找服务。客户端通过进程调用服务，不走TCP或UDP，方便调试使用。

    + 服务治理(失败模式与负载均衡)
        + rpcx 支持 故障模式:(对应NewXClient函数第二个参数)
            + Failfast：如果调用失败，立即返回错误
            + Failover：选择其他节点，直到达到最大重试次数
            + Failtry：选择相同节点并重试，直到达到最大重试次数
        + rpcx 提供了许多选择器:(对应NewXClient函数第三个参数)
            + Random选择器 
        + client.NewXClient("Arith",client.Failtry,client.RandomSelect,d, client.DefaultOption)
    + 广播与群发
        + Broadcast 表示向所有服务器发送请求，只有所有服务器正确返回时才会成功。
        + Fork 表示向所有服务器发送请求，只要任意一台服务器正确返回就成功。此时FailMode 和 SelectMode的设置是无效的。

+ 注册中心
    + Peer2Peer
        ```
        点对点是最简单的一种注册中心的方式，事实上没有注册中心，客户端直接得到唯一的服务器的地址，连接服务
        d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
        xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
        defer xclient.Close()

        注意:rpcx使用network @ Host: port格式表示一项服务。在network 可以 tcp ， http ，unix ，quic或kcp。该Host可以所主机名或IP地址。
        NewXClient必须使用服务名称作为第一个参数，然后使用failmode，selector，discovery和其他选项
        ```
    
    + MultipleServers
        ```
        上面的方式只能访问一台服务器，假设我们有固定的几台服务器提供相同的服务，我们可以采用这种方式。
        d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr1}, {Key: *addr2}})
        xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
        defer xclient.Close()

        你必须在MultipleServersDiscovery 中设置服务信息和元数据。如果添加或删除了某些服务，你可以调用MultipleServersDiscovery.Update来动态更新服务。
        func (d *MultipleServersDiscovery) Update(pairs []*KVPair)
        ```

    + ZooKeeper

    + Etcd

    + Consul

    + mDNS

    + Inprocess
        ```
        这个Registry用于进程内的测试。 在开发过程中，可能不能直接连接线上的服务器直接测试，而是写一些mock程序作为服务，这个时候就可以使用这个registry, 测试通过在部署的时候再换成相应的其它registry.

        在这种情况下， client和server并不会走TCP或者UDP协议，而是直接进程内方法调用,所以服务器代码是和client代码在一起的。

        func InProcesRPC(){
        	//注册服务
        	s := server.NewServer()
        	c := client.InprocessClient
        	s.Plugins.Add(c)
        	s.RegisterName("Arith",new(Arith),"")
        	go func() {
        		s.Serve("tcp",":12345") 		//s.Server 会挂起服务，导致阻塞，因此用协程
        	}()

        	//设置客户端
        	d := client.NewInprocessDiscovery()
        	xclient := client.NewXClient("Arith",client.Failtry,client.RandomSelect,d,client.DefaultOption)
        	defer xclient.Close()

        	args := &Args{
        		A: 10,
        		B: 20,
        	}

        	//调用服务
        	for i:=0; i<10000;i++{
        		reply := &Reply{}
        		n := i%4
        		switch n {
        		case 0:
        			err := xclient.Call(context.Background(),"Add",args,reply)
        			if err != nil{
        				log.Fatalf("failed to call: %v", err)
        			}
        			log.Printf("%d + %d = %d", args.A, args.B, reply.C)
        			break;
        		case 1:
        			err := xclient.Call(context.Background(),"Sub",args,reply)
        			if err != nil{
        				log.Fatalf("failed to call: %v", err)
        			}
        			log.Printf("%d - %d = %d", args.A, args.B, reply.C)
        			break;
        		case 2:
        			err := xclient.Call(context.Background(),"Mul",args,reply)
        			if err != nil{
        				log.Fatalf("failed to call: %v", err)
        			}
        			log.Printf("%d * %d = %d", args.A, args.B, reply.C)
        			break;
        		case 3:
        			err := xclient.Call(context.Background(),"Div",args,reply)
        			if err != nil{
        				log.Fatalf("failed to call: %v", err)
        			}
        			log.Printf("%d / %d = %d", args.A, args.B, reply.C)
        			break;
        		}
        	}
        }
        ```

+ 特性
    + 编解码(json/protoBuff/XML/...)
        ```
        Option参数可以设置自己的编码器
        func NewXClient(servicePath string, failMode FailMode, selectMode  SelectMode, discovery ServiceDiscovery, option Option)
        ```
        + SerializeNone
            + 这种编解码器不会对数据进行编解码，并且要求数据是 []byte 类型的数据。
        + JSON
            + 对性能要求不是非常高的场景，可以使用这种编解码。
        + Protobuf
        + MsgPack
            + 默认的编解码器
        + 定制编解码器

    +  失败模式 (NewClient函数的第二个参数)
        + Failfast
            + 在这种模式下， 一旦调用一个节点失败， rpcx立即会返回错误。
        + Failover
            + 在这种模式下, rpcx如果遇到错误，它会尝试调用另外一个节点， 直到服务节点能正常返回信息，或者达到最大的重试次数。 重试测试Retries在参数Option中设置， 缺省设置为3。
        + Failtry
            + 在这种模式下， rpcx如果调用一个节点的服务出现错误， 它也会尝试，但是还是选择这个节点进行重试， 直到节点正常返回数据或者达到最大重试次数。
        + Failbackup
            + 在这种模式下， 如果服务节点在一定的时间内不返回结果， rpcx客户端会发送相同的请求到另外一个节点， 只要这两个节点有一个返回， rpcx就算调用成功。这个设定的时间配置在 Option.BackupLatency 参数中。

    + Fork模式(XClient中一种远程调用模式,Call,Go)
        ```
        Fork 表示向所有服务器发送请求，只要任意一台服务器正确返回就成功。此时FailMode 和 SelectMode的设置是无效的。
        func main() {
            ……

            xclient := client.NewXClient("Arith", client.Failover, client.RoundRobin, d, client.DefaultOption)
            defer xclient.Close()

            args := &example.Args{
                A: 10,
                B: 20,
            }

            for {
                reply := &example.Reply{}
                err := xclient.Fork(context.Background(), "Mul", args, reply)
                if err != nil {
                    log.Fatalf("failed to call: %v", err)
                }

                log.Printf("%d * %d = %d", args.A, args.B, reply.C)
                time.Sleep(1e9)
            }
        }
        ```

    + 广播模式(同fork call go)
        ```
        可以将一个请求发送到这个服务的所有节点。 如果所有的节点都正常返回，没有错误的话， Broadcast将返回其中的一个节点的返回结果。 如果有节点返回错误的话，Broadcast将返回这些错误信息中的一个。
        func main() {
            ……

            xclient := client.NewXClient("Arith", client.Failover, client.RoundRobin, d, client.DefaultOption)
            defer xclient.Close()

            args := &example.Args{
                A: 10,
                B: 20,
            }

            for {
                reply := &example.Reply{}
                err := xclient.Broadcast(context.Background(), "Mul", args, reply)
                if err != nil {
                    log.Fatalf("failed to call: %v", err)
                }

                log.Printf("%d * %d = %d", args.A, args.B, reply.C)
                time.Sleep(1e9)
            }
        }
        ```

    + 路由 (NewClient函数的第三个参数)
        + 随机
            ```
            NewXClient(1,2, client.RandomSelect ,4,5)
            从配置的节点中随机选择一个节点。
            最简单，但是有时候单个节点的负载比较重。这是因为随机数只能保证在大量的请求下路由的比较均匀，并不能保证在很短的时间内负载是均匀的。
            ```
        + 轮询
            ```
            NewXClient(1,2, client.RoundRobin ,4,5)
            使用轮询的方式，依次调用节点，能保证每个节点都均匀的被访问。在节点的服务能力都差不多的时候适用。
            ```
            + WeightedRoundRobin
                ```
                NewXClient(1,2, client.WeightedRoundRobin ,4,5)
                使用Nginx 平滑的基于权重的轮询算法。
                比如如果三个节点a、b、c的权重是{ 5, 1, 1 }, 这个算法的调用顺序是 { a, a, b, a, c, a, a }, 相比较 { c, b, a, a, a, a, a }, 虽然权重都一样，但是前者更好，不至于在一段时间内将请求都发送给a。
                ```
        + 网络质量优先(Ping)
            ```
            NewXClient(1,2, client.WeightedICMP ,4,5)
            首先客户端会基于ping(ICMP)探测各个节点的网络质量，越短的ping时间，这个节点的权重也就越高。但是，我们也会保证网络较差的节点也有被调用的机会。
            ```
        + 一致性哈希
            ```
            NewXClient(1,2, client.ConsistentHash ,4,5)
            使用 JumpConsistentHash 选择节点， 相同的servicePath, serviceMethod 和 参数会路由到同一个节点上。 JumpConsistentHash 是一个快速计算一致性哈希的算法，但是有一个缺陷是它不能删除节点，如果删除节点，路由就不准确了，所以在节点有变动的时候它会重新计算一致性哈希。
            ```
        + 地理位置优先
            ```
            NewXClient(1,2, client.ConsistentHash ,4,5)
            如果我们希望的是客户端会优先选择离它最新的节点， 比如在同一个机房。 如果客户端在北京， 服务在上海和美国硅谷，那么我们优先选择上海的机房。

            它要求服务在注册的时候要设置它所在的地理经纬度。
            如果两个服务的节点的经纬度是一样的， rpcx会随机选择一个。
            比必须使用下面的方法配置客户端的经纬度信息：
                func (c *xClient) ConfigGeoSelector(latitude, longitude float64)
            ```
        
        + 定制路由规则
            ```
            type Selector interface {
                Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string
                UpdateServer(servers map[string]string)
            }
            ```
    + NewXClient最后的参数Options
        + 超时
            > 超时机制可以保护服务调用陷入无限的等待之中。超时定义了服务的最长等待时间，如果在给定的    时间没有相应，服务调用就进入下一个状态，或者重试、或者立即返回错误。
            + 配置client
                ```
                option := client.DefaultOption
	            option.ReadTimeout = 10 * time.Second
                xclient := client.NewXClient(1, 2, 3, 4, option)
                ```
            + 配置context
                ```
                func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
                ```
        + 元数据
            ```
            客户端和服务器端可以互相传递元数据。
            元数据不是服务请求和服务响应的业务数据，而是一些辅助性的数据。
            元数据是一个键值队的列表，键和值都是字符串， 类似 http.Header。
            ```
        + 心跳
            ```
            option := client.DefaultOption
            option.Heartbeat = true
            option.HeartbeatInterval = time.Second
            ```
        + 分组
            ```
            如果你为服务设置了设置group， 只有在这个group的客户端才能访问这些服务
            ```
        + 服务状态(state)
            ```
            state 是另外一个元数据。 如果你在元数据中设置了state=inactive, 客户端将不能访问这些服   务，即使这些服务是"活"着的。
            你可以使用临时禁用一些服务，而不是杀掉它们， 这样就实现了服务的降级。
            ```
        
+ 插件
    + Metrics插件

    + 限流插件

    + 别名插件

    + 身份认证
        + 服务端
            ```
            func main() {
                flag.Parse()

                s := server.NewServer()
                s.RegisterName("Arith", new(example.Arith), "")
                s.AuthFunc = auth
                s.Serve("reuseport", *addr)
            }

            func auth(ctx context.Context, req *protocol.Message, token string) error {
            
                if token == "bearer tGzv3JOkF0XG5Qx2TlKWIA" {
                    return nil
                }

                return errors.New("invalid token")
            }
            ```
        + 客户端
            ```
            func main() {
                flag.Parse()

                d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")

                option := client.DefaultOption
                option.ReadTimeout = 10 * time.Second

                xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, option)
                defer xclient.Close()

                //xclient.Auth("bearer tGzv3JOkF0XG5Qx2TlKWIA")
                xclient.Auth("bearer abcdefg1234567")

                args := &example.Args{
                    A: 10,
                    B: 20,
                }

                reply := &example.Reply{}
                ctx := context.WithValue(context.Background(), share.ReqMetaDataKey, make(map[string]string))
                err := xclient.Call(ctx, "Mul", args, reply)
                if err != nil {
                    log.Fatalf("failed to call: %v", err)
                }

                log.Printf("%d * %d = %d", args.A, args.B, reply.C)

            }
            ```
    + 