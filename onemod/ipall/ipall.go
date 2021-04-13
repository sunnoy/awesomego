/*
 *@Description
 *@author          lirui
 *@create          2021-04-13 18:07
 */
package main

import (
	"fmt"
	"math/big"
	"net"
)

func main() {
	ss := new(big.Int).SetInt64(88)

	ip := addIPOffset(ss, 99)

	fmt.Printf("ip is %v", ip)
}

func addIPOffset(base *big.Int, offset int) net.IP {
	return net.IP(big.NewInt(0).Add(base, big.NewInt(int64(offset))).Bytes())
}
