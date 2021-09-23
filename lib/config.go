package goshleep

import (
	"log"
	"os/user"

	"github.com/spf13/viper"
)

// AutoConfig automatically builds the configuration file.
func AutoConfig() {
	usr, _ := user.Current()

	viper.SetDefault("storage", usr.HomeDir+"/.config/goshleep/")
	viper.SetDefault("prefix", "+")
	viper.SetDefault("cmdPrefix", "+")
	viper.SetDefault("eightballMessages", []string{"Senpai, pls no ;-;",
		"Take a wild guess...",
        "Without a doubt", "No", "Yes", "You'll be the judge", "Sure",
        "Of course", "No way", "No... (╯°□°）╯︵ ┻━┻", "Very doubtful",
        "Most likely", "Might be possible" })

	viper.SetConfigName("conf")                    // name of config file (without extension)
	viper.SetConfigType("yaml")                    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$HOME/.config/goshleep/") // call multiple times to add many search paths
	viper.AddConfigPath("/etc/goshleep/")          // path to look for the config file in
	viper.AddConfigPath(".")                       // optionally look for config in the working directory
	viper.SafeWriteConfig()
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.WriteConfig()
			// Config file not found; ignore error if desired
		} else {
			log.Println(err)
			// Config file was found but another error was produced
		}
	}
	// Config file found and successfully parsed

}
