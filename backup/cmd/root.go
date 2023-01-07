/*
Copyright © 2023 Teppei Sudo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type BackupConfig struct {
	Resources []Resource
}

type Resource struct {
	Namespace string
	Kinds     []string
}

var cfgFile string
var outputdir string
var backupConfig BackupConfig

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kube-tracker",
	Short: "Simple application to take backup k8s resources",
	Long:  `Simple application to take backup k8s resources`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(backupConfig)
		mkdir(outputdir)
		for _, r := range backupConfig.Resources {
			var command *exec.Cmd
			var commandStdErr bytes.Buffer
			var commandStdOut bytes.Buffer
			var fileName string
			if r.Namespace == "" {
				fileName = "cluster"
				mkdir(path.Join(outputdir, fileName))
				command = exec.Command("kubectl", "get", strings.Join(r.Kinds, ","))
			} else {
				fileName = r.Namespace
				mkdir(path.Join(outputdir, fileName))
				command = exec.Command("kubectl", "-n", r.Namespace, "get", strings.Join(r.Kinds, ","))
			}
			command.Stderr = &commandStdErr
			command.Stdout = &commandStdOut
			command.Run()
			writeToFile(commandStdOut.String(), path.Join(outputdir, fileName, "summary"))
			// fmt.Println(commandStdOut.String())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kube-tracker.yaml)")
	rootCmd.PersistentFlags().StringVar(&outputdir, "dir", "", "output directory path")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		// Search config in home directory with name ".hoge" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".kube-tracker.yaml")
	}

	viper.SetConfigType("yml")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("couldn't read config file")
	}
	fmt.Println("Using config file:", viper.ConfigFileUsed())

	if err := viper.Unmarshal(&backupConfig); err != nil {
		log.Fatal(err)
	}
}

func GetRawConfig() (config api.Config) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.RawConfig()
	if err != nil {
		log.Fatal("couldn't get kubeconfig")
	}
	return config
}

func mkdir(dir string) {
	if _, err := os.Stat(dir); err != nil {
		if err := os.Mkdir(dir, 0666); err != nil {
			log.Println("cannot make new dir: ", err)
		}
	}
}

func writeToFile(input string, dstPath string) {
	dst, err := os.Create(dstPath)
	if err != nil {
		log.Println("cannot create new file: ", err)
	}
	defer dst.Close()
	data := []byte(input)
	_, err = dst.Write(data)
	if err != nil {
		log.Println("cannot write data to file: ", err)
	}
}
