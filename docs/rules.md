# General Rules

## Parameters

- When retrieving values from the configuration, it's essential to use the commandToConfigString function to ensure consistency and avoid hardcoding. This approach enables easier maintenance and updates to the configuration.
  #### Correct Usage
  ```go
  var exampleCmd = &cobra.Command{
  	Use:   "example",
  	Short: "foo bar",
  	Long:  `hello world`,
  	Run: func(cmd *cobra.Command, args []string) {
            cn := commandToConfigString(*cmd)
            fmt.Println(viper.GetString(cn+".message"))
    },
  }
  ```
  #### Incorrect Usage (Avoid)
  ```go
  var exampleCmd = &cobra.Command{
  	Use:   "example",
  	Short: "foo bar",
  	Long:  `hello world`,
  	Run: func(cmd *cobra.Command, args []string) {
            fmt.Println(viper.GetString("example.message"))
      },
  }
  ```
