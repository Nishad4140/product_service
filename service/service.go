package service

import (
	"context"
	"fmt"

	"github.com/Nishad4140/product_service/adapter"
	"github.com/Nishad4140/product_service/entities"
	"github.com/Nishad4140/proto_files/pb"
	"github.com/opentracing/opentracing-go"
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
