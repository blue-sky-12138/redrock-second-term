package main

import (
	"SecondTerm/Homework-7/workloadManager/control"
	"SecondTerm/Homework-7/workloadManager/loadBalance"
)

func main() {
	loadBalance.BalanceInit()
	control.RoutersEntrance()
}
