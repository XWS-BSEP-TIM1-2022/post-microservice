package startup

import (
	"fmt"
	postService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/post"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/token"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
	otgo "github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"post-microservice/application"
	"post-microservice/infrastructure/api"
	"post-microservice/infrastructure/persistance"
	"post-microservice/model"
	"post-microservice/startup/config"
)

type Server struct {
	config     *config.Config
	tracer     otgo.Tracer
	closer     io.Closer
	jwtManager *token.JwtManager
}

func NewServer(config *config.Config) *Server {
	tracer, closer := tracer.Init(config.PostServiceName)
	otgo.SetGlobalTracer(tracer)
	jwtManager := token.NewJwtManagerDislinkt(config.ExpiresIn)
	return &Server{
		config:     config,
		tracer:     tracer,
		closer:     closer,
		jwtManager: jwtManager,
	}
}

func (server *Server) GetTracer() otgo.Tracer {
	return server.tracer
}

func (server *Server) GetCloser() io.Closer {
	return server.closer
}

func (server *Server) Start() {
	mongoClient := server.initMongoClient()

	postStore := server.initPostStore(mongoClient)
	postService := server.initPostService(postStore)

	commentStore := server.initCommentStore(mongoClient)
	commentService := server.initCommentService(commentStore)

	postHandler := server.initPostHandler(postService, commentService)
	server.startGrpcServer(postHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistance.GetClient(server.config.PostDBHost, server.config.PostDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) startGrpcServer(postHandler *api.PostHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	log.Println(fmt.Sprintf("started grpc server on localhost:%s", server.config.Port))
	postService.RegisterPostServiceServer(grpcServer, postHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func (server *Server) initPostStore(client *mongo.Client) model.PostStore {
	store := persistance.NewPostMongoDBStore(client)
	return store
}

func (server *Server) initPostService(store model.PostStore) *application.PostService {
	return application.NewPostService(store)
}

func (server *Server) initPostHandler(postService *application.PostService, commentService *application.CommentService) *api.PostHandler {
	return api.NewPostHandler(postService, commentService)
}

func (server *Server) initCommentStore(client *mongo.Client) model.CommentStore {
	store := persistance.NewCommentMongoDBStore(client)
	return store
}

func (server *Server) initCommentService(store model.CommentStore) *application.CommentService {
	return application.NewCommentService(store)
}
