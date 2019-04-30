run:
	go run .
reset:
	echo "doesnot work :("
	read apikey
	echo $apikey
	cat << EOF > ./aa.json
	{"Apikey": $apikey}
	EOF