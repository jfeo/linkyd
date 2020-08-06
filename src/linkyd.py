# -*- coding: utf-8 -*-

from flask import Flask, flash, request, render_template, redirect, url_for, make_response
from link import Link, LinkCollection

app = Flask(__name__)
app.secret_key = b'_5#y2L"F4Q8z\n\xec]/'

links = LinkCollection()


@app.route('/', methods=['GET'])
def webapp_index():
    return render_template('index.html', links=links.items(), name=request.cookies.get('name'))

@app.route('/add', methods=['POST'])
def webapp_add():
    form = request.form
    error = False
    resp = make_response(redirect(url_for('webapp_index')))
    if 'name' not in form or not form['name']:
        flash('Missing required value <strong>name</strong>.', 'danger')
        error = True
    if 'uri' not in form or not form['uri']:
        flash('Missing required value <strong>uri</strong>.', 'danger')
        error = True
    if not error:
        link = Link(uri=form['uri'], name=form['name'])
        links.append(link)
        flash('Added the link! Hurray!', 'success')
        resp.set_cookie('name', form['name'])
    return resp

@app.route('/delete/<int:link_id>', methods=['GET'])
def webapp_delete(link_id):
    try:
        link = links[link_id]
        del links[link_id]
    except IndexError:
        flash('Could not delete the link', 'danger')
        return redirect(url_for('webapp_index'))
    flash('Deleted the link! Hurray!', 'success')
    return redirect(url_for('webapp_index'))


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
        return {'message': 'missing or invalid json'}, 400
    if 'name' not in data:
        return {'message': 'missing required value \'name\''}, 400
    if 'uri' not in data:
        return {'message': 'missing required value \'uri\''}, 400
    link = Link(uri=data['uri'], name=data['name'])
    links.append(link)
    return {'message': 'successfully added'}


@app.route('/api/links/<int:link_id>', methods=['DELETE'])
def delete_link(link_id):
    link = links[link_id]
    del links[link_id]
    return {'message': 'successfully deleted', 'link': link.serialize()}
