package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"text/template"

	"github.com/mikerybka/util"
)

func main() {
	// Gather inputs
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "<pkg>")
		return
	}
	pkg := os.Args[1]

	// Read server.go file
	fset := token.NewFileSet()
	path := filepath.Join(util.HomeDir(), "src", pkg, "server.go")
	serverFile, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Generate client.go file
	client, err := NewClient(serverFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	path = filepath.Join(util.HomeDir(), "src", pkg, "client.go")
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	template.Must(template.New("client.go").Parse(`package {{ .PkgName }}

type Client struct {
	ServerURL string
}

func (c *Client) send(method string, input, output any) error {
	b, err := json.Marshal(input)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", filepath.Join(c.ServerURL, method), bytes.NewReader(b))
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if !res.StatusCode != 200 {
		b, _ = io.ReadAll(res.Body)
		return fmt.Errorf("%s: %s", res.Status, b)
	}
	return json.NewDecoder(r.Body).Decode(output)
}

{{ range .Methods }}

{{ end }}

`)).Execute(f, client)

	// Run go fmt
}

func defaultName(t string) string {
	switch t {
	case "string":
		return "s"
	case "*net/http.Request":
		return "r"
	case "net/http.ResponseWriter":
		return "w"
	case "error":
		return "err"
	case "int":
		return "i"
	case "bool":
		return "ok"
	default:
		panic(fmt.Sprintf("unknown type %s", t))
	}
}

func NewClient(serverFile *ast.File) (*Client, error) {

	return &Client{}, nil
}

type Client struct {
	PkgName string
	Methods []Method
}

type Method struct {
	Name    string
	Inputs  []Field
	Outputs []Field
}

type Field struct {
	Name string
	Type string
}
