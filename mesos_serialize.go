package main

import (
	"fmt"
	pb "github.com/QubitProducts/jobbo/proto"
	mesos "github.com/QubitProducts/jobbo/proto/mesos/v1"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"regexp"
)

func convertCommand(cmd *pb.Command) *mesos.CommandInfo {
	uris := []*mesos.CommandInfo_URI{}
	for _, uri := range cmd.Uris {
		uris = append(uris, &mesos.CommandInfo_URI{
			Value:      proto.String(uri.Value),
			Executable: proto.Bool(uri.Executable),
			Extract:    proto.Bool(uri.Extract),
			Cache:      proto.Bool(uri.Cache),
			OutputFile: proto.String(uri.OutputFile),
		})
	}

	environment := &mesos.Environment{
		Variables: []*mesos.Environment_Variable{},
	}
	for k, v := range cmd.Environment {
		environment.Variables = append(environment.Variables, &mesos.Environment_Variable{
			Name:  proto.String(k),
			Value: proto.String(v),
		})
	}

	return &mesos.CommandInfo{
		Uris:        uris,
		Environment: environment,
		Shell:       proto.Bool(cmd.Shell),
		Value:       proto.String(cmd.Value),
		Arguments:   cmd.Arguments,
		User:        proto.String(cmd.User),
	}
}

func convertContainer(container *pb.Container) *mesos.ContainerInfo {
	volumes := []*mesos.Volume{}
	for _, vol := range container.Volumes {
		mode := mesos.Volume_RO
		if vol.ReadWrite {
			mode = mesos.Volume_RW
		}

		volumes = append(volumes, &mesos.Volume{
			Mode:          mode.Enum(),
			ContainerPath: proto.String(vol.ContainerPath),
			HostPath:      proto.String(vol.HostPath),
		})
	}

	portMappings := []*mesos.ContainerInfo_DockerInfo_PortMapping{}
	for _, port := range container.PortMappings {
		portMappings = append(portMappings, &mesos.ContainerInfo_DockerInfo_PortMapping{
			HostPort:      proto.Uint32(port.HostPort),
			ContainerPort: proto.Uint32(port.ContainerPort),
			Protocol:      proto.String(port.Protocol),
		})
	}

	return &mesos.ContainerInfo{
		Type:    mesos.ContainerInfo_DOCKER.Enum(),
		Volumes: volumes,
		Docker: &mesos.ContainerInfo_DockerInfo{
			Image:          proto.String(container.Image),
			Network:        mesos.ContainerInfo_DockerInfo_BRIDGE.Enum(),
			PortMappings:   portMappings,
			Privileged:     proto.Bool(false),
			ForcePullImage: proto.Bool(container.ForcePullImage),
		},
	}
}

func convertResources(resources *pb.Resources) []*mesos.Resource {
	return []*mesos.Resource{
		&mesos.Resource{
			Name: proto.String("cpus"),
			Type: mesos.Value_SCALAR.Enum(),
			Scalar: &mesos.Value_Scalar{
				Value: proto.Float64(float64(resources.Cpus)),
			},
		},
		&mesos.Resource{
			Name: proto.String("mem"),
			Type: mesos.Value_SCALAR.Enum(),
			Scalar: &mesos.Value_Scalar{
				Value: proto.Float64(float64(resources.Mem)),
			},
		},
		&mesos.Resource{
			Name: proto.String("disk"),
			Type: mesos.Value_SCALAR.Enum(),
			Scalar: &mesos.Value_Scalar{
				Value: proto.Float64(float64(resources.Disk)),
			},
		},
	}
}

func taskID(job *pb.Job) string {
	return fmt.Sprintf("%v.%v", job.Metadata.Name, job.Metadata.Uid)
}

var taskIDRegex = regexp.MustCompile(`^(.+)\.([^.]+)$`)

func decodeTaskID(taskID string) (string, string, error) {
	mtch := taskIDRegex.FindStringSubmatch(taskID)
	if len(mtch) != 3 || mtch[1] == "" || mtch[2] == "" {
		return "", "", errors.Errorf("taskID did not parse: expected: <job name>.<job uid>, actual: %v", taskID)
	}

	return mtch[1], mtch[2], nil
}

func convertJob(job *pb.Job) *mesos.TaskInfo {
	var container *mesos.ContainerInfo
	if job.Container != nil {
		container = convertContainer(job.Container)
	}

	return &mesos.TaskInfo{
		Name: proto.String(job.Metadata.Name),
		TaskId: &mesos.TaskID{
			Value: proto.String(taskID(job)),
		},
		Resources: convertResources(job.Resources),
		Container: container,
		Command:   convertCommand(job.Command),
	}

}

func convertMesosStatus(state mesos.TaskState) (pb.TaskStatus, error) {
	switch state {
	case mesos.TaskState_TASK_STAGING:
		return pb.TaskStatus_SCHEDULED, nil
	case mesos.TaskState_TASK_STARTING:
		return pb.TaskStatus_SCHEDULED, nil
	case mesos.TaskState_TASK_RUNNING:
		return pb.TaskStatus_RUNNING, nil
	case mesos.TaskState_TASK_KILLING:
		return pb.TaskStatus_RUNNING, nil
	case mesos.TaskState_TASK_FINISHED:
		return pb.TaskStatus_FINISHED, nil
	case mesos.TaskState_TASK_FAILED:
		return pb.TaskStatus_FAILED, nil
	case mesos.TaskState_TASK_KILLED:
		return pb.TaskStatus_FAILED, nil
	case mesos.TaskState_TASK_LOST:
		return pb.TaskStatus_FAILED, nil
	case mesos.TaskState_TASK_ERROR:
		return pb.TaskStatus_FAILED, nil
	default:
		return 0, errors.Errorf("Unrecognised mesos task state %v", state)
	}
}
