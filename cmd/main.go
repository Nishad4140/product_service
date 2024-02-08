package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/Nishad4140/product_service/db"
	"github.com/Nishad4140/product_service/initializer"
	"github.com/Nishad4140/product_service/service"
	servicediscoveryconsul "github.com/Nishad4140/product_service/servicediscovery_consul"
	"github.com/Nishad4140/proto_files/pb"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err.Error())
	}

	addr := os.Getenv("DATABASE_ADDR")

	DB, err := db.InitDB(addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	services := initializer.Initialize(DB)

	server := grpc.NewServer()

	pb.RegisterProductServiceServer(server, services)

	lis, err := net.Listen("tcp", ":3000")

	if err != nil {
		log.Fatalf("Failed to listen on 3000 : %v", err)
	}

	log.Printf("Product server is listening on port 3000")

	go func() {
		time.Sleep(5 * time.Second)

		servicediscoveryconsul.RegisterService()
	}()

	healthService := &service.HealthChecker{}

	grpc_health_v1.RegisterHealthServer(server, healthService)

	tracer, closer := initTracer()

	defer closer.Close()

	service.RetrieveTracer(tracer)

	if err = server.Serve(lis); err != nil {
		log.Fatalf("Failed to connect on port 3000 : %v", err)
	}

}

func initTracer() (tracer opentracing.Tracer, closer io.Closer) {
	jaegerEndpoint := "http://localhost:14268/api/traces"

	cfg := &config.Configuration{
		ServiceName: "product-service",
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: jaegerEndpoint,
		},
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("updated")
	return
}
