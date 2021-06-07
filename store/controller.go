func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	
	// get all products
	products := c.Repository.GetProducts()
	data, _ := json.Marshal(products)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

func (c *Controller) AddProduct(w http.ResponseWriter, r *http.Request) {
	var product Product

	// https://golang.org/pkg/io/#LimitedReader
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		log.Fatalln("Error AddProduct", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error AddProduct", err)

	}

	if err := json.Unmarshal(body, &product); err != nil {
		w.WriteHeader(422) // unprocessable entity

		if json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error AddProduct unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	success := c.Repository.AddProduct(product)
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)
    return
}