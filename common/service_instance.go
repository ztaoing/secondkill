package common

type ServiceInstance struct {
	Host      string
	Port      int
	Weight    int //权重
	CurWeight int //当前权重
	GrpcPort  int //grpc port
}
