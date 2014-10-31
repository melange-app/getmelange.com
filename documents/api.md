Melange Javascript API
---

# Introduction

This document is intended to be a reference for the Melange Javascript API
for use in plugin development.

This is version `0.0.5` of this document.

# Data Types

### Address

    {
        alias: "hleath@airdispatch.me",
        fingerprint: "003208549081345hadsfjkh",
    }

### Message

    {
      to: [{ alias: "hleath@airdispatch.me" }],
      name: "chat/1",
      date: (new Date()).toISOString(),
      public: false,
      self: false,
      components: {
        "airdispat.ch/chat/body": "Hello, world",
      },
    }

# Message Management

### melange.createMessage(msg, callback)

### melange.updateMessage(msg, id, callback)

### melange.findMessages(fields, predicate, callback, realtime)

### melange.downloadMessage(address, id, callback)

### melange.downloadPublicMessages(fields, predicate, addr, callback)

# Viewer Management

### melange.viewer(callback)

### melange.refreshViewer(msg)

# Melange Management

### melange.currentUser(callback)

### melange.openLink(url)

# Data Proxy

# Angular Extensions

### mlgToField

### .angularCallback
