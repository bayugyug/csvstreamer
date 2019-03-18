## csvstreamer

* [x] Simple golang utility library for parsing/writing CSV file using channel stream


### Download the utility library


```sh
		go get -v github.com/bayugyug/csvstreamer
```

### Sanity check


```sh	    
		go test -v ./...
		
		go test ./... -bench=. -test.benchmem -v 2>/dev/null
		 
```


### Read the CSV file (simple)

```go	    
		package main

		import (
				"fmt"
				"os"
				"strings"
				"time"

				"github.com/bayugyug/csvstreamer"
		)

		func main() {
				start := time.Now()

				if len(os.Args) <= 1 {
						fmt.Println("Oops, path parameter is invalid")
						return
				}
				dirPath := strings.TrimSpace(os.Args[1])
				fmt.Println("Path", dirPath, "#", len(os.Args))

				csv := csvstreamer.New(dirPath)
				result, errc := csv.Simple()

				idx := 0
				for _, p := range result {
						idx++
						fmt.Println(idx, "result:", p.Row)
				}
				//sanity check
				if errc != nil {
						fmt.Println("ERROR:", errc)
				}
				fmt.Println(csv.Version())
				fmt.Println("Since", time.Since(start))
				fmt.Println("Done")

		}


		 
```



### Read the CSV file (via channel)

```go	    
		package main

		import (
				"fmt"
				"os"
				"strings"
				"time"

				"github.com/bayugyug/csvstreamer"
		)

		func main() {
				start := time.Now()

				if len(os.Args) <= 1 {
						fmt.Println("Oops, path parameter is invalid")
						return
				}
				dirPath := strings.TrimSpace(os.Args[1])
				fmt.Println("Path", dirPath, "#", len(os.Args))

				done := make(chan struct{})
				defer close(done)

				csv := csvstreamer.New(dirPath)
				result, errc := csv.Parse(done)
				idx := 0
				for p := range result {
						idx++
						fmt.Println(idx, "result:", p.Row)
				}
				//sanity check
				if err := <-errc; err != nil {
						fmt.Println("ERROR:", err)
				}
				fmt.Println(csv.Version())
				fmt.Println("Since", time.Since(start))
				fmt.Println("Done")
		}

		 
```
	

### Write into the CSV file 

```go	    
		package main

		import (
				"fmt"
				"os"
				"time"

				"github.com/bayugyug/csvstreamer"
		)

		func main() {
				start := time.Now()

				//test data
				records := [][]string{
						{"first_name", "last_name", "username"},
				}
				for i := 0; i <= 100; i++ {
						records = append(records, []string{
								fmt.Sprintf("first::04d%d", i+1),
								fmt.Sprintf("last::04d%d", i+1),
								fmt.Sprintf("user::04d%d", i+1),
						})
				}

				csv := csvstreamer.New("save-to-raw.csv")
				for _, rec := range records {
						//Save ( true<append>, rec<[]string>)
						if ok, err := csv.Save(true, rec); !ok || err != nil {
								fmt.Println("ERROR:", err)
								os.Exit(1)
						}
				}
				fmt.Println("Since", time.Since(start))
				fmt.Println("Done")
		}
		 
```
### Notes

	

### Reference


### License

[MIT](https://bayugyug.mit-license.org/)

