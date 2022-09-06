package main

import (
	"context"
	"fmt"
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
		StreamId:   request.StreamId,
		TopicId:    request.TopicId,
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
		StreamId:         info.StreamId,
		TopicId:          info.TopicId,
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
				StreamId:         info.StreamId,
				TopicId:          info.TopicId,
				Message:          info.Message,
			}
		}),
	}, nil
}

func (j *JobServer) CronSync(ctx context.Context, request *pb.Empty) (*pb.Empty, error) {
	j.Taavi.CronSync(false)
	return &pb.Empty{}, nil
}

func (j *JobServer) GetScheduledJobs(ctx context.Context, req *pb.Empty) (*pb.AllJobsResponse, error) {
	jobs := j.Taavi.ScheduledJobs()
	fmt.Println(jobs)
	fmt.Println(len(jobs))
	return &pb.AllJobsResponse{}, nil
}

func (j *JobServer) DeleteJob(ctx context.Context, req *pb.JobIdRequest) (*pb.Empty, error) {
	id := int(req.Id)
	if id == 0 {
		return nil, fmt.Errorf("invalid id")
	}

	info := CronInfo{
		SchedulerId: id,
	}
	if tx := j.Taavi.db.Find(&info); tx.Error != nil {
		return nil, tx.Error
	}

	if err := j.Taavi.cancelByUUID(info.ScheduleId); err != nil {
		return nil, err
	}
	if tx := j.Taavi.db.Delete(&info); tx.Error != nil {
		return nil, tx.Error
	}

	return &pb.Empty{}, nil
}
