generate-proto:
	cd proto && \
	  protoc \
	  --go_out=plugins=grpc:. \
	  jobbo.proto
	cd proto && \
	  protoc \
	  --go_out=. \
	  mesos/v1/mesos.proto
	cd proto && \
	  protoc \
	  --go_out=Mmesos/v1/mesos.proto=github.com/QubitProducts/jobbo/proto/mesos/v1:. \
	  mesos/v1/scheduler/scheduler.proto
