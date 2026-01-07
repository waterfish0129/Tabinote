package main

import "github.com/waterfish0129/Tabinote/cmd"

// @title Tabinote API
// @version 0.0.1
// @description Made By Walter.Fish

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description  請在輸入框輸入 "Bearer {你的Token}"
func main() {
	defer cmd.Clean()
	cmd.Start()

}
