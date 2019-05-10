package cmd

import (
	"fmt"
	"net"
	"os"
	"os/exec"

	"github.com/cybozu-go/topolvm/lvmd"
	"github.com/cybozu-go/topolvm/lvmd/proto"
	"github.com/cybozu-go/well"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var config struct {
	vgName     string
	socketName string
}

const defaultSocketName = "/run/topolvm/lvmd.sock"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lvmd",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return subMain()
	},
}

func subMain() error {
	err := well.LogConfig{}.Apply()
	if err != nil {
		return err
	}

	err = checkVG(config.vgName)
	if err != nil {
		return err
	}

	lis, err := net.Listen("unix", config.socketName)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	proto.RegisterLVServiceServer(grpcServer, lvmd.NewLVService(config.vgName))
	proto.RegisterVGServiceServer(grpcServer, lvmd.NewVGService(config.vgName))
	return grpcServer.Serve(lis)
}

func checkVG(vg string) error {
	return exec.Command("vgs", vg).Run()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&config.vgName, "volumegroup", "", "LVM volume group name")
	rootCmd.Flags().StringVar(&config.socketName, "listen", defaultSocketName, "Unix domain socket name")
}
