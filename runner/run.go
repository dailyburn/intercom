package main

import "bitbucket.org/dailyburn/intercom"

func main() {
	c := intercom.NewIntercomClient("", "", 100)
	c.UpdateUser(nil)
}
