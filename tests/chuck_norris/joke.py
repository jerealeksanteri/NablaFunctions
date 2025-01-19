import requests

response = requests.get("https://api.chucknorris.io/jokes/random")

print(response.text)