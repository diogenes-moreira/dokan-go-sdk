package types

import "time"

// ProductType represents the type of a product
type ProductType string

const (
	ProductTypeSimple   ProductType = "simple"
	ProductTypeGrouped  ProductType = "grouped"
	ProductTypeExternal ProductType = "external"
	ProductTypeVariable ProductType = "variable"
)

// ProductStatus represents the status of a product
type ProductStatus string

const (
	ProductStatusDraft   ProductStatus = "draft"
	ProductStatusPending ProductStatus = "pending"
	ProductStatusPublish ProductStatus = "publish"
)

// CatalogVisibility represents the catalog visibility of a product
type CatalogVisibility string

const (
	CatalogVisibilityVisible CatalogVisibility = "visible"
	CatalogVisibilityCatalog CatalogVisibility = "catalog"
	CatalogVisibilitySearch  CatalogVisibility = "search"
	CatalogVisibilityHidden  CatalogVisibility = "hidden"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusOnHold     OrderStatus = "on-hold"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
	OrderStatusFailed     OrderStatus = "failed"
)

// Product represents a Dokan product
type Product struct {
	ID                int                `json:"id,omitempty"`
	Name              string             `json:"name"`
	Slug              string             `json:"slug,omitempty"`
	Permalink         string             `json:"permalink,omitempty"`
	DateCreated       *time.Time         `json:"date_created,omitempty"`
	DateCreatedGMT    *time.Time         `json:"date_created_gmt,omitempty"`
	DateModified      *time.Time         `json:"date_modified,omitempty"`
	DateModifiedGMT   *time.Time         `json:"date_modified_gmt,omitempty"`
	Type              ProductType        `json:"type"`
	Status            ProductStatus      `json:"status"`
	Featured          bool               `json:"featured"`
	CatalogVisibility CatalogVisibility  `json:"catalog_visibility"`
	Description       string             `json:"description"`
	ShortDescription  string             `json:"short_description"`
	SKU               string             `json:"sku"`
	Price             string             `json:"price,omitempty"`
	RegularPrice      string             `json:"regular_price"`
	SalePrice         string             `json:"sale_price,omitempty"`
	DateOnSaleFrom    *time.Time         `json:"date_on_sale_from,omitempty"`
	DateOnSaleFromGMT *time.Time         `json:"date_on_sale_from_gmt,omitempty"`
	DateOnSaleTo      *time.Time         `json:"date_on_sale_to,omitempty"`
	DateOnSaleToGMT   *time.Time         `json:"date_on_sale_to_gmt,omitempty"`
	PriceHTML         string             `json:"price_html,omitempty"`
	OnSale            bool               `json:"on_sale,omitempty"`
	Purchasable       bool               `json:"purchasable,omitempty"`
	TotalSales        int                `json:"total_sales,omitempty"`
	Virtual           bool               `json:"virtual"`
	Downloadable      bool               `json:"downloadable"`
	Categories        []ProductCategory  `json:"categories,omitempty"`
	Tags              []ProductTag       `json:"tags,omitempty"`
	Images            []ProductImage     `json:"images,omitempty"`
	Attributes        []ProductAttribute `json:"attributes,omitempty"`
	DefaultAttributes []ProductAttribute `json:"default_attributes,omitempty"`
	Variations        []int              `json:"variations,omitempty"`
	GroupedProducts   []int              `json:"grouped_products,omitempty"`
	MenuOrder         int                `json:"menu_order"`
	MetaData          []MetaData         `json:"meta_data,omitempty"`
}

// ProductCategory represents a product category
type ProductCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ProductTag represents a product tag
type ProductTag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ProductImage represents a product image
type ProductImage struct {
	ID       int    `json:"id,omitempty"`
	Src      string `json:"src"`
	Name     string `json:"name,omitempty"`
	Alt      string `json:"alt,omitempty"`
	Position int    `json:"position,omitempty"`
}

// ProductAttribute represents a product attribute
type ProductAttribute struct {
	ID        int      `json:"id,omitempty"`
	Name      string   `json:"name"`
	Position  int      `json:"position,omitempty"`
	Visible   bool     `json:"visible"`
	Variation bool     `json:"variation"`
	Options   []string `json:"options"`
}

// Order represents a Dokan order
type Order struct {
	ID                 int            `json:"id,omitempty"`
	ParentID           int            `json:"parent_id,omitempty"`
	Number             string         `json:"number,omitempty"`
	OrderKey           string         `json:"order_key,omitempty"`
	CreatedVia         string         `json:"created_via,omitempty"`
	Version            string         `json:"version,omitempty"`
	Status             OrderStatus    `json:"status"`
	Currency           string         `json:"currency"`
	DateCreated        *time.Time     `json:"date_created,omitempty"`
	DateCreatedGMT     *time.Time     `json:"date_created_gmt,omitempty"`
	DateModified       *time.Time     `json:"date_modified,omitempty"`
	DateModifiedGMT    *time.Time     `json:"date_modified_gmt,omitempty"`
	DiscountTotal      string         `json:"discount_total,omitempty"`
	DiscountTax        string         `json:"discount_tax,omitempty"`
	ShippingTotal      string         `json:"shipping_total,omitempty"`
	ShippingTax        string         `json:"shipping_tax,omitempty"`
	CartTax            string         `json:"cart_tax,omitempty"`
	Total              string         `json:"total,omitempty"`
	TotalTax           string         `json:"total_tax,omitempty"`
	PricesIncludeTax   bool           `json:"prices_include_tax,omitempty"`
	CustomerID         int            `json:"customer_id,omitempty"`
	CustomerIPAddress  string         `json:"customer_ip_address,omitempty"`
	CustomerUserAgent  string         `json:"customer_user_agent,omitempty"`
	CustomerNote       string         `json:"customer_note,omitempty"`
	Billing            *Address       `json:"billing,omitempty"`
	Shipping           *Address       `json:"shipping,omitempty"`
	PaymentMethod      string         `json:"payment_method,omitempty"`
	PaymentMethodTitle string         `json:"payment_method_title,omitempty"`
	TransactionID      string         `json:"transaction_id,omitempty"`
	DatePaid           *time.Time     `json:"date_paid,omitempty"`
	DatePaidGMT        *time.Time     `json:"date_paid_gmt,omitempty"`
	DateCompleted      *time.Time     `json:"date_completed,omitempty"`
	DateCompletedGMT   *time.Time     `json:"date_completed_gmt,omitempty"`
	CartHash           string         `json:"cart_hash,omitempty"`
	LineItems          []LineItem     `json:"line_items,omitempty"`
	TaxLines           []TaxLine      `json:"tax_lines,omitempty"`
	ShippingLines      []ShippingLine `json:"shipping_lines,omitempty"`
	FeeLines           []FeeLine      `json:"fee_lines,omitempty"`
	CouponLines        []CouponLine   `json:"coupon_lines,omitempty"`
	Refunds            []Refund       `json:"refunds,omitempty"`
	MetaData           []MetaData     `json:"meta_data,omitempty"`
}

// Address represents a billing or shipping address
type Address struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Company   string `json:"company,omitempty"`
	Address1  string `json:"address_1"`
	Address2  string `json:"address_2,omitempty"`
	City      string `json:"city"`
	State     string `json:"state"`
	Postcode  string `json:"postcode"`
	Country   string `json:"country"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

// LineItem represents an order line item
type LineItem struct {
	ID          int        `json:"id,omitempty"`
	Name        string     `json:"name"`
	ProductID   int        `json:"product_id"`
	VariationID int        `json:"variation_id,omitempty"`
	Quantity    int        `json:"quantity"`
	TaxClass    string     `json:"tax_class,omitempty"`
	Subtotal    string     `json:"subtotal"`
	SubtotalTax string     `json:"subtotal_tax"`
	Total       string     `json:"total"`
	TotalTax    string     `json:"total_tax"`
	Taxes       []TaxLine  `json:"taxes,omitempty"`
	MetaData    []MetaData `json:"meta_data,omitempty"`
	SKU         string     `json:"sku,omitempty"`
	Price       float64    `json:"price,omitempty"`
}

// TaxLine represents a tax line
type TaxLine struct {
	ID               int        `json:"id,omitempty"`
	RateCode         string     `json:"rate_code"`
	RateID           int        `json:"rate_id"`
	Label            string     `json:"label"`
	Compound         bool       `json:"compound"`
	TaxTotal         string     `json:"tax_total"`
	ShippingTaxTotal string     `json:"shipping_tax_total"`
	MetaData         []MetaData `json:"meta_data,omitempty"`
}

// ShippingLine represents a shipping line
type ShippingLine struct {
	ID          int        `json:"id,omitempty"`
	MethodTitle string     `json:"method_title"`
	MethodID    string     `json:"method_id"`
	Total       string     `json:"total"`
	TotalTax    string     `json:"total_tax"`
	Taxes       []TaxLine  `json:"taxes,omitempty"`
	MetaData    []MetaData `json:"meta_data,omitempty"`
}

// FeeLine represents a fee line
type FeeLine struct {
	ID        int        `json:"id,omitempty"`
	Name      string     `json:"name"`
	TaxClass  string     `json:"tax_class,omitempty"`
	TaxStatus string     `json:"tax_status"`
	Total     string     `json:"total"`
	TotalTax  string     `json:"total_tax"`
	Taxes     []TaxLine  `json:"taxes,omitempty"`
	MetaData  []MetaData `json:"meta_data,omitempty"`
}

// CouponLine represents a coupon line
type CouponLine struct {
	ID          int        `json:"id,omitempty"`
	Code        string     `json:"code"`
	Discount    string     `json:"discount"`
	DiscountTax string     `json:"discount_tax"`
	MetaData    []MetaData `json:"meta_data,omitempty"`
}

// Refund represents a refund
type Refund struct {
	ID     int    `json:"id"`
	Reason string `json:"reason,omitempty"`
	Total  string `json:"total"`
}

// Store represents a Dokan store
type Store struct {
	ID             int                          `json:"id"`
	StoreName      string                       `json:"store_name"`
	FirstName      string                       `json:"first_name"`
	LastName       string                       `json:"last_name"`
	Email          string                       `json:"email"`
	Phone          string                       `json:"phone,omitempty"`
	ShowEmail      bool                         `json:"show_email,omitempty"`
	Address        *Address                     `json:"address,omitempty"`
	Location       string                       `json:"location,omitempty"`
	Banner         string                       `json:"banner,omitempty"`
	Icon           string                       `json:"icon,omitempty"`
	Gravatar       string                       `json:"gravatar,omitempty"`
	ShopURL        string                       `json:"shop_url,omitempty"`
	ProductsURL    string                       `json:"products_url,omitempty"`
	TocsURL        string                       `json:"tocs_url,omitempty"`
	Featured       bool                         `json:"featured,omitempty"`
	Rating         *Rating                      `json:"rating,omitempty"`
	Enabled        bool                         `json:"enabled,omitempty"`
	Registered     string                       `json:"registered,omitempty"`
	PaymentMethods map[string]map[string]string `json:"payment,omitempty"`
	Social         map[string]string            `json:"social,omitempty"`
}

// Rating represents store rating information
type Rating struct {
	Rating string `json:"rating"`
	Count  int    `json:"count"`
}

// MetaData represents metadata
type MetaData struct {
	ID    int         `json:"id,omitempty"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// ListParams represents common list parameters
type ListParams struct {
	Page    int    `url:"page,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
	Search  string `url:"search,omitempty"`
	OrderBy string `url:"orderby,omitempty"`
	Order   string `url:"order,omitempty"`
}

// ListResponse represents a paginated list response
type ListResponse struct {
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
}

// ProductListParams represents parameters for listing products
type ProductListParams struct {
	ListParams
	Status      []ProductStatus `url:"status,omitempty"`
	Type        []ProductType   `url:"type,omitempty"`
	Featured    *bool           `url:"featured,omitempty"`
	Category    []int           `url:"category,omitempty"`
	Tag         []int           `url:"tag,omitempty"`
	MinPrice    *float64        `url:"min_price,omitempty"`
	MaxPrice    *float64        `url:"max_price,omitempty"`
	StockStatus string          `url:"stock_status,omitempty"`
	SKU         string          `url:"sku,omitempty"`
}

// OrderListParams represents parameters for listing orders
type OrderListParams struct {
	ListParams
	Status         []OrderStatus `url:"status,omitempty"`
	Customer       int           `url:"customer,omitempty"`
	Product        int           `url:"product,omitempty"`
	After          *time.Time    `url:"after,omitempty"`
	Before         *time.Time    `url:"before,omitempty"`
	ModifiedAfter  *time.Time    `url:"modified_after,omitempty"`
	ModifiedBefore *time.Time    `url:"modified_before,omitempty"`
}

// StoreListParams represents parameters for listing stores
type StoreListParams struct {
	ListParams
	Featured *bool `url:"featured,omitempty"`
	Enabled  *bool `url:"enabled,omitempty"`
}
