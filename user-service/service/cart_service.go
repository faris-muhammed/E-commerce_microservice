package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/faris-muhammed/e-commerce/user-service/models"
	"github.com/faris-muhammed/e-commerce/user-service/repository"
	cartpb "github.com/faris-muhammed/e-protofiles/cart"
	productpb "github.com/faris-muhammed/e-protofiles/product"
)

type CartServiceServer struct {
	cartpb.UnimplementedCartServiceServer
	repo           repository.CartRepository
	productService productpb.ProductServiceClient
}

func NewCartService(repo repository.CartRepository, productService productpb.ProductServiceClient) *CartServiceServer {
	return &CartServiceServer{
		repo:           repo,
		productService: productService,
	}
}

func (s *CartServiceServer) AddToCart(ctx context.Context, req *cartpb.AddToCartRequest) (*cartpb.AddToCartResponse, error) {
	item := models.Cart{
		UserID:    req.UserId,
		ProductID: req.ProductId,
		Quantity:  req.Quantity,
	}

	// Create the request for GetProduct
	getProductReq := &productpb.GetProductRequest{
		ProductId: uint64(req.ProductId),
	}

	// Check if the product exists by calling the Seller Service
	productResponse, err := s.productService.GetProduct(ctx, getProductReq)
	if err != nil {
		return nil, err // Return the error if checking the product fails
	}

	// If the product does not exist
	if productResponse == nil || productResponse.Product == nil {
		return nil, errors.New("product does not exist")
	}

	// Add item to the cart
	if err := s.repo.AddToCart(item); err != nil {
		return nil, errors.New("failed to add item to cart")
	}

	return &cartpb.AddToCartResponse{
		Message: "Item added to cart successfully",
	}, nil
}

func (s *CartServiceServer) EditCart(ctx context.Context, req *cartpb.EditCartRequest) (*cartpb.EditCartResponse, error) {
	if err := s.repo.EditCartItem(req.UserId, req.ProductId, req.Quantity); err != nil {
		return nil, errors.New("failed to edit cart item")
	}

	return &cartpb.EditCartResponse{
		Message: "Cart item updated successfully",
	}, nil
}

func (s *CartServiceServer) RemoveCart(ctx context.Context, req *cartpb.RemoveCartRequest) (*cartpb.RemoveCartResponse, error) {
	if err := s.repo.RemoveCartItem(req.UserId); err != nil {
		return nil, errors.New("failed to remove cart item")
	}

	return &cartpb.RemoveCartResponse{
		Message: "Item removed from cart successfully",
	}, nil
}

func (s *CartServiceServer) RemoveCartByID(ctx context.Context, req *cartpb.RemoveCartByIDRequest) (*cartpb.RemoveCartByIDResponse, error) {
	// Call the repository to delete the cart item for the given user and product
	err := s.repo.RemoveCartByID(req.UserId, req.ProductId)
	if err != nil {
		fmt.Printf("Error removing cart item: %v\n", err)
		return nil, errors.New("failed to remove cart item")
	}

	// Return a success response
	return &cartpb.RemoveCartByIDResponse{
		Message: fmt.Sprintf("Cart item for user %d and product %d removed successfully", req.UserId, req.ProductId),
	}, nil
}


func (s *CartServiceServer) ListCart(ctx context.Context, req *cartpb.ListCartRequest) (*cartpb.ListCartResponse, error) {
	// Fetch cart items for the user
	items, err := s.repo.ListCartItems(req.UserId)
	if err != nil {
		fmt.Printf("Error fetching cart items: %v\n", err)
		return nil, errors.New("failed to list cart items")
	}

	var cartItems []*cartpb.CartItem
	var total float64 // Variable to calculate the total price of the cart

	for _, item := range items {
		// Fetch product details from the product service
		getProductReq := &productpb.GetProductRequest{
			ProductId: uint64(item.ProductID),
		}

		productResponse, err := s.productService.GetProduct(ctx, getProductReq)
		if err != nil {
			// Log the error but continue processing other items
			fmt.Printf("Error fetching product details for ProductID %d: %v\n", item.ProductID, err)
			continue
		}

		// Ensure the product response is valid
		if productResponse != nil && productResponse.Product != nil {
			product := productResponse.Product

			// Calculate item total price
			itemTotalPrice := product.Price * float64(item.Quantity)

			// Add the item to the cart
			cartItems = append(cartItems, &cartpb.CartItem{
				Quantity: item.Quantity,
				Product: &cartpb.Product{
					Id:                product.Id,
					Name:              product.Name,
					AvailableQuantity: uint32(product.Stock),
					Price:             product.Price,
				},
				SubTotal: itemTotalPrice, // Include item total price
			})
			fmt.Println("ProductStock",product.Stock)
			// Accumulate the total cart price
			total += itemTotalPrice
		}
	}

	// Return the response with cart items and total price
	return &cartpb.ListCartResponse{
		Items: cartItems,
		TotalAmount: total,
	}, nil
}


// func (s *CartServiceServer) ListProducts(ctx context.Context, req *cartpb.ListProductsRequest) (*cartpb.ListProductsResponse, error) {
// 	// Forward request to seller-service
// 	productResponse, err := s.productService.ListProducts(ctx, &productpb.ListProductsRequest{})
// 	if err != nil {
// 		fmt.Printf("Error from seller service: %v\n", err)
// 		return nil, errors.New("failed to list products from seller-service")
// 	}
// 	fmt.Printf("Products received from seller-service: %+v\n", productResponse.Products)
// 	// Map response to user-service response format
// 	var productItems []*cartpb.Product
// 	for _, product := range productResponse.Products {
// 		productItems = append(productItems, &cartpb.Product{
// 			Id:                product.Id,
// 			Name:              product.Name,
// 			Price:             product.Price,
// 			AvailableQuantity: uint32(product.Stock),
// 		})
// 	}
// 	fmt.Printf("Mapped products: %+v\n", productItems)

// 	return &cartpb.ListProductsResponse{
// 		Products: productItems,
// 	}, nil
// }
