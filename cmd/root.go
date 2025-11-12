package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/chege/azfind/internal/completion"
	"github.com/chege/azfind/internal/fzfui"
	"github.com/chege/azfind/internal/syncer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var listCache bool
var doSync bool
var doCompletion bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "azf [query]",
	Short: "Fast Azure Resource Finder",
	Long:  `"azf" is a fast CLI for searching, filtering, and opening Azure resources in the browser or running supported actions.`,
	Args:  cobra.ArbitraryArgs, // allow arbitrary args for search input
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		if doCompletion {
			shell := "bash"
			if len(args) > 0 {
				shell = args[0]
			}
			return generateCompletionScript(cmd.Root(), shell, os.Stdout)
		}

		if doSync {
			return syncer.SyncAll(ctx)
		}

		if listCache {
			return nil // list-cache will be reimplemented later
		}

		return fzfui.RunSearch(ctx, args)
	},
}

// generateCompletionScript prints shell completion to w based on the given shell.
func generateCompletionScript(cmd *cobra.Command, shell string, w io.Writer) error {
	switch shell {
	case "bash":
		return cmd.GenBashCompletion(w)
	case "zsh":
		return cmd.GenZshCompletion(w)
	case "fish":
		return cmd.GenFishCompletion(w, true)
	case "powershell", "pwsh":
		return cmd.GenPowerShellCompletionWithDesc(w)
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.azfind.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolVar(&listCache, "list-cache", false, "List cached Azure resources")
	rootCmd.Flags().BoolVar(&doSync, "sync", false, "Synchronize Azure resources into local cache")
	rootCmd.Flags().BoolVar(&doCompletion, "completion", false, "Generate dynamic name completions")

	_ = rootCmd.Flags().MarkHidden("completion")

	rootCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ctx := context.Background()
		list, _ := completion.Generate(ctx, toComplete)
		return list, cobra.ShellCompDirectiveNoFileComp
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".azfind" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".azfind")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
