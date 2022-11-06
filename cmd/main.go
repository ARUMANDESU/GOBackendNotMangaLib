package main

import (
	"flag"
	"log"
	"notmangalib.com/pkg"
	"notmangalib.com/pkg/myhttp"
	"os"
)

func main() {
	addr := flag.String("addr", ":5000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	server := &myhttp.Server{}
	server.StartServer(addr, errorLog, infoLog, pkg.Routes())
}
