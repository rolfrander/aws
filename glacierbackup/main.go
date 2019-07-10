package main

import (
	"flag"
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

func backupDirectory(conf config, dir string) error {
	//fmt.Println("zipping " + dir)
	filename, err := zipDir(dir)
	if err != nil {
		return err
	}
	//fmt.Println("filename: " + filename)

	sess, err := getSession()
	if err != nil {
		return err
	}

	description := fmt.Sprintf("backup of %s on %s", dir, time.Now())
	err = upload(sess, filename, conf.GlacierVault, description)
	if err != nil {
		return err
	}

	return nil
}

func doLoginGoogle(conf *config) error {
	var tok *oauth2.Token
	tok, err := loginGoogle(conf.GoogleClientID, conf.GoogleClientSecret)
	if err != nil {
		return err
	}
	conf.GoogleRefreshToken = tok.RefreshToken
	return nil
}

func doLoginAws(conf *config) error {
	return nil
}

func main() {

	pLoginGoogle := flag.Bool("logingoogle", false, "Login to google sheets")
	pLoginAws := flag.Bool("loginaws", false, "Login to AWS")
	pConfigFile := flag.String("cfg", ".glacier.yaml", "path to config file")

	flag.Parse()

	var conf config
	err := conf.readConfig(*pConfigFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if *pLoginGoogle || *pLoginAws {
		if *pLoginGoogle {
			err = doLoginGoogle(&conf)
			if err != nil {
				fmt.Println("Error logging in to google:")
				fmt.Println(err.Error())
			}
		}
		if *pLoginAws {
			err = doLoginAws(&conf)
			if err != nil {
				fmt.Println("Error logging in to aws:")
				fmt.Println(err.Error())
			}
		}
		conf.writeConfig(*pConfigFile)
	} else {
		isloggedin := true
		if conf.GoogleRefreshToken == "" {
			fmt.Println("not logged in to google, use -logingoogle")
			isloggedin = false
		}
		if conf.AwsRefreshToken == "" {
			fmt.Println("not logged in to amazon, use -loginaws")
			isloggedin = false
		}
		if !isloggedin {
			return
		}

		if len(flag.Args()) > 1 {
			dir := flag.Args()[0]
			err = backupDirectory(conf, dir)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
