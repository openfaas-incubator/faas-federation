package routing

import (
	"net/http"
	"net/url"
	"sort"
)

// Result to hold the result from each request including an Index
// which will be used for sorting the results after they come in
type Result struct {
	Index    int
	Response http.Response
	Err      error
}

// Get sends requests in parallel but only up to a certain
// limit, and furthermore it's only parallel up to the amount of CPUs but
// is always concurrent up to the concurrency limit
func Get(urls []*url.URL, concurrencyLimit int) []Result {

	// this buffered channel will block at the concurrency limit
	semaphoreChan := make(chan struct{}, concurrencyLimit)

	// this channel will not block and collect the http request results
	resultsChan := make(chan *Result)

	// make sure we close these channels when we're done with them
	defer func() {
		close(semaphoreChan)
		close(resultsChan)
	}()

	// keen an Index and loop through every url we will send a request to
	for i, u := range urls {

		// start a go routine with the Index and url in a closure
		go func(i int, u *url.URL) {

			// this sends an empty struct into the semaphoreChan which
			// is basically saying add one to the limit, but when the
			// limit has been reached block until there is room
			semaphoreChan <- struct{}{}

			// send the request and put the response in a result struct
			// along with the Index so we can sort them later along with
			// any error that might have occoured
			res, err := http.Get(u.String())
			result := &Result{i, *res, err}

			// now we can send the result struct through the resultsChan
			resultsChan <- result

			// once we're done it's we read from the semaphoreChan which
			// has the effect of removing one from the limit and allowing
			// another goroutine to start
			<-semaphoreChan

		}(i, u)
	}

	// make a slice to hold the results we're expecting
	var results []Result

	// start listening for any results over the resultsChan
	// once we get a result append it to the result slice
	for {
		result := <-resultsChan
		results = append(results, *result)

		// if we've reached the expected amount of urls then stop
		if len(results) == len(urls) {
			break
		}
	}

	// let's sort these results real quick
	sort.Slice(results, func(i, j int) bool {
		return results[i].Index < results[j].Index
	})

	// now we're done we return the results
	return results
}
