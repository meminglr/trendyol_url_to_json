import requests
import json
import sys

def scrape_product(url):
    headers = {
        "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
        "Accept-Language": "tr-TR,tr;q=0.9,en-US;q=0.8,en;q=0.7",
    }

    try:
        response = requests.get(url, headers=headers)
        if response.status_code != 200:
            print(f"Hata: Sayfa yüklenemedi (Durum Kodu: {response.status_code})")
            return None
        
        text = response.text
        prefix = 'window["__envoy_product-info__PROPS"]='
        start_pos = text.find(prefix)
        
        if start_pos == -1:
            print("Hata: Ürün verisi sayfa içerisinde bulunamadı.")
            return None
            
        start_json = start_pos + len(prefix)
        decoder = json.JSONDecoder()
        data, _ = decoder.raw_decode(text[start_json:])
        product = data.get("product", {})
        
        # Verileri ayıkla
        extracted = {
            "ID": product.get("id"),
            "İsim": product.get("name"),
            "Marka": product.get("brand", {}).get("name"),
            "Fiyat": product.get("merchantListing", {}).get("winnerVariant", {}).get("price", {}).get("discountedPrice", {}).get("text"),
            "Satıcı": product.get("merchantListing", {}).get("merchant", {}).get("name"),
            "Kategori": product.get("category", {}).get("name"),
            "Stok Durumu": "Var" if product.get("inStock") else "Yok",
            "Resimler": product.get("images", [])
        }
        
        return extracted

    except Exception as e:
        print(f"Bir hata oluştu: {e}")
        return None

if __name__ == "__main__":
    if len(sys.argv) > 1:
        product_url = sys.argv[1]
    else:
        product_url = input("Ürün linkini giriniz: ")
    
    result = scrape_product(product_url)
    if result:
        print("\n--- Ürün Bilgileri ---")
        for key, value in result.items():
            if key == "Resimler":
                print(f"{key}: {len(value)} adet bulundu.")
            else:
                print(f"{key}: {value}")
        
        # Save to file optionally
        with open("product_info.json", "w", encoding="utf-8") as f:
            json.dump(result, f, indent=2, ensure_ascii=False)
        print("\nBilgiler 'product_info.json' dosyasına kaydedildi.")
