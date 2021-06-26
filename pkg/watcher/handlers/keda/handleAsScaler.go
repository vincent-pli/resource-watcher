package handlers

import (

	// cloudevents "github.com/cloudevents/sdk-go"
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/vincent-pli/resource-watcher/pkg/watcher/handlers/keda/externalscaler"

	"github.com/go-logr/logr"
	"google.golang.org/grpc"
	"k8s.io/client-go/tools/cache"
)

type KedaHandler struct {
	TargetInstance *corev1.ObjectReference,
	PollingInterval uint,
	Log logr.Logger,
}

var _ cache.ResourceEventHandler = (*KedaHandler)(nil)

func New(targetInstance *corev1.ObjectReference, pollingInterval uint, log logr.Logger) (*KedaHandler, error) {
	//create svc, scaledObject before really start
	


	return &KedaHandler {
		TargetInstance: targetInstance,
		PollingInterval: pollingInterval,
		Log: log,
	}

	return nil, nil
}

func (c *KedaHandler) Start() error {
	grpcServer := grpc.NewServer()
	lis, _ := net.Listen("tcp", ":6000")
	pb.RegisterExternalScalerServer(grpcServer, &ExternalScaler{})

	fmt.Println("listenting on :6000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

func (c *KedaHandler) OnAdd(obj interface{}) {

}

func (c *KedaHandler) OnUpdate(oldObj, newObj interface{}) {

}

func (c *KedaHandler) OnDelete(obj interface{}) {

}

func (e *ExternalScaler) IsActive(ctx context.Context, scaledObject *pb.ScaledObjectRef) (*pb.IsActiveResponse, error) {
	// longitude := scaledObject.ScalerMetadata["longitude"]
	// latitude := scaledObject.ScalerMetadata["latitude"]

	// if len(longitude) == 0 || len(latitude) == 0 {
	// 	return nil, status.Error(codes.InvalidArgument, "longitude and latitude must be specified")
	// }

	// startTime := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	// endTime := time.Now().Format("2006-01-02")
	// radiusKM := 500
	// query := fmt.Sprintf("format=geojson&starttime=%s&endtime=%s&longitude=%s&latitude=%s&maxradiuskm=%d", startTime, endTime, longitude, latitude, radiusKM)

	// resp, err := http.Get(fmt.Sprintf("https://earthquake.usgs.gov/fdsnws/event/1/query?%s", query))
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	// payload := USGSResponse{}
	// err = json.Unmarshal(body, &payload)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	// count := 0
	// for _, f := range payload.Features {
	// 	if f.Properties.Mag > 1.0 {
	// 		count++
	// 	}
	// }

	// return &pb.IsActiveResponse{
	// 	Result: count > 2,
	// }, nil
}

func (e *ExternalScaler) GetMetricSpec(context.Context, *pb.ScaledObjectRef) (*pb.GetMetricSpecResponse, error) {
	// return &pb.GetMetricSpecResponse{
	// 	MetricSpecs: []*pb.MetricSpec{{
	// 		MetricName: "eqThreshold",
	// 		TargetSize: 10,
	// 	}},
	// }, nil
}

func (e *ExternalScaler) GetMetrics(_ context.Context, metricRequest *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	// longitude := metricRequest.ScaledObjectRef.ScalerMetadata["longitude"]
	// latitude := metricRequest.ScaledObjectRef.ScalerMetadata["latitude"]

	// if len(longitude) == 0 || len(latitude) == 0 {
	// 	return nil, status.Error(codes.InvalidArgument, "longitude and latitude must be specified")
	// }

	// earthquakeCount, err := getEarthQuakeCount(longitude, latitude, 1.0)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	// return &pb.GetMetricsResponse{
	// 	MetricValues: []*pb.MetricValue{{
	// 		MetricName:  "earthquakeThreshold",
	// 		MetricValue: int64(earthquakeCount),
	// 	}},
	// }, nil
}

func getEarthQuakeCount(longitude, latitude string, magThreshold float64) (int, error) {
	return 0, nil
}

func (e *ExternalScaler) StreamIsActive(scaledObject *pb.ScaledObjectRef, epsServer pb.ExternalScaler_StreamIsActiveServer) error {
	// longitude := scaledObject.ScalerMetadata["longitude"]
	// latitude := scaledObject.ScalerMetadata["latitude"]

	// if len(longitude) == 0 || len(latitude) == 0 {
	// 	return status.Error(codes.InvalidArgument, "longitude and latitude must be specified")
	// }

	// for {
	// 	select {
	// 	case <-epsServer.Context().Done():
	// 		// call cancelled
	// 		return nil
	// 	case <-time.Tick(time.Hour * 1):
	// 		earthquakeCount, err := getEarthQuakeCount(longitude, latitude, 1.0)
	// 		if err != nil {
	// 			// log error
	// 		} else if earthquakeCount > 2 {
	// 			err = epsServer.Send(&pb.IsActiveResponse{
	// 				Result: true,
	// 			})
	// 		}
	// 	}
	// }
}
