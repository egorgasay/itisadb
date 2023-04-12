package servernumber

func Get() {
	var snum int32
	r := db.QueryRow("SELECT server_number FROM info")
	if err = r.Scan(&snum); err != nil {
		return nil, err
	}

	if snum == 0 {
		_, err := db.Exec("UPDATE info SET server_number = ?", *server)
		if err != nil {
			return nil, err
		}
	}
}
