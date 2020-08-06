linkyd
======

A small HTTP REST server, for storing and sharing links and other small bits of
text. Based on Flask.

For example to share links on a local network, with only trusted entities
(yikes).

Use absolutely at your own risk :-)

Install
-------

This project is made with *Pipenv*. Install dependencies with

```
pipenv install
```

Usage
-----

Can be run (at least for development purposes) with:

```
export FLASK_APP=src/linkyd.py
pipenv run flask run
```
