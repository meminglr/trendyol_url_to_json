package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProductResponse struct {
	ID          interface{} `json:"id"`
	Name        string      `json:"name"`
	Brand       string      `json:"brand"`
	Price       string      `json:"price"`
	Merchant    string      `json:"merchant"`
	Category    string      `json:"category"`
	StockStatus string      `json:"stock_status"`
	Images      []string    `json:"images"`
}

func main() {
	r := gin.Default()

	r.GET("/scrape", func(c *gin.Context) {
		productURL := c.Query("url")
		if productURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "url parametresi eksik"})
			return
		}

		data, err := scrapeProduct(productURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, data)
	})

	fmt.Println("Sunucu 8080 portunda çalışıyor...")
	r.Run(":8080")
}

func scrapeProduct(url string) (*ProductResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "tr-TR,tr;q=0.9,en-US;q=0.8,en;q=0.7")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("sayfa yüklenemedi: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	content := string(body)
	prefix := "window[\"__envoy_product-info__PROPS\"]="
	startPos := strings.Index(content, prefix)

	if startPos == -1 {
		return nil, fmt.Errorf("ürün verisi bulunamadı")
	}

	jsonStart := content[startPos+len(prefix):]
	
	// We need to find the end of the JSON object. 
	// Since it's followed by a script tag close or another variable, 
	// and raw_decode in python handles it, we can use a simpler approach or a proper JSON parser.
	// Go's json.Unmarshal will fail if there is trailing data.
	// We'll find the first ';' after the JSON start if it exists, or use the whole thing if it's fine.
	
	// A more robust way in Go is to use a decoder that stops at the end of the object.
	var rawData map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(jsonStart))
	if err := decoder.Decode(&rawData); err != nil {
		return nil, fmt.Errorf("JSON ayrıştırma hatası: %v", err)
	}

	product, ok := rawData["product"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("ürün detayı bulunamadı")
	}

	// Extract nested fields safely
	res := &ProductResponse{
		ID:   product["id"],
		Name: getString(product, "name"),
	}

	if brand, ok := product["brand"].(map[string]interface{}); ok {
		res.Brand = getString(brand, "name")
	}

	if merchantListing, ok := product["merchantListing"].(map[string]interface{}); ok {
		if winner, ok := merchantListing["winnerVariant"].(map[string]interface{}); ok {
			if price, ok := winner["price"].(map[string]interface{}); ok {
				res.Price = getString(price, "discountedPrice")
				// Sometimes it's a map, python script showed .get("text")
				if priceText, ok := price["discountedPrice"].(map[string]interface{}); ok {
					res.Price = getString(priceText, "text")
				}
			}
		}
		if merchant, ok := merchantListing["merchant"].(map[string]interface{}); ok {
			res.Merchant = getString(merchant, "name")
		}
	}

	if category, ok := product["category"].(map[string]interface{}); ok {
		res.Category = getString(category, "name")
	}

	if inStock, ok := product["inStock"].(bool); ok && inStock {
		res.StockStatus = "Var"
	} else {
		res.StockStatus = "Yok"
	}

	if images, ok := product["images"].([]interface{}); ok {
		for _, img := range images {
			if s, ok := img.(string); ok {
				res.Images = append(res.Images, s)
			}
		}
	}

	return res, nil
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
