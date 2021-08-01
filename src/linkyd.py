# -*- coding: utf-8 -*-

from flask import Flask, flash, request, render_template, redirect, url_for, jsonify
from flask import make_response
from link import Link, LinkCollection
from text import TEXT
import os

app = Flask(__name__)
app.secret_key = b'_5#y2L"F4Q8z\n\xec]/'

try:
    links = LinkCollection.load(os.open('dump.json', 'r'))
    os.remove("dump.json")
except:
    links = LinkCollection()

@app.route('/', methods=['GET'])
def webapp_index():
    name = request.cookies.get('name')
    return render_template('index.html', TEXT=TEXT, links=links.items(), name=name)


@app.route('/add', methods=['POST'])
def webapp_add():
    form = request.form
    error = False
    resp = make_response(redirect(url_for('webapp_index')))
    if 'name' not in form or not form['name']:
        flash(TEXT['FLASH_MISSING'].format('name'), 'danger')
        error = True
    if 'uri' not in form or not form['uri']:
        flash(TEXT['FLASH_MISSING'].format('uri'), 'danger')
        error = True
    if not error:
        link = Link(uri=form['uri'], name=form['name'])
        links.append(link)
        flash(TEXT['FLASH_ADDED_SUCCESS'], 'success')
        resp.set_cookie('name', form['name'])
    return resp


@app.route('/delete/<int:link_id>', methods=['GET'])
def webapp_delete(link_id):
    try:
        del links[link_id]
    except KeyError:
        flash(TEXT['FLASH_DELETE_ERROR'], 'danger')
        return redirect(url_for('webapp_index'))
    flash(TEXT['FLASH_DELETE_SUCCESS'], 'success')
    return redirect(url_for('webapp_index'))

@app.route('/dump', methods=['GET'])
def webapp_dump():
    try:
        links.dump()
        flash(TEXT['FLASH_DUMP_SUCCESS'], 'success')
    except Exception as e:
        flash(TEXT['FLASH_DUMP_ERROR'].format(repr(e)), 'danger')
    finally:
        return redirect(url_for('webapp_index'))

@app.route('/restore', methods=['GET'])
def webapp_restore():
    return render_template('restore.html', TEXT=TEXT)

@app.route('/restore', methods=['POST'])
def webapp_post_restore():
    global links
    if "dumpFile" not in request.files:
        flash(TEXT['FLASH_RESTORE_NO_FILE'], 'danger')
        return redirect(url_for('webapp_restore'))
    
    dumpFile = request.files['dumpFile']
    try:
        loaded_links = LinkCollection.load(dumpFile)
        links = loaded_links
        flash(TEXT['FLASH_RESTORE_SUCCESS'], 'success')
        return redirect(url_for('webapp_index'))
    except Exception as e:
        flash(TEXT['FLASH_RESTORE_ERROR'].format(repr(e)), 'danger')
        return redirect(url_for('webapp_restore'))

@app.route('/api/links', methods=['GET'])
def get_links():
    return {
        'links': list([link.serialize() for (link_id, link) in links.items()]),
        'count': len(links)
    }

@app.route('/api/links', methods=['POST'])
def post_links():
    data = request.get_json()
    if data is None:
        return {'message': TEXT['REST_MSG_JSON_ERROR']}, 400
    if 'name' not in data:
        return {'message': TEXT['REST_MSG_JSON_MISSING'].format('name')}, 400
    if 'uri' not in data:
        return {'message': TEXT['REST_MSG_JSON_MISSING'].format('uri')}, 400
    link = Link(uri=data['uri'], name=data['name'])
    links.append(link)
    return {'message': TEXT['REST_MSG_ADDED_SUCCESS']}


@app.route('/api/links/<int:link_id>', methods=['DELETE'])
def delete_link(link_id):
    try:
        link = links[link_id]
        del links[link_id]
    except KeyError:
        return {'message': TEXT['REST_MSG_DELETE_ERROR']}
    return {'message': TEXT['REST_MSG_DELETE_SUCCESS'], 'link': link.serialize()}

@app.route('/api/dump', methods=['GET'])
def get_dump():
    return jsonify(links.serialize())