package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/unclever-labs/xphilx/xphilx"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var xphilxCfg xphilx.Config
var logLevel string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xphilx",
	Short: "exfiltrate layer 7 payloads to s3",
	Long:  `exfiltrate layer 7 payloads to s3`,
	Run:   xphilxRun,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "config file (default is $HOME/.tmp.yaml)")

	rootCmd.Flags().IntVarP(&xphilxCfg.LogsPerFile, "logs-per-file", "l", 1000, "Logs per file to send to s3")
	rootCmd.Flags().Int32VarP(&xphilxCfg.SnapLength, "snap-length", "s", 1600, "SnapLen for pcap packet capture")
	rootCmd.Flags().StringVarP(&xphilxCfg.Port, "port", "p", "80", "Port to capture packets on")
	rootCmd.Flags().StringVarP(&xphilxCfg.Interface, "interface", "i", "eth0", "Network interface to listen on")
	rootCmd.Flags().StringVarP(&xphilxCfg.S3BucketPath, "s3-bucket-path", "b", "", "Name and path of s3 bucket in format s3://bucket-here/path/here/")

	setFlagsFromEnv(rootCmd)
	setPFlagsFromEnv(rootCmd)
}

func setPFlagsFromEnv(cmd *cobra.Command) {
	// Courtesy of https://github.com/coreos/pkg/blob/master/flagutil/env.go
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		key := strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))
		if val := os.Getenv(key); val != "" {
			if err := cmd.PersistentFlags().Set(f.Name, val); err != nil {
				fmt.Println("Failed setting flag from env:", err)
			}
		}
	})
}

func setFlagsFromEnv(cmd *cobra.Command) {
	// Courtesy of https://github.com/coreos/pkg/blob/master/flagutil/env.go
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		key := strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))
		if val := os.Getenv(key); val != "" {
			if err := cmd.Flags().Set(f.Name, val); err != nil {
				fmt.Println("Failed setting flag from env:", err)
			}
		}
	})
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func xphilxRun(cmd *cobra.Command, args []string) {
	xphilx.Run(xphilxCfg)
}
