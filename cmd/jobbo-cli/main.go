package main

import (
	"context"
	"fmt"
	pb "github.com/QubitProducts/jobbo/proto"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"os"
	"strings"
)

var (
	serverAddr string
	insecure   bool

	quiet bool

	client pb.JobboClient
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "jobbo-cli",
		Short: "CLI client to the Jobbo server",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			opts := []grpc.DialOption{}
			if insecure {
				opts = append(opts, grpc.WithInsecure())
			}
			conn, err := grpc.Dial(serverAddr, opts...)
			if err != nil {
				return errors.Wrap(err, "Could not connect to jobbo server")
			}

			client = pb.NewJobboClient(conn)
			return nil
		},
	}
	rootCmd.PersistentFlags().StringVarP(&serverAddr, "addr", "a", "", "Location of jobbo server")
	rootCmd.PersistentFlags().BoolVar(&insecure, "insecure", false, "Connect to the jobbo server without encryption")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Minimise output")

	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(rmCmd())
	rootCmd.AddCommand(logsCmd())

	rootCmd.Execute()
}

func listCmd() *cobra.Command {
	var (
		status string
		name   string
		uid    string
	)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List instances managed by the Jobbo server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ts, ok := pb.TaskStatus_value[strings.ToUpper(status)]
			taskStatus := pb.TaskStatus(ts)
			if !ok {
				if status == "" {
					taskStatus = pb.TaskStatus_UNKNOWN
				} else {
					return errors.Errorf("Unknown task status %v", status)
				}
			}

			resp, err := client.ListInstances(context.Background(), &pb.ListInstancesRequest{
				Status: taskStatus,
				Name:   name,
				Uid:    uid,
			})
			if err != nil {
				return errors.Wrap(err, "Could not list instances")
			}

			for _, inst := range resp.Instances {
				if quiet {
					fmt.Printf("%v\n", inst.Job.Metadata.Uid)
				} else {
					switch inst.State.Status {
					case pb.TaskStatus_FAILED:
						fmt.Printf("%v.%v - %v - %v\n", inst.Job.Metadata.Name, inst.Job.Metadata.Uid, inst.State.Status, inst.State.FailureReason)
					default:
						fmt.Printf("%v.%v - %v\n", inst.Job.Metadata.Name, inst.Job.Metadata.Uid, inst.State.Status)
					}
				}
			}
			return nil
		},
	}

	listCmd.Flags().StringVarP(&status, "status", "s", "", "Task status to filter on")
	listCmd.Flags().StringVarP(&name, "name", "n", "", "Task name to filter on")
	listCmd.Flags().StringVarP(&uid, "uid", "u", "", "Task uid to filter on")

	return listCmd
}

func rmCmd() *cobra.Command {
	var (
		status string
		uid    string
	)

	rmCmd := &cobra.Command{
		Use:   "rm",
		Short: "Removes an instance managed by the Jobbo server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ts, ok := pb.TaskStatus_value[strings.ToUpper(status)]
			taskStatus := pb.TaskStatus(ts)
			if !ok {
				if status == "" {
					taskStatus = pb.TaskStatus_UNKNOWN
				} else {
					return errors.Errorf("Unknown task status %v", status)
				}
			}

			resp, err := client.RemoveInstance(context.Background(), &pb.RemoveInstanceRequest{
				Status: taskStatus,
				Uid:    uid,
			})
			if err != nil {
				return errors.Wrap(err, "Could not remove instance")
			}

			if quiet {
				fmt.Printf("%v\n", resp.Instance.Job.Metadata.Uid)
			} else {
				fmt.Printf("Removed %v.%v\n", resp.Instance.Job.Metadata.Name, resp.Instance.Job.Metadata.Uid)
			}

			return nil
		},
	}

	rmCmd.Flags().StringVarP(&status, "status", "s", "", "Task status to filter on")
	rmCmd.Flags().StringVarP(&uid, "uid", "u", "", "Task uid to filter on")

	return rmCmd
}

func logsCmd() *cobra.Command {
	var (
		uid string
	)

	logsCmd := &cobra.Command{
		Use:   "logs",
		Short: "Returns the logs for a task",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := client.GetInstanceLogs(context.Background(), &pb.GetInstanceLogsRequest{
				Uid: uid,
			})
			if err != nil {
				return errors.Wrap(err, "Could not get instance logs")
			}

			os.Stderr.Write(resp.Stderr)
			os.Stdout.Write(resp.Stdout)

			return nil
		},
	}

	logsCmd.Flags().StringVarP(&uid, "uid", "u", "", "Task uid to filter on")

	return logsCmd
}
