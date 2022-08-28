from mimetypes import init
import pycurl
from os.path import exists
import requests
import threading
import re
import colorama

claimstatus = False
induser = 0
print('Checking for a target list...')
if exists('targets.txt'):
    print('Target list has been found!')
    users = []
    claimstatus = True
    for line in open('targets.txt', 'r'):
        users.append(line.strip())
    email = input('Email: ')
    session_id = input('Session ID: ')
    csrf_token = input('CSRF Token: ')
    threads = int(input('Threads: '))
else:
    print('No target list detected:')
    target = input('Target: ')
    email = input('Email: ')
    session_id = input('Session ID: ')
    csrf_token = input('CSRF Token: ')
    threads = int(input('Threads: '))

def claim(token, email, ssid, users=[]):
    global target
    global induser

    c = pycurl.Curl()
    c.setopt(pycurl.URL, 'https://solo.to/account/update-info/1')
    c.setopt(pycurl.HTTPHEADER, ["host: solo.to",
        "connection: keep-alive",
        'sec-ch-ua:  "Not A;Brand";v="99", "Chromium";v="102", "Google Chrome";v="102"',
        f"x-csrf-token: {token}",
        "sec-ch-ua-mobile: ?0",
        "user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36",
        "content-type: application/x-www-form-urlencoded; charset=UTF-8",
        "accept: */*",
        "x-requested-with: XMLHttpRequest",
        'sec-ch-ua-platform: "macOS"',
        "origin: https://solo.to",
        "sec-fetch-site: same-origin",
        "sec-fetch-mode: cors",
        "sec-fetch-dest: empty",
        "referer: https://solo.to/account",
        "accept-encoding: utf-8",
        "accept-language: en-GB,en-US;q=0.9,en;q=0.8",
        f"cookie: soloto_session={ssid}"])

    if claimstatus == True:
        for i in range(len(users)):
            try:
                c.setopt(pycurl.POST, 1)
                c.setopt(pycurl.POSTFIELDS, data)
                c.perform()
            except:
                print('error')
    else:
        try:
            data = f"_token={token}&email={email}&username={target}&domain="
            c.setopt(pycurl.POST, 1)
            c.setopt(pycurl.POSTFIELDS, data)
            c.perform()

        except:
            print('poor')

if claimstatus == True:
    for i in range(threads):
        t = threading.Thread(target=claim, args=(csrf_token, email, session_id, users))
        t.daemon = True
        t.start()
if claimstatus == False:
    for i in range(threads):
        t = threading.Thread(target=claim, args=(csrf_token, email, session_id, target))
        t.daemon = True
        t.start()
