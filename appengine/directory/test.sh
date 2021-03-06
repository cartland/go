
echo "GET"
curl -H "Content-Type: application/json" http://localhost:8888 -f -X GET
echo

echo "PUT"
curl -H "Content-Type: application/json" -d '{"name":"directory","location":"http://localhost:8888"}' http://localhost:8888 -f -X PUT 
echo

echo "PUT"
curl -H "Content-Type: application/json" -d '{"name":"service","location":"https://loadbalance-golang.appspot.com/reliable"}' http://localhost:8888 -f -X PUT 
echo

echo "PUT"
curl -H "Content-Type: application/json" -d '{"name":"service","location":"https://loadbalance-golang.appspot.com/flaky"}' http://localhost:8888 -f -X PUT 
echo

echo "PUT"
curl -H "Content-Type: application/json" -d '{"name":"garbage","location":"garbage"}' http://localhost:8888 -f -X PUT 
echo

echo
sleep 1


echo "GET"
curl -H "Content-Type: application/json" http://localhost:8888 -f -X GET
echo

echo "HEARTBEAT"
curl -H "Content-Type: application/json" http://localhost:8888/heartbeat -f -X GET
echo

echo
sleep 1


echo "GET"
curl -H "Content-Type: application/json" http://localhost:8888 -f -X GET
echo
