type partialResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

// create a buffer that fufills the response writer interface, for the purposes of rending partial responses to ajax requests
func (pw *partialResponseWriter) Write(p []byte) (int, error) {
	return pw.buf.Write(p)
}

// this middleware detects if it is an alpine-ajax request, then filters for requested sections.  If not, it serves the whole  page.
func alpineAjax(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Create a response wrapper:
		isAjax, err := strconv.ParseBool(r.Header.Get("X-Alpine-Request"))
		if isAjax && err == nil {
			fmt.Println("DEBUG: Alpine Ajax request")
			myResponseWriter := &partialResponseWriter{
				ResponseWriter: w,
				buf:            &bytes.Buffer{},
			}
			next.ServeHTTP(myResponseWriter, r)
			//only return the css selectors requested by the ajax request, and the x-sync too
			doc, err := html.Parse(myResponseWriter.buf)
			if err != nil {
				fmt.Println("Error:", err)
				next.ServeHTTP(w, r)
				return
			}
			targets := strings.Split(r.Header.Get("X-Alpine-Target"), " ")
			for _, targetid := range targets {
				var targetNode *html.Node
				var crawler func(*html.Node)
				crawler = func(node *html.Node) {
					if node.Type == html.ElementNode {
						for _, a := range node.Attr {
							if a.Key == "x-sync" {
								targetNode = node
								return
							}
							if a.Val == targetid {
								targetNode = node
								return
							}
						}
					}
					for child := node.FirstChild; child != nil; child = child.NextSibling {
						crawler(child)
					}
				}
				crawler(doc)
				if targetNode != nil {
					//html.Render(os.Stdout, targetNode)
					html.Render(myResponseWriter, targetNode)
				}
			}
			if _, err := io.Copy(w, myResponseWriter.buf); err != nil {
				log.Printf("Alpine Ajax middleware failed: %v", err)
			}
		} else {
			fmt.Println("DEBUG: Non-ajax request")
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}
