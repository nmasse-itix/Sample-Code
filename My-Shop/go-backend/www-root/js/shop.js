var apibase = "/api/shop";
var dialog;

require([ "dojo/ready", "dojo/request/xhr", "dojo/dom-construct", "dojo/dom", "dojo/on", "dojo/dom-style", "dojo/dom-attr", "dijit/Dialog" ],
        function(ready, xhr, domConstruct, dom, on, domStyle, domAttr, Dialog) {

  function setOnCategoryClickHandler(node, catid, catname) {
    on(node, "click", function() {
      loadProducts(catid);
      setSubTitle(catname);
    });
  }

  function setOnProductClickHandler(node, id) {
    on(node, "click", function() {
      displayProduct(id);
    });
  }

  function setSearchHandler() {
    on(dom.byId("search_textbox"), "keyup", doSearch);
    on(dom.byId("search_textbox"), "blur", function () {
      console.log("SEARCH TEXT BOX >> Blur");
      window.setInterval(function () { console.log("SEARCH TEXT BOX >> Je la cache"); domStyle.set("search_results", "visibility", "hidden"); }, 100);
    });
    on(dom.byId("search_textbox"), "focus", function () {
      console.log("SEARCH TEXT BOX >> Focus");
      var searchCriteria = domAttr.get("search_textbox", "value");
      if (searchCriteria != "") {
        domStyle.set("search_results", "visibility", "visible");
      }
    });
  }

  function setBuyHandler(node, id) {
    on(node, "click", function() {
      buyProduct(id);
    });
  }

  function buyProduct(id) {
    xhr(apibase + "/product/" + encodeURI(id) + "/buy",
      { handleAs: "json" }
    ).then(function (data) {
        dialog.set("title", "Commande acceptée");
        dialog.set("content", "La commande est partie. Vous allez très prochainement recevoir le produit.");
        console.log("test...");
        if (data.DownloadUrl != null && data.DownloadUrl != "") {
          console.log("popup !");
          window.location.href = data.DownloadUrl;
        }
        dialog.show();
        displayProduct(id); // Refresh UI
      }, function (err) {
        if (err.response != null && err.response.status == 409) {
          dialog.set("title", "Erreur");
          dialog.set("content", "Le produit n'est plus en stock. Désolé.");
          dialog.show();
        } else {
          dialog.set("title", "OOPS");
          dialog.set("content", "Erreur interne. Désolé.");
          dialog.show();
        }
        console.log(err);
        displayProduct(id); // Refresh UI
      }, function (evt) {

    });
  }

  function doSearch(evt) {
    var searchCriteria = domAttr.get("search_textbox", "value");
    if (searchCriteria == "") {
      domStyle.set("search_results", "visibility", "hidden");
      return;
    } else {
      domStyle.set("search_results", "visibility", "visible");
    }

    xhr(apibase + "/search/" + encodeURI(searchCriteria),
      { handleAs: "json" }
    ).then(function (data) {
        domConstruct.empty("search_results");
        var placeholder = dom.byId("search_results");
        for (var i = 0; i < data.length; i++) {
          var div = domConstruct.create("div", {}, placeholder);
          domConstruct.create("img", { width: "32px", src: "/img/" + data[i].Image }, div);
          domConstruct.create("span", { textContent: data[i].Name }, div);
          setOnProductClickHandler(div, data[i].Id);
        }
        if (data.length == 0) {
          domConstruct.create("span", { textContent: "Aucun résultat", 'class': "no_result" }, placeholder);
        }
      }, function (err) {
        console.log(err);
      }, function (evt) {

    });
  }

  function displayProduct(id) {
    xhr(apibase + "/product/" + encodeURI(id),
      { handleAs: "json" }
    ).then(function (data) {
        domConstruct.empty("products_pane");
        var placeholder = dom.byId("products_pane");
        var div = dojo.create("div", { 'class': 'product_detail' }, placeholder);
        dojo.create("h1", { 'class': "product_name", textContent: data.Name }, div);
        var div2 = dojo.create("div", {}, div)
        dojo.create("img", { src: "/img/" + data.Image, height: "200px" }, div2);
        dojo.create("span", { 'class': "product_price", textContent: data.Price + " €" }, div2);
        if (data.Stock != "-1") {
          dojo.create("span", { 'class': "product_stock", textContent: "En Stock: " + data.Stock }, div2);
        }
        if (data.VendorName != "") {
          dojo.create("span", { 'class': "sold_by", textContent: "Vendu par: " + data.VendorName }, div2);
        }

        var buy_button = dojo.create("div", { 'class': "buy_button" }, div2);
        dojo.create("div", { textContent: "Acheter !" }, buy_button);
        setBuyHandler(buy_button, data.Id);

        dojo.create("div", { 'class': "product_description", textContent: data.Description }, div);
      }, function (err) {
        if (err.response != null && err.response.status == 404) {
          dialog.set("title", "Erreur");
          dialog.set("content", "Le produit a été retiré de la vente. Désolé.");
          dialog.show();
        } else {
          dialog.set("title", "OOPS");
          dialog.set("content", "Erreur interne. Désolé.");
          dialog.show();
        }
        console.log(err);
      }, function (evt) {

    });
  }

  function setSubTitle(name) {
    domConstruct.empty("subtitle");
    dom.byId("subtitle").textContent = name;
  }

  function loadProducts(catid) {
    var queryString = "";
    if (catid != null) {
      queryString = "?category=" + encodeURI(catid);
    }
    xhr(apibase + "/product/" + queryString,
      { handleAs: "json" }
    ).then(function (data) {
        domConstruct.empty("products_pane");
        var placeholder = dom.byId("products_pane");
        var table = domConstruct.create("table", { 'class': "product_table" }, placeholder);
        var current_tr = null;
        var i = 0;
        for (; i < data.length; i++) {
          if (i % 3 == 0) {
            current_tr = domConstruct.create("tr", {}, table);
          }
          var td = domConstruct.create("td", {}, current_tr);
          domConstruct.create("img", { width: "100px", src: "/img/" + data[i].Image }, td);
          domConstruct.create("span", { textContent: data[i].Name }, td);
          setOnProductClickHandler(td, data[i].Id);
        }
        // Fill remaining columns if less than 3 products
        for (; i < 3; i++) {
          domConstruct.create("td", {}, current_tr);
        }
      }, function (err) {
        dialog.set("title", "OOPS");
        dialog.set("content", "Erreur interne. Désolé.");
        dialog.show();
        console.log(err);
      }, function (evt) {

    });
  }

  function loadCategories() {
    xhr(apibase + "/category/",
      { handleAs: "json" }
    ).then(function (data) {
      var placeholder = dom.byId("category_list_placeholder");
      domConstruct.empty("category_list_placeholder");
      var all_products_node = domConstruct.create("span", { textContent: "Tous les produits", 'class': "category_item" }, placeholder);
      setOnCategoryClickHandler(all_products_node, null, "Tous les produits");
      for (var i = 0; i < data.length; i++) {
        var catid = data[i].Id;
        var catname = data[i].Name;
        var node = domConstruct.create("span", { textContent: data[i].Name, 'class': "category_item" }, placeholder);
        setOnCategoryClickHandler(node, catid, catname);
      }
    }, function (err) {
      dialog.set("title", "OOPS");
      dialog.set("content", "Erreur interne. Désolé.");
      dialog.show();
      console.log(err);
    }, function (evt) {

    });
  }

  ready(function() {
    loadCategories();
    loadProducts(null);
    setSubTitle("Tous les produits");
    setSearchHandler();
    dialog = new Dialog({
      id: "global_dialog",
      title: "...",
      content: "...",
      style: "width: 500px; display: none;"
    });
  });
});
