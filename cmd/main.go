package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"project-survey-generator/internal/configuration"
	"project-survey-generator/internal/services"
	"project-survey-generator/internal/surveymarkup"
	"project-survey-generator/internal/surveymarkup/model/pb"
	"syscall"
)

func main() {
	parser := configuration.NewParser()
	config, err := parser.Parse("appsettings.json")
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	serviceProvider := services.NewProvider(config)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.SurveyGeneratorPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterSurveyMarkupGeneratorServer(s, surveymarkup.NewServer(config, serviceProvider.GetDbRepo(), serviceProvider.GetMinifier(), serviceProvider.GetGenerator()))

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
