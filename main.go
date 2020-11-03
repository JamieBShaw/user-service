package main

import (
	"context"
	"flag"
	"github.com/JamieBShaw/user-service/protob"
	"github.com/JamieBShaw/user-service/repository/postgres"
	"github.com/JamieBShaw/user-service/service"
	internalGrpc "github.com/JamieBShaw/user-service/transport/grpc"
	internalhttp "github.com/JamieBShaw/user-service/transport/http"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	googlegrpc "google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	log            = logrus.New()
	router         = mux.NewRouter()
	port           = os.Getenv("PORT")
	DB_USER        = os.Getenv("PGUSER")
	DB_PASSWORD     = os.Getenv("PGPASSWORD")
	DB_NAME = os.Getenv("PGDATABASE")
	grpc           = flag.Bool("grpc", false, "service will use grpc (http2) as the transport layer")
)

func main() {
	flag.Parse()

	dbConnection := pg.Connect(&pg.Options{
		User:     DB_USER,
		Password: DB_PASSWORD,
		Database: DB_NAME,
		OnConnect: func(ctx context.Context, cn *pg.Conn) error {
			_, err := cn.Exec("CREATE TABLE IF NOT EXISTS users (" +
				"id bigserial primary key," +
				"username varchar(40) unique,"+
				"admin bool,"+
				"created_at timestamp default now() not null," +
				"updated_at timestamp default now() not null" +
				");")
			if err != nil {
				log.Errorf("error creating table on service startup: %v", err)
				return err
			}
			return nil
		},
	})

	cc, err := googlegrpc.Dial("0.0.0.0:8081", googlegrpc.WithInsecure())
	if err != nil {
		log.Fatalf("error getting connection grpc client: %v", err)
	}

	log.Println("Starting GPRC Client dial on port 8081")

	authClient := protob.NewAuthServiceClient(cc)

	defer dbConnection.Close()

	repo := postgres.NewRepository(log, dbConnection)
	userService := service.NewUserService(repo)

	if *grpc {
		log.Infof("Starting GRPC User Service running on port: %v", port)

		lis, err := net.Listen("tcp", "0.0.0.0:"+port)
		if err != nil {
			log.Fatal("Failed to listen", err)
		}

		s := googlegrpc.NewServer()
		srv := internalGrpc.NewGrpcServer(userService)
		protob.RegisterUserServiceServer(s, srv)

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

	} else {
		handler := internalhttp.NewHttpHandler(userService, router, authClient)

		srv := &http.Server{
			Addr:         "0.0.0.0:" + port,
			Handler:      handler,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		go func() {
			log.Infof("Starting HTTP User Service running on port: %v", port)
			if err := srv.ListenAndServe(); err != nil {
				log.Println(err)
			}
		}()


		c := make(chan os.Signal, 1)
		// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
		// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
		signal.Notify(c, os.Interrupt)

		// Block until we receive our signal.
		<-c

		// Register a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		srv.Shutdown(ctx)
		// Optionally, you could run srv.Shutdown in a goroutine and block on
		// <-ctx.Done() if your application should wait for other services
		// to finalize based on context cancellation.
		log.Println("shutting down")
		os.Exit(0)
	}
}
