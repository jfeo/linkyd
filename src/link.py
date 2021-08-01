# -*- coding: utf-8 -*-

"""
linkyd.links
~~~~~~~~~~~~

This module contains the Link and LinkCollection classes.
"""

from datetime import datetime
from threading import Lock
from json import dump, load

class Link:
    """A link object that has been shared on the server.

    Contains a uri, the name of whoever added it, and the time it was added.
    """

    def __init__(self, uri: str, name: str):
        self.id = None
        self.uri = uri
        self.name = name
        self.added = datetime.now()

    def display_uri(self):
        """Returns a display uri."""
        if not self.uri.startswith('http'):
            return 'https://' + self.uri
        return self.uri

    def serialize(self):
        """Returns a serializable object."""
        return {
            'id': self.id,
            'uri': self.uri,
            'name': self.name,
            'added': self.added.isoformat()
        }
    
    @staticmethod
    def parse(dump):
        link = Link(dump["uri"], dump["name"])
        link.id = dump["id"]
        link.added = datetime.fromisoformat(dump["added"])
        return link


class LinkCollection():
    """A thread safe collection of links, that supports appending links."""

    def __init__(self):
        self.lock = Lock()
        self.auto_increment = 1
        self.links = {}
        self.name_index = {}

    def append(self, link: Link):
        """Add a new link to the link collection. Returns None."""
        with self.lock:
            link.id = self.auto_increment
            self.links[self.auto_increment] = link
            if link.name in self.name_index:
                self.name_index[link.name].add(self.auto_increment)
            else:
                self.name_index[link.name] = set([self.auto_increment])
            self.auto_increment += 1

    def items(self):
        return self.links.items()
    
    def serialize(self):
        return [self[link_id].serialize() for link_id in self]
    
    @staticmethod
    def parse(dump):
        links = LinkCollection()
        for link in dump:
            links.append(Link.parse(link))
        return links

    @staticmethod
    def load():
        with open("dump.json", "r") as dump_file:
            return LinkCollection.parse(load(dump_file))
    
    def dump(self):
        with self.lock:
            with open("dump.json", "w") as dump_file:
                dump(self.serialize(), dump_file)

    def __getitem__(self, link_id: int):
        return self.links[link_id]

    def __delitem__(self, link_id: int):
        with self.lock:
            link = self.links[link_id]
            # remove from name index
            self.name_index[link.name].remove(link_id)
            del self.links[link_id]

    def __len__(self):
        return len(self.links)

    def __iter__(self):
        return self.links.__iter__()
