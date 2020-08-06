# -*- coding: utf-8 -*-

"""
linkyd.links
~~~~~~~~~~~~

This module contains the Link and LinkCollection classes.
"""

from datetime import datetime
from threading import Lock


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
