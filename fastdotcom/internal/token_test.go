package internal

import "testing"

func TestExtractJSPath(t *testing.T) {
	fakeHTML := `
<!DOCTYPE html>
<html>
    <head>
        <title>fastester.com</title>
    </head>
    <body>
        <script src="/app-888deadbeef888.js"></script>
    </body>
</html>
`
	if jsPath := extractJSPath(fakeHTML); jsPath != "/app-888deadbeef888.js" {
		t.Errorf("got: %v", jsPath)
	}
}

func TestExtractToken(t *testing.T) {
	fakeJS := `
;(function(){
const params={asdf:'FFFFF',token:'FooBarBaz'};
})();
`
	if token := extractToken(fakeJS); token != "FooBarBaz" {
		t.Errorf("got: %v", token)
	}
}
