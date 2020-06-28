package mogo

func ExampleConn() {
	db, err := Conn("127.0.0.1:27017/test")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func ExampleQuery() {

}

func ExampleCreate() {

}

func ExampleUpdate() {

}
