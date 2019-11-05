# my_tracker

Malaysia's Shipment Tracking API

### Supported vendor
- PosLaju
- J & T Express

### USAGE

```
package main

import (
	"fmt"
	"github.com/arma7x/my_tracker"
)

func main() {
	fmt.Println(my_tracker.PosLaju("TRACK_CODE_HERE"))
	fmt.Println(my_tracker.JnT("TRACK_CODE_HERE"))
}
```
