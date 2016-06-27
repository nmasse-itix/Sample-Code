curl -H "Content-Type: application/json" \
     -d '{ "Name": "A New Product", "Category": "maison", "Image": "logo.png", "Description": "My brand new product", "Price": 1, "Stock": 1 }' \
     -X POST \
     -D - \
     http://localhost:8787/api/shop/product/

