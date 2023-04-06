# Yet Another Crappy Scanner (YACS)
YACS is my poor attempt on using minequery with masscan

If you ever want to use this thing then build it

## How do I build YACS?
```
go mod tidy
go build .
```
Easy, right?

## A little guide for you to understand how to use it.
```
./golang-scanner --range (ip range) --port-range (I bet you know how to use this one) --rate (masscan rate thingy) --output (file to save all of the servers)
```
Not hard at all.
