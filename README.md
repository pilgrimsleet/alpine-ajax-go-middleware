# alpine-ajax-go-middleware
Go standard library middleware to automatically trim alpine-ajax responses to only inclued targeted ids

This middleware is written to be paired with alpine-ajax, when used with the intended pattern of progressive enhancement.  It examines the headers to see if it is an alpine-ajax request, then trims down the generated page to only inclued the targeted ids as well as elements with "x-sync."  When it is not an alpine-ajax request, this middleware has no effect.
