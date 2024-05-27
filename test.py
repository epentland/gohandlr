import requests

url = 'http://localhost:8087/user/43'

data = {
    'name': 'John Doe',
    'email': 'hello@gmail.com',
    'age': 25
}

response = requests.post(url, json=data)

print(response.json())