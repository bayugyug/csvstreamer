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
				csvFilename := strings.TrimSpace(os.Args[1])
				csv := csvstreamer.New(csvFilename)
				result, errc := csv.Simple()

				for idx, p := range result {
						fmt.Println(idx+1, "result:", p.Row)
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
				csvFilename := strings.TrimSpace(os.Args[1])
				done := make(chan struct{})
				defer close(done)

				csv := csvstreamer.New(csvFilename)
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
				for i := 0; i < 100; i++ {
						records = append(records, []string{
								fmt.Sprintf("first::%04d", i+1),
								fmt.Sprintf("last::%04d", i+1),
								fmt.Sprintf("user::%04d", i+1),
						})
				}
				csv := csvstreamer.New("save-to-raw.csv")
				if err := csv.Append(records...); err != nil {
					fmt.Println("ERROR:", err)
					os.Exit(1)
				}
				fmt.Println("Since", time.Since(start))
				fmt.Println("Done")
		}
		 
```



### Write into the CSV file (by batch)

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
				for i := 0; i < 100; i++ {
						records = append(records, []string{
								fmt.Sprintf("bybatch::first::%04d", i+1),
								fmt.Sprintf("bybatch::last::%04d", i+1),
								fmt.Sprintf("bybatch::user::%04d", i+1),
						})
				}
				csv := csvstreamer.New("save-to-raw.csv")
				if err := csv.AppendBatch(records, 20); err != nil {
					fmt.Println("ERROR:", err)
					os.Exit(1)
				}
				fmt.Println("Since", time.Since(start))
				fmt.Println("Done")
		}
		 
```


### Notes

	

### Reference


### License

[MIT](https://bayugyug.mit-license.org/)

