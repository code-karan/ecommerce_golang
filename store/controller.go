package store

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "fmt"
    "net/http"
    "strings"

    "github.com/gorilla/context"
    "github.com/dgrijalva/jwt-go"

)

type Controller struct {
    Repository Repository
}

func (c *Controller) GetToken(w http.ResponseWriter, req *http.Request) {
    var user User
    _ = json.NewDecoder(req.Body).Decode(&user)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": user.Username,
        "password": user.Password,
    })

    if user.Username == "" || user.Password == "" {
        fmt.Fprintf(w, "Please provide a username and a password")
        return
    }

    log.Println("Username: " + user.Username);
    log.Println("Password: " + user.Password);

    tokenString, error := token.SignedString([]byte("secret"))
    if error != nil {
        fmt.Println(error)
    }
    json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
}

func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        authorizationHeader := req.Header.Get("authorization")
        if authorizationHeader != "" {
            bearerToken := strings.Split(authorizationHeader, " ")
            if len(bearerToken) == 2 {
                token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
                    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, fmt.Errorf("There was an error")
                    }
                    return []byte("secret"), nil
                })
                if error != nil {
                    json.NewEncoder(w).Encode(Exception{Message: error.Error()})
                    return
                }
                if token.Valid {
                    log.Println("TOKEN WAS VALID")
                    context.Set(req, "decoded", token.Claims)
                    next(w, req)
                } else {
                    json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
                }
            }
        } else {
            json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
        }
    })
}


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

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln("Error AddProduct", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error AddProduct", err)

	}

	if err := json.Unmarshal(body, &product); err != nil {
		w.WriteHeader(422) // error unprocessable entity

		if json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error AddProduct unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
    fmt.Println(product)

	success := c.Repository.AddProduct(product)
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)
    fmt.Fprintf(w, "Product added successfully")
    return
}

func (c *Controller) UpdateProduct(w http.ResponseWriter, r *http.Request) {
    var product Product
	body, err := ioutil.ReadAll(r.Body)

    if err != nil {
        fmt.Fprintf(w, "Couldn't read request body")
    }

    if err := json.Unmarshal(body, &product); err != nil {
		w.WriteHeader(422) // error unprocessable entity

		if json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error AddProduct unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

    success := c.Repository.UpdateProduct(product)
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}