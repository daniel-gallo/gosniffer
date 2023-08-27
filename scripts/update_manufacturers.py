import os
import re
import sqlite3
from urllib import request

DB_URL = 'https://github.com/joemiller/mac-to-vendor/raw/main/db/macaddrs.sqlite'
DB_FILENAME = 'macaddrs.sqlite'

OUTPUT_FILENAME = '../lanscanner/manufacturers.txt'


def download_db():
    request.urlretrieve(DB_URL, DB_FILENAME)


def get_vendors():
    con = sqlite3.connect(DB_FILENAME)

    for mac_prefix, manufacturer in con.execute('SELECT * FROM vendors'):
        # There is a row with addr='#' and vendor=<null>
        if manufacturer is None:
            continue

        bytes = re.findall('..', mac_prefix.lower())
        mac_prefix = ':'.join(bytes)

        yield mac_prefix, manufacturer

    con.close()


def db_to_file():
    with open(OUTPUT_FILENAME, 'w') as f:
        for mac_prefix, manufacturer in get_vendors():
            f.write(f'{mac_prefix}\t{manufacturer}\n')


def delete_db():
    os.remove(DB_FILENAME)


if __name__ == '__main__':
    download_db()
    db_to_file()
    delete_db()
