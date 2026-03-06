# Trendyol Product Scraper API (Go)

Bu proje, Python ile yazdığınız kazıma (scraping) mantığının Go diline dönüştürülmüş ve bir API endpoint'i haline getirilmiş halidir. Bulut sunucunuzda Docker kullanarak kolayca çalıştırabilirsiniz.

## Özellikler
- **Hızlı ve Hafif**: Go dilinin performansı ile düşük kaynak tüketimi.
- **REST API**: `/scrape?url=...` endpoint'i üzerinden veri döner.
- **Docker Desteği**: Dağıtım için hazır Dockerfile.

## Nasıl Çalıştırılır?

### 1. Yerel Olarak (Go Yüklü İse)
```bash
cd go_scraper
go run main.go
```

### 2. Docker İle (Tavsiye Edilen)
Bulut sunucunuzda Docker yüklü ise:
```bash
# Image'ı oluştur
docker build -t trendyol-scraper .

# Konteyner'ı çalıştır
docker run -p 8080:8080 trendyol-scraper
```

## Kullanım
API çalıştıktan sonra tarayıcıdan veya `curl` ile test edebilirsiniz:

```bash
curl "http://localhost:8080/scrape?url=https://www.trendyol.com/brand/product-p-12345"
```

## Yanıt Formatı (JSON)
```json
{
  "id": 123456,
  "name": "Ürün Adı",
  "brand": "Marka",
  "price": "100.00 TL",
  "merchant": "Satıcı Adı",
  "category": "Kadın",
  "stock_status": "Var",
  "images": ["url1", "url2"]
}
```
