import requests


r = requests.get('http://localhost:8100/all_users')
print(r.text)


# data = '{"resourceType":"Patient","id":"foo-bar","active":true,"name":[{"use":"official","family":["Плис","Александрович"],"given":["Виктор"]}],"gender":"male","birthDate":"1996-06-26","address":[{"type":"both","use":"home","city":"Сибай","line":["ул. Революционная, д. 10"],"district":"Россия"}],"telecom":[{"use":"mobile","rank":1,"value":"89631338188","system":"phone"}]}'
# r = requests.post('http://localhost:8100/user/foo-bar', data={'patient': data})
# print(r.text)
#
#
# r = requests.get('http://localhost:8100/user/foo-bar')
# print(r.text)