package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	source = "/home/tony/gocode/src/github.com/ajtaylor/corvomq/"
	target = "/srv/corvomq"
)

func main() {
	website := flag.Bool("website", false, "Deploys full website")
	websiteHTML := flag.Bool("websiteHTML", false, "Deploys website HTML only")
	websiteServer := flag.Bool("websiteServer", false, "Deploys website server only")
	webapp := flag.Bool("webapp", false, "Deploys webapp")
	css := flag.Bool("css", false, "Deploys CSS")
	api := flag.Bool("api", false, "Deploys API")
	docs := flag.Bool("docs", false, "Deploys website documentation")

	flag.Parse()

	env := os.Environ()

	// var wg sync.WaitGroup

	if *css {
		// deployCSS(&wg, env)
		deployCSS(env)
	}

	if *website {
		// deployWebsite(&wg, env)
		deployWebsite(env)
	}

	if *websiteHTML {
		// deployWebsite(&wg, env)
		deployWebsiteHTML(env)
	}

	if *websiteServer {
		// deployWebsite(&wg, env)
		deployWebsiteServer(env)
	}

	if *webapp {
		deployWebapp(env)
	}

	if *api {
		deployAPI(env)
	}

	if *docs {
		deployDocumentation(env)
	}

	// wg.Wait()
}

func printError(source string, err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v ==> Error: %s\n", source, err.Error()))
	}
}

func printOutput(source string, outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("%v ==> Output:\n\t%v\n", source, string(outs))
	}
}

func runCmd(source string, cmd *exec.Cmd) {
	output, err := cmd.CombinedOutput()
	printError(source, err)
	printOutput(source, output)
}

// func deployCSS(wg *sync.WaitGroup, env []string) {
func deployCSS(env []string) {
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	s := "Deploy CSS"
	// log.Println(s)
	log.Println(1)
	// os.Chdir(source + "www/static/css")
	// cmd := exec.Command("/bin/sh", "-c", "postcss ./src/corvomq.css --no-map -u postcss-import -u postcss-inherit -o ./dist/corvomq.css")
	os.Chdir(source + "css")
	cmd := exec.Command("/bin/sh", "-c", "./wt compile ./corvomq.scss -b ../build/css/")
	runCmd(s, cmd)
	os.MkdirAll(target+"/www/static/css", os.ModePerm)
	// cmd = exec.Command("/bin/sh", "-c", "mkdir -p "+target+"/www/static/css")
	// runCmd(s, cmd)
	cmd = exec.Command("/bin/sh", "-c", "mv ../build/css/corvomq.css "+target+"/www/static/css/corvomq.css")
	runCmd(s, cmd)
	// cmd = exec.Command("/bin/sh", "-c", "cp ../build/css/corvomq.css "+target+"/webapp/static/css/corvomq.css")
	// runCmd(s, cmd)
	cmd = exec.Command("/bin/sh", "-c", "sudo systemctl restart nginx")
	runCmd(s, cmd)
	// }()
}

// func deployWebsite(wg *sync.WaitGroup, env []string) {
func deployWebsite(env []string) {
	deployWebsiteHTML(env)
	deployWebsiteServer(env)
}

// func deployWebsite(wg *sync.WaitGroup, env []string) {
func deployWebsiteHTML(env []string) {
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	s := "Deploy website HTML"
	// log.Println(s)
	// Build HTML
	os.Chdir(source + "www/pug")
	cmd := exec.Command("/bin/sh", "-c", "pug ./*.pug -o ../html")
	runCmd(s, cmd)
	// Make website HTML directory
	os.MkdirAll(target+"/www/html", os.ModePerm)
	// Clear website HTML directory
	cmd = exec.Command("/bin/sh", "-c", "rm -r "+target+"/www/html/*")
	runCmd(s, cmd)
	// Move HTML to website
	cmd = exec.Command("/bin/sh", "-c", "mv ../html/*.html "+target+"/www/html/")
	runCmd(s, cmd)
	cmd = exec.Command("/bin/sh", "-c", "sudo systemctl restart nginx")
	runCmd(s, cmd)
	// }()
}

func deployDocumentation(env []string) {
	s := "Deploy documentation"
	// Build HTML
	os.Chdir(source + "www/pug")
	cmd := exec.Command("/bin/sh", "-c", "pug -O ./options.js ./documentation -o ../html/documentation")
	runCmd(s, cmd)
	// Make documentation directory
	os.MkdirAll(target+"/www/html/documentation", os.ModePerm)
	// Clear documentation directory
	cmd = exec.Command("/bin/sh", "-c", "rm -r "+target+"/www/html/documentation/*")
	runCmd(s, cmd)
	// Move HTML to website
	cmd = exec.Command("/bin/sh", "-c", "mv ../html/documentation/ "+target+"/www/html/")
	runCmd(s, cmd)
	cmd = exec.Command("/bin/sh", "-c", "sudo systemctl restart nginx")
	runCmd(s, cmd)
	// }()
}

// func deployWebsite(wg *sync.WaitGroup, env []string) {
func deployWebsiteServer(env []string) {
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	s := "Deploy website server"
	// Build web server
	os.Chdir(source + "www")
	cmd := exec.Command("/bin/sh", "-c", "go build -i -o www github.com/ajtaylor/corvomq/www")
	runCmd(s, cmd)
	// Move webserver to app
	// gopath := os.Getenv("GOPATH")
	// cmd = exec.Command("/bin/sh", "-c", "cp "+gopath+"/bin/www "+target+"/www/")
	cmd = exec.Command("/bin/sh", "-c", "mv ./www "+target+"/www/")
	runCmd(s, cmd)
	// Copy daemon unit
	cmd = exec.Command("/bin/sh", "-c", "sudo cp ./corvomq-www.service /etc/systemd/system/corvomq-www.service")
	runCmd(s, cmd)
	// Restart running www
	cmd = exec.Command("/bin/sh", "-c", "sudo systemctl daemon-reload")
	runCmd(s, cmd)
	cmd = exec.Command("/bin/sh", "-c", "sudo systemctl restart corvomq-www")
	runCmd(s, cmd)
	// }()
}

func deployWebapp(env []string) {
	var cmd *exec.Cmd
	s := "Deploy webapp"
	// Build http server
	os.Chdir(source + "webapp")
	// cmd = exec.Command("/bin/sh", "-c", "go build -i -o webapp github.com/ajtaylor/corvomq/webapp")
	// runCmd(s, cmd)
	// Move http to app
	// cmd = exec.Command("/bin/sh", "-c", "mv ./webapp "+target+"/webapp/")
	// runCmd(s, cmd)
	// Compile ES6 to js
	cmd = exec.Command("/bin/sh", "-c", "webpack --config ./webpack.config.js")
	runCmd(s, cmd)
	// Make app js directory
	os.MkdirAll(target+"/www/static/javascript", os.ModePerm)
	// cmd = exec.Command("/bin/sh", "-c", "mkdir -p "+target+"/www/static/javascript")
	// runCmd(s, cmd)
	// Move js to app
	cmd = exec.Command("/bin/sh", "-c", "mv ../build/webapp/main.bundle.js "+target+"/www/static/javascript/app.js")
	runCmd(s, cmd)
	// Delete js
	// cmd = exec.Command("/bin/sh", "-c", "rm ./*.js")
	// runCmd(s, cmd)
	// java -jar ~/closure-compiler/closure-compiler-v20170218.jar --js ./app.js --js_output_file app.min.js
	// CSS
	// cmd = exec.Command("/bin/sh", "-c", "cp ./static/css/corvomq.css "+target+"/webapp/static/css/corvomq.css")
	// runCmd(s, cmd)
	// Build HTML
	os.Chdir(source + "webapp/pug")
	cmd = exec.Command("/bin/sh", "-c", "pug ./app.pug -o ../html")
	runCmd(s, cmd)
	// Make app HTML directory
	os.MkdirAll(target+"/www/html", os.ModePerm)
	// cmd = exec.Command("/bin/sh", "-c", "mkdir -p "+target+"/www/html")
	// runCmd(s, cmd)
	// Move HTML to app
	cmd = exec.Command("/bin/sh", "-c", "mv ../html/app.html "+target+"/www/html/")
	runCmd(s, cmd)
	// Restart running www
	cmd = exec.Command("/bin/sh", "-c", "sudo systemctl restart corvomq-www")
	runCmd(s, cmd)
}

func deployAPI(env []string) {
	s := "Deploy API"
	// Build API server
	os.Chdir(source + "api")
	cmd := exec.Command("/bin/sh", "-c", "go build -i -o api github.com/ajtaylor/corvomq/api")
	runCmd(s, cmd)
	// Stop API server
	cmd = exec.Command("/bin/sh", "-c", "sudo systemctl stop corvomq-api")
	runCmd(s, cmd)
	// Copy API server to app
	cmd = exec.Command("/bin/sh", "-c", "mv ./api "+target+"/api/")
	runCmd(s, cmd)
	// Start API server
	cmd = exec.Command("/bin/sh", "-c", "sudo systemctl start corvomq-api")
	runCmd(s, cmd)
}
