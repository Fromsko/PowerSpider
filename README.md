# PowerSpider

This is a project for school private net where help our to get ele-money data.

# Examples

```go
package main

import (
    "log"
    "PowerSpider/config"
    "PowerSpider/core"
)

func main() {
    // Init configure file
    config.InitConfig(
        config.Config{
            User  : "202127530334",
            Pwd   : "102018",
            Timer : config.Timer{
                TimeUint: "hourse",
                TimeInfo: 2,
            },
            ResDir: "res",
            Proxy : "http://localhost:7980"
            RedictUrl : "http://10.13.14.20:9999/"
        }
    )
    // Start appliction
    if err := core.Start(); err != nil {
        log.Fprint("[Error] %s", err)
    }
}
```
