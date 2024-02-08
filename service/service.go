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

func (product *ProductService) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.AddProductResponse, error) {

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

func (product *ProductService) GetProduct(ctx context.Context, req *pb.GetProductByID) (*pb.AddProductResponse, error) {

	span := Tracer.StartSpan("get products by using id")
	defer span.Finish()

	res, err := product.Adapter.GetProduct(uint(req.Id))
	if err != nil {
		return nil, err
	}

	if res.Name == "" {
		return nil, fmt.Errorf("there is no product in the given id")
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

func (product *ProductService) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest) (*pb.AddProductResponse, error) {

	span := Tracer.StartSpan("update quantity of product")
	defer span.Finish()

	var res *pb.AddProductResponse

	if req.Increase {

		result, err := product.Adapter.IncrementStock(uint(req.Id), int(req.Quantity))
		if err != nil {
			return nil, err
		}

		res = &pb.AddProductResponse{
			Id:       uint32(result.Id),
			Name:     result.Name,
			Price:    int32(result.Price),
			Quantity: int32(result.Quantity),
		}
	} else {
		result, err := product.Adapter.DecrementStock(uint(req.Id), int(req.Quantity))
		if err != nil {
			return nil, err
		}
		res = &pb.AddProductResponse{
			Id:       uint32(result.Id),
			Name:     result.Name,
			Price:    int32(result.Price),
			Quantity: int32(result.Quantity),
		}
	}

	return res, nil
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
