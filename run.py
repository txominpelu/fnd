import requests
import json
import xml.etree.ElementTree as ET
import html2text
import sys
import os

def htmlstr2text(html):
    h = html2text.HTML2Text()
    # Ignore converting links from HTML
    h.ignore_links = True
    return h.handle(html)

def add_field(item, field, entry):
    val = [x.text.strip() for x in item.iter(field)]
    if len(val) > 0:
        entry[field] = val[0]
    else:
        entry[field] = ""

def get_entries():
    entries = []
    resp = requests.get("https://remoteok.io/remote-dev-jobs.rss")
    xml_resp = resp.text
    root = ET.fromstring(xml_resp)

    fields = ['title', 'company', 'description']
    for item in root[0].iter('item'):
        entry = {}
        add_field(item, 'guid', entry)
        add_field(item, 'title', entry)
        add_field(item, 'company', entry)
        add_field(item, 'location', entry)
        add_field(item, 'description', entry)
        if 'description' in entry:
            entry['description'] = htmlstr2text(entry['description'])
            entry['first_line_description'] = first_line(entry['description'])
        entries.append(entry)
    return entries

cache_file = "{}/.rok".format(os.environ['HOME'])

def first_line(s):
    return s.replace('\n','')

def list_entries():
    entries = get_entries()
    out = {}
    i = 0
    for entry in entries:
        print("{guid}.{title} - {company}.{first_line_description}".format(**entry))
        out[entry['guid']] = entry
        i += 1
    with open(cache_file, 'w+') as f:
        json.dump(out, f)

def get_description(line):
    guid = line.partition('.')[0]
    with open(cache_file, 'r') as f:
        entries = json.load(f)
        print("{} - {}".format(entries[guid]['title'], entries[guid]['company']))
        print("------------------------------------")
        print(entries[guid]['description'])


if __name__ == '__main__':
    if sys.argv[1] == 'list':
        list_entries()
    elif sys.argv[1] == 'description':
        get_description(sys.argv[2])
