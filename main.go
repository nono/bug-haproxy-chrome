package main

import (
	"fmt"
	"net/http"
	"time"
)

const port = 7000
const html = `
<!DOCTYPE html>
<html>
<title>Bug haproxy - Chrome</title>
<h1>Bug haproxy - Chrome</h1>
<form action="/upload" method="POST">
<input type="file" id="file-input" />
</form>
<script>
var input = document.getElementById("file-input");
input.addEventListener('change', function() {
  console.log('change');
  var file = input.files[0];
  var fd = new FormData();
  fd.append('upload', file);
  fetch('/upload', { method: 'POST', body: fd }).catch(console.error.bind(console));
});
</script>
</html>
`

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, html)
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(409)
		w.Write([]byte("Conflict!"))
	})

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
