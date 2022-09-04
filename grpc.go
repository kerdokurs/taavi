package main

import (
	"context"
	"log"
	"net/http"

	pb "kerdo.dev/taavi/pkg"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
)

const Port = "9000"

func (t *Taavi) runGRPC() {
	server := TaaviServer{Taavi: t}
	s := JobServer{
		TaaviServer: &server,
	}

	mux := runtime.NewServeMux()
	err := pb.RegisterJobServiceHandlerServer(context.Background(), mux, &s)
	if err != nil {
		log.Fatalf("Could not register job service: %v\n", err)
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})
	cmux := c.Handler(mux)

	log.Printf("Starting backend server on port %s\n", Port)

	if err := http.ListenAndServe("0.0.0.0:"+Port, cmux); err != nil {
		log.Fatalf("Could not start backend server: %v\n", err)
	}
}

type TaaviServer struct {
	Taavi *Taavi
}

type JobServer struct {
	pb.UnimplementedJobServiceServer
	*TaaviServer
}

func (j *JobServer) CreateNewJob(ctx context.Context, request *pb.NewJobRequest) (*pb.JobResponse, error) {
	info := CronInfo{
		TimeString: request.CronTime,
		Type:       CronType(request.Type),
		ServerId:   request.ServerId,
		ChannelId:  request.ChannelId,
		Message:    request.Message,
	}
	tx := j.Taavi.db.Create(&info)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &pb.JobResponse{
		Id:               int32(info.ID),
		CronTime:         info.TimeString,
		Type:             int32(info.Type),
		NextScheduleTime: int32(info.ScheduledAt.UnixMilli()),
		ServerId:         info.ServerId,
		ChannelId:        info.ChannelId,
		Message:          info.Message,
	}, nil
}

func (j *JobServer) GetAllJobs(ctx context.Context, request *pb.AllJobsRequest) (*pb.AllJobsResponse, error) {
	var infos []CronInfo
	tx := j.Taavi.db.Find(&infos)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &pb.AllJobsResponse{
		Jobs: Map(infos, func(info CronInfo) *pb.JobResponse {
			return &pb.JobResponse{
				Id:               int32(info.ID),
				CronTime:         info.TimeString,
				Type:             int32(info.Type),
				NextScheduleTime: int32(info.ScheduledAt.UnixMilli()),
				ServerId:         info.ServerId,
				ChannelId:        info.ChannelId,
				Message:          info.Message,
			}
		}),
	}, nil
}

func (j *JobServer) CronSync(ctx context.Context, request *pb.Empty) (*pb.Empty, error) {
	j.Taavi.CronSync(false)
	return &pb.Empty{}, nil
}
