package surveymarkup

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"project-survey-generator/internal/configuration"
	"project-survey-generator/internal/dbcache"
	"project-survey-generator/internal/surveymarkup/minifier"
	"project-survey-generator/internal/surveymarkup/model/pb"
	"project-survey-generator/internal/utils"
)

type Server struct {
	pb.UnimplementedSurveyMarkupGeneratorServer

	srvGeneratorAddr string
	dbRepo           *dbcache.Repo
	generator        *Generator
	minifier         *minifier.Service
}

func NewServer(config *configuration.AppConfiguration, dbRepo *dbcache.Repo, minifier *minifier.Service, generator *Generator) *Server {
	return &Server{
		srvGeneratorAddr: config.SurveyProceederUrl,
		dbRepo:           dbRepo,
		generator:        generator,
		minifier:         minifier,
	}
}

func (s *Server) GenerateMarkup(ctx context.Context, rq *pb.GenerateMarkupRequest) (*pb.GenerateMarkupResponse, error) {
	unit := s.dbRepo.GetUnitById(int(rq.UnitId))
	if unit == nil {
		return nil, status.Error(codes.NotFound, "")
	}

	surveys := s.dbRepo.GetUnitSurveysWithIds(unit.Id, utils.ConvertInts[int](rq.SurveyIds))
	if len(surveys) != len(rq.SurveyIds) {
		return nil, status.Error(codes.NotFound, "")
	}

	markup, err := s.generator.Generate(unit, surveys, rq.Language)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	markup = s.minifier.Minify(markup)

	return &pb.GenerateMarkupResponse{Markup: markup}, nil
}
