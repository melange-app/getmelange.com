Creating a Melange App
---

# Introduction

Melange Applications are a fast and easy way to develop secure,
decentralized systems. In this guide, you will learn how Melange
treats plugins, and how to create a simple application.

# Plugin Information

### Installing Plugins

To install a plugin that you have developed, you must simply drop the
folder with the unique Plugin Id into the Plugins folder.

- OS X: `~/Library/Application Support/Melange/plugins`
- Windows: `%APPDATA%\Melange\plugins`

### Plugin File Structure

    - com.getmelange.developer_plugin
    |
    |- package.json
    |- index.html
    |- HTML Files
    |- JS Files
    |- CSS Files
    |- etc.

Your plugin will be loaded just like a website (using `index.html` as
the main page). You can access resources using relative paths
(`/styles.css`). The only difference between a regular website and a
Melange plugin is the `package.json` file.

# Manifest File

### package.json

The `package.json` file is modeled after the `package.json` files in
npm repositories with a few additions for Melange.

Example:

    {
        "id": "com.github.melange-app.plugin-notes",
        "name": "Notes",
        "description": "Publish notes to people who follow you.",
        "version": "0.0.5",
        "permissions": {
            "read-message": ["airdispat.ch/notes/title", "airdispat.ch/notes/body", "airdispat.ch/notes/date"],
            "send-message": ["airdispat.ch/notes/title", "airdispat.ch/notes/body", "airdispat.ch/notes/date"]
        },
        "viewers": {
            "default": {
                "view": "viewer/viewer.html",
                "type": ["airdispat.ch/notes/title", "airdispat.ch/notes/body"]
            }
        },
        "tiles": {
            "status-updater": {
                "name": "Quick Status Updater",
                "description": "Quickly publish status updates at the top of your dashboard.",
                "view": "tile/tile.html",
                "click": true,
                "size": "100%x100"
            }
        },
        "hideSidebar": true,
        "author": {
            "name": "Hunter Leath",
            "email": "h@hunterleath.com"
        },
        "homepage": "http://airdispat.ch/plugins/notes"
    }

Detailed Description:

- `id`: The unique identifying string of the plugin. In general, this
  is set like an Apple Bundle Identifier (using a reverse-domain
  string). For example, if your personal website is `getmelange.com`
  and you are developing a plugin called `new_plugin`, you would set
  the `id` field to `com.getmelange.new_plugin`.
- `name`: The displayed name for the plugin.
- `description`: The displayed description for the plugin.
- `version`: The [Semantic Version](http://semver.org/) for the
  plugin.
- `permission`: An object that represents the permissions for the
  plugin. Read the section below on `permissions`.
- `viewers`: An object that represents the viewers for the
  plugin. Read the section below on `viewers`.
- `tiles`: An object that represents the tiles for the plugin. Read
  the section below on`viewers`.
- `hideSidebar`: **OPTIONAL** If this is set to `true`, then your
  plugin will not show up in the sidebar for Melange like a normal
  plugin. Use this option if you plan on creating a plugin that only
  contains tiles or viewers.
- `author`: Used in the same way as the npm repository author,
  `author.name` is displayed to the user as the author of the plugin.
- `homepage`: The homepage for the plugin, if one exists.

### Permissions

The permissions object is used to grant your plugin permission to read
and write users' Melange messages. Since all Messages contain a list
of named components (see the API guide), permissions involve
specifying which named components the plugin is able to read and
write.

Currently, only two permissions are allowed:

- `read-message`: Should be an array of named components that
  correspond to the components this plugin intends to read.
- `send-message`: Should be an array of named components that
  correspond to the components this plugin intends to create.

If the plugin attempts to call `findMessages` for a component name
that it does not have access to, an exception will be raised. If a
plugin would receive a message in `findMessages` that contains data it
does not have permission to read, it will be stripped before being
passed to the plugin. An exception will be raised if a plugin attempts
to create a message with components it does not have permission to
write.

### Viewers

Viewers are a mechanism in Melange that allows plugins to display
custom content in users' newsfeeds and profiles.

### Tiles

Tiles are parts of a plugin that users may decide to pin to the top of
their newsfeed in the dashboard. These "Tiles" may do whatever the
plugin author wishes. For example, the "Status" tile allows users to
post simple messages while the "News" tile displays static information
that the user may click to read more.
