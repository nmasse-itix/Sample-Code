package main

import (
  "code.google.com/p/gorest"
  "net/http"
  "fmt"
  "strings"
  "io/ioutil"
  "encoding/json"
)

type Category struct {
  Id string
  Name string
}

var Categories []Category = []Category {
  { "fringues", "Habillage" },
  { "cuisine", "Cuisine" },
  { "digital", "Digital" },
  { "maison", "Bricolage" }  }

type Product struct {
  Id int
  Name string
  Category string
  Image string
  Description string
  Price float32
  Stock int
  VendorId string
  VendorName string
  VendorProductId string
  IsDigital bool
}

type BuyResponse struct {
  ResponseCode string
  DownloadUrl string
}

type CallbackResponse struct {
  ResponseCode string `json:"code"`
  RedirectUrl  string `json:"redirect_url"`
}

type BuyCallback struct {
  VendorId string
  VendorName string
  VendorProductId string
}

var Products []Product = []Product {
  { 0, "T-Shirt", "fringues", "brice.jpg", "Le T-Shirt de Brice de Nice.", 99.9, 1, "", "", "", false },
  { 1, "Pull col roulé", "fringues", "pull.jpg", "Un pull à col roulé de couleur marron.", 2.00, 10, "", "", "", false },
  { 2, "Cocotte minute", "cuisine", "cocotte.jpg", "La cocotte minute 'Presto'.", 45.00, 2, "", "", "", false },
  { 3, "Marteau-Piqueur", "maison", "mp.jpg", "Le marteau piqueur 'DESTRUCTOR 2000'.", 600, 5, "", "", "", false } }
//  { 4, "Visseuse Ultrasonique", "maison", "visseuse.jpg", "La visseuse-dévisseuse de chez Méga Store.", 600, 5, "mega-store", "Méga Store", "0001", false },
//  { 5, "DVD de Harry Poter", "digital", "dvd.jpg", "L'histoire de Ari l'empotteur au pays des merveilles.", 29.9, -1, "zouba-books", "Zouba Books", "12345", true }  }

func main() {
  gorest.RegisterService(new(MyShopService)) // Register our service
  http.Handle("/api/",gorest.Handle())
  http.Handle("/", http.FileServer(http.Dir("www-root")))
  http.ListenAndServe(":8787", nil)
}

// REST Service Definition
type MyShopService struct {
    gorest.RestService `root:"/api/" consumes:"application/json" produces:"application/json"`
    getCategories  gorest.EndPoint `method:"GET" path:"/shop/category/" output:"[]Category"`
    getProductsByCategory  gorest.EndPoint `method:"GET" path:"/shop/product/?{category:string}" output:"[]Product"`
    searchProducts  gorest.EndPoint `method:"GET" path:"/shop/search/{criteria:string}" output:"[]Product"`
    getProduct  gorest.EndPoint `method:"GET" path:"/shop/product/{id:int}" output:"Product"`
    addProduct  gorest.EndPoint `method:"POST" path:"/market/product/" postdata:"Product"`
    buyProduct  gorest.EndPoint `method:"GET" path:"/shop/product/{id:int}/buy" output:"BuyResponse"`
}

func(serv MyShopService) SearchProducts(criteria string) []Product {
  fmt.Println(">>> SearchProducts: criteria = ", criteria)
  if (criteria == "") {
    return Products
  }

  sliceofcriteria := strings.Split(criteria, " ")

  FilteredProducts := []Product {}
  for _, p := range Products {
    var selected bool = false
    name := strings.ToLower(p.Name)
    desc := strings.ToLower(p.Description)

    for _, c := range sliceofcriteria {
      c = strings.ToLower(c)
      if strings.Contains(name, c) || strings.Contains(desc, c) {
        selected = true
      }
    }

    if selected {
      FilteredProducts = append(FilteredProducts, p)
    }
  }

  return FilteredProducts
}

func(serv MyShopService) GetProducts() []Product {
  fmt.Println(">>> GetProducts")
  return Products
}

func(serv MyShopService) BuyProduct(id int) (resp BuyResponse) {
  fmt.Println(">>> BuyProduct: id = ", id)
  if id > len(Products) - 1 {
    serv.ResponseBuilder().SetResponseCode(404).Overide(true)
    return
  }

  if Products[id].Stock == 0 {
    serv.ResponseBuilder().SetResponseCode(409).Overide(true)
    return
  }

  if Products[id].Stock >0 {
    Products[id].Stock--
  }

  return_code := "order-accepted"
  redirect_url := ""
  if Products[id].VendorId != "" {
    elements := []string { "http://api.the-vendor.test:8080/api/vendor", Products[id].VendorId, "callback", Products[id].VendorProductId, "1" }
    callback := strings.Join(elements, "/")
    fmt.Println(">>> BuyProduct: firing callback to vendor", Products[id].VendorName, "with URL =", callback)

    client := &http.Client{}
    req, err := http.NewRequest("GET", callback, nil)
    resp, err := client.Do(req)
    if err != nil {
      fmt.Println(">>> BuyProduct: ERROR", err)
      return_code = "error"
      serv.ResponseBuilder().SetResponseCode(500).Overide(true)
    } else {
      defer resp.Body.Close()
      fmt.Println(">>> BuyProduct: response Status:", resp.Status)
      if resp.Status != "200 OK" {
        body, _ := ioutil.ReadAll(resp.Body)
        fmt.Println(">>> BuyProduct: response Body:", string(body))
        return_code = "error"
        serv.ResponseBuilder().SetResponseCode(500).Overide(true)
      } else {
        if Products[id].IsDigital {
          fmt.Println(">>> BuyProduct: Decoding JSON response")
          json_resp := new(CallbackResponse)
          json.NewDecoder(resp.Body).Decode(json_resp)
          redirect_url = json_resp.RedirectUrl;
        } else {
          fmt.Println(">>> BuyProduct: Ignoring JSON response")
        }
      }
    }
  }

  fmt.Println(">>> BuyProduct: return_code =", return_code, "redirect_url =", redirect_url)

  resp = BuyResponse { return_code, redirect_url }
  return
}


func(serv MyShopService) GetProduct(id int) (p Product) {
  fmt.Println(">>> GetProduct: id =", id)
  if id > len(Products) - 1 {
    serv.ResponseBuilder().SetResponseCode(404).Overide(true)
    return
  }
  p = Products[id]
  return
}

func(serv MyShopService) AddProduct(posted Product) {
  fmt.Println(">>> AddProduct: posted =", posted)
  id := len(Products)
  posted.Id = id
  Products = append(Products, posted);
  serv.ResponseBuilder().Created("/api/shop/product/"+string(id))
}

func(serv MyShopService) GetProductsByCategory(category string) []Product {
  fmt.Println(">>> GetProductsByCategory: category =", category)
  if (category == "") {
    return Products
  }

  FilteredProducts := []Product {}
  for _, p := range Products {
    if p.Category == category {
      FilteredProducts = append(FilteredProducts, p)
    }
  }

  return FilteredProducts
}

func(serv MyShopService) GetCategories() []Category {
  fmt.Println(">>> GetCategories")
  return Categories
}
