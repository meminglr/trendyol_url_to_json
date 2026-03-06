from bs4 import BeautifulSoup
import re

with open("product_page.html", "r", encoding="utf-8") as f:
    soup = BeautifulSoup(f, 'html.parser')

scripts = soup.find_all('script')
for idx, s in enumerate(scripts):
    content = s.string if s.string else ""
    if len(content) > 1000:
        print(f"Script {idx}: length {len(content)}, id: {s.get('id')}, type: {s.get('type')}")
        # print first 200 chars to see what it looks like
        print(content[:200].strip())
        print("-" * 50)
