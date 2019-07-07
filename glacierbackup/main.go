package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	var conf Config
	err := conf.readConfig("/home/rolfn/.glacier.yaml")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(os.Args) > 1 {
		dir := os.Args[1]
		//fmt.Println("zipping " + dir)
		filename, err := zipDir(dir)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		//fmt.Println("filename: " + filename)

		sess, err := getSession()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		description := fmt.Sprintf("backup of %s on %s", dir, time.Now())
		err = upload(sess, filename, conf.Vault, description)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	}

}
