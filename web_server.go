package main

import (
	"net/http"
	"strconv"

	"github.com/NikhilSharmaWe/marketplace/proto"
	"github.com/julienschmidt/httprouter"
)

func (app *application) setupWebserver() {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)

	router.HandlerFunc(http.MethodPost, "/user", app.HandleCreateUser)
	router.HandlerFunc(http.MethodGet, "/user/:id", app.HandleGetUser)
	router.HandlerFunc(http.MethodGet, "/neighbour/:userId", app.HandleGetNearestNeighbour)

	router.HandlerFunc(http.MethodPost, "/shop", app.HandleCreateShop)
	router.HandlerFunc(http.MethodGet, "/shop/:id", app.HandleGetShop)
	router.HandlerFunc(http.MethodGet, "/shopForUser/:userId/:maxDist", app.HandleGetShopForUser)

	router.HandlerFunc(http.MethodPost, "/product", app.HandleCreateProduct)
	router.HandlerFunc(http.MethodGet, "/product/:id", app.HandleGetProduct)
	router.HandlerFunc(http.MethodPost, "/addProduct", app.HandleAddProductToShop)
	router.HandlerFunc(http.MethodGet, "/serviceableProducts/:shopId", app.HandleGetServiceableProducts)

	router.HandlerFunc(http.MethodGet, "/inventory/:shopId/:productId", app.HandleGetInventory)
	router.HandlerFunc(http.MethodPost, "/inventory", app.HandleUpdateInventory)

	app.goApiBoot.WebServer.Handler = router
}

func (app *application) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	createReq := &proto.CreateUserRequest{}
	err := app.readJSON(w, r, createReq)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	user, err := app.grpcClient.CreateUser(app.ctx, createReq)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	app.logger.Printf("new user created %+v", user)
	app.writeJSON(w, http.StatusOK, map[string]any{"created": "user"})
}

func (app *application) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id := params.ByName("id")

	protoReq := proto.GetRequest{
		Id: id,
	}

	user, err := app.grpcClient.GetUserByID(app.ctx, &protoReq)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	app.logger.Printf("got user[%s]", id)
	app.writeJSON(w, http.StatusOK, map[string]any{
		"user": user,
	})
}

func (app *application) HandleCreateShop(w http.ResponseWriter, r *http.Request) {
	createReq := &proto.CreateShopRequest{}
	err := app.readJSON(w, r, createReq)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	shop, err := app.grpcClient.CreateShop(app.ctx, createReq)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	app.logger.Printf("new shop created %+v", shop)
	app.writeJSON(w, http.StatusOK, map[string]any{"created": "shop"})

}

func (app *application) HandleGetShop(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id := params.ByName("id")

	protoReq := proto.GetRequest{
		Id: id,
	}

	shop, err := app.grpcClient.GetShopByID(app.ctx, &protoReq)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	app.logger.Printf("got shop[%s]", id)
	app.writeJSON(w, http.StatusOK, map[string]any{
		"shop": shop,
	})
}

func (app *application) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	createReq := &proto.CreateProductRequest{}
	err := app.readJSON(w, r, createReq)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	product, err := app.grpcClient.CreateProduct(app.ctx, createReq)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	app.logger.Printf("new product created %+v", product)
	app.writeJSON(w, http.StatusOK, map[string]any{"created": "product"})
}

func (app *application) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id := params.ByName("id")

	protoReq := proto.GetRequest{
		Id: id,
	}

	product, err := app.grpcClient.GetProductByID(app.ctx, &protoReq)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	app.logger.Printf("got product[%s]", id)
	app.writeJSON(w, http.StatusOK, map[string]any{
		"product": product,
	})
}

func (app *application) HandleAddProductToShop(w http.ResponseWriter, r *http.Request) {
	addProduct := &proto.AddServiceableProductRequest{}
	err := app.readJSON(w, r, addProduct)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	shop, err := app.grpcClient.AddServiceableProduct(app.ctx, addProduct)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	app.logger.Printf("added new product[%s] to shop[%s]", addProduct.ProductId, addProduct.ShopId)
	app.writeJSON(w, http.StatusOK, map[string]any{
		"added product": map[string]any{
			"productAddedId": addProduct.ProductId,
			"shop":           shop,
		},
	})
}

func (app *application) HandleGetServiceableProducts(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id := params.ByName("shopId")

	protoReq := proto.GetServiceableProductsRequest{
		ShopId: id,
	}

	products, err := app.grpcClient.GetServiceableProducts(app.ctx, &protoReq)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	app.logger.Printf("got serviciable products for shop[%s]", id)
	app.writeJSON(w, http.StatusOK, products)
}

func (app *application) HandleGetInventory(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	shopId := params.ByName("shopId")
	productId := params.ByName("productId")

	protoReq := proto.GetInventoryRequest{
		ShopId:    shopId,
		ProductId: productId,
	}

	inventory, err := app.grpcClient.GetInventory(app.ctx, &protoReq)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if inventory.Quantity == 0 {
		val := 0
		inventory.Quantity = int32(val)
	}

	app.logger.Printf("got inventory of product[%s] for shop[%s]", productId, shopId)
	app.writeJSON(w, http.StatusOK, inventory)
}

func (app *application) HandleUpdateInventory(w http.ResponseWriter, r *http.Request) {
	updateReq := &proto.UpdateInventoryRequest{}

	err := app.readJSON(w, r, updateReq)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	inventory, err := app.grpcClient.UpdateInventory(app.ctx, updateReq)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	app.logger.Printf("updated inventory of product[%s] for shop[%s]", updateReq.ProductId, updateReq.ShopId)
	app.writeJSON(w, http.StatusOK, map[string]any{
		"updatedInventory": inventory,
	})
}

func (app *application) HandleGetNearestNeighbour(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	userId := params.ByName("userId")

	protoReq := proto.GetNearestNeighbourRequest{
		UserId: userId,
	}

	neighbour, err := app.grpcClient.GetNearestNeighbour(app.ctx, &protoReq)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	app.logger.Printf("got the nearest neighbour for user[%s]", userId)
	app.writeJSON(w, http.StatusOK, map[string]any{
		"nearestNeighbour": neighbour,
	})
}

func (app *application) HandleGetShopForUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	userId := params.ByName("userId")
	maxDist, err := strconv.ParseFloat(params.ByName("maxDist"), 64)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	protoReq := proto.GetShopForUserRequest{
		UserId:          userId,
		MaxDistanceInKM: maxDist,
	}

	shop, err := app.grpcClient.GetShopForUser(app.ctx, &protoReq)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	app.logger.Printf("got the shops for user[%s] preference proximity", userId)
	app.writeJSON(w, http.StatusOK, map[string]any{
		"nearbyShops": shop,
	})
}
