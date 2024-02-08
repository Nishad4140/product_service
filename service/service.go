package service

import (
	"context"
	"fmt"

	"github.com/Nishad4140/product_service/adapter"
	"github.com/Nishad4140/product_service/entities"
	"github.com/Nishad4140/proto_files/pb"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

var (
	Tracer opentracing.Tracer
)

func RetrieveTracer(tr opentracing.Tracer) {
	Tracer = tr
}

type ProductService struct {
	Adapter adapter.AdapterInterface
	pb.UnimplementedProductServiceServer
}

func NewProductService(adapter adapter.AdapterInterface) *ProductService {
	return &ProductService{
		Adapter: adapter,
	}
}

func (product *ProductService) AddPorduct(ctx context.Context, req *pb.AddProductRequest) (*pb.AddProductResponse, error) {

	span := Tracer.StartSpan("add products grpc")

	defer span.Finish()

	if req.Name == "" {
		return nil, fmt.Errorf("name cant be empty")
	}

	reqEntity := entities.Products{
		Name:     req.Name,
		Price:    int(req.Price),
		Quantity: int(req.Quantity),
	}

	res, err := product.Adapter.AddProduct(reqEntity)
	if err != nil {
		return nil, err
	}

	return &pb.AddProductResponse{
		Id:       uint32(res.Id),
		Name:     res.Name,
		Price:    int32(res.Price),
		Quantity: int32(res.Quantity),
	}, nil
}

func (product *ProductService) GetAllProducts(em *empty.Empty, srv pb.ProductService_GetAllProductsServer) error {

	span := Tracer.StartSpan("get all products grpc")
	defer span.Finish()

	products, err := product.Adapter.GetAllProducts()
	if err != nil {
		return err
	}
	for _, prod := range products {
		if err = srv.Send(&pb.AddProductResponse{
			Id:       uint32(prod.Id),
			Name:     prod.Name,
			Price:    int32(prod.Price),
			Quantity: int32(prod.Quantity),
		}); err != nil {
			return err
		}
	}
	return nil
}

type HealthChecker struct {
	grpc_health_v1.UnimplementedHealthServer
}

func (s *HealthChecker) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *HealthChecker) Watch(in *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watching is not supported")
}
