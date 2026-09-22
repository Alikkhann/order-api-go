package order

import "errors"

 var ErrOrderNotFound = errors.New("Order not found")
 var ErrUserNotFound = errors.New("User not found")
 var ErrNoProductsFound = errors.New("no products found")