Melange Javascript API
---

# Introduction

This document is intended to be a reference for the Melange Javascript API
for use in plugin development.

This is version `0.0.10` of this document.

# Using the Melange API

The Melange API can be included in any Melange application with the
following HTML:

    <script src="http://common.melange/js/melange/0.0.1"></script>

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

The Message data-type is used to represent a message sent to someone
in a Melange application. It is utilized in all of the API calls below
as `msg`.

The following is a description of the fields:

- `to`: an array of Address types that the message will be sent
  to. This may be an empty array if the message will be public.
- `name`: the name (string) of the message being sent. **Currently, it
  is the developer's responsibility to ensure that this name does not
  conflict with any other plugin.** We recommend that you prefix the
  `id` of the message with the name of the plugin. For example, a
  message from the chat application with identification number 1 could
  have `name` = "chat/1"
- `date`: the date/time that the message was sent in ISO
  format. Generally, this should be set to `(new
  Date()).toISOString()` as shown above.
- `public`: this field determines whether or not the receiver(s) will
  get a notification of the message. If the receiver will not get a
  notification, they must be "subscribed" to the sender in order to
  see the message.
- `self`: *Not required for outgoing messages.* For messages you
  receive from the Melange API, messages sent by the current user will
  have `self` set to `true`.
- `components`: a dictionary of named components to their string
  values. Effectively, this is the body of the message, each "item"
  that you wish to track in the message should have its own field in
  the dictionary.

### API Response

    {
      error: {
        code: 0,
        message: "",
      }
    }

This is a generic response received by the API for "action"
requests. If the `error.code` field is non-zero, then an error has
occurred. The `error.message` field will contain information.

# Message Management

These methods are generally for interacting with Melange messages. You
will primarily use them when building Melange plugins.

### melange.createMessage(msg, callback)

`createMessage` takes a Message type and a function callback of type
`function(api_response)`. It will attempt to publish the `msg` given
(according to the type defined above) as the current user, and it will
invoke the callback upon completion (or error).

Example:

    melange.createMessage(
        {
            to: [{ alias: "hleath@airdispatch.me" }],
            name: "documentation/test-create",
            date: (new Date()).toISOString(),
            public: false,
            components: {
                "getmelange.com/developer/method": "createMessage",
            },
        }, function(response) {
            // Check for errors!
            if (response.error.code !== 0)
                alert(response.error.message);
            })

### melange.updateMessage(msg, id, callback)

`updateMessage` will replace a message named `id` with the new version
`msg`. It takes three arguments:

1. The new Message as `msg`.
2. The `name` of the message that you are replacing as `id`.
3. A callback that receives a generic API response as `callback`.

**Note: The new message's name must match the id of the message you
  are replacing (`msg.name == id`).**

Example:

    melange.updateMessage(
        {
            to: [{ alias: "hleath@airdispatch.me" }],
            name: "documentation/test-create",
            date: (new Date()).toISOString(),
            public: false,
            components: {
                "getmelange.com/developer/method": "updateMessage",
            },
        }, "documentation/test-create",
            function(response) {
            // Check for errors!
            if (response.error.code !== 0)
                alert(response.error.message);
            })

### melange.findMessages(fields, predicate, callback [, realtime])

`findMessages` will return an array of Messages sent or received by
the user matching the given `fields` and `predicate`. It takes the
following arguments:

1. `fields`: An array of component names (see Message data type
   definition). All Messages returned by `findMessages` will have
   components with names that match those specified in `fields`. An
   optional field may be specified by starting the field name with a
   `?` in the array.
2. `predicate`: **CURRENTLY UNUSED** Just set to `{}`. In the future,
   it will be used to filter messages further.
3. `callback`: A function called on success (in which case, it will be
   passed an array of Messages) or failure (in which case, it will be
   passed an API Response with a non-zero error code).
4. `realtime`: **OPTIONAL** If supplied, `realtime` is a callback that
   accepts one argument - a Message. The callback will be called every
   time a new message is *received* that has the necessary fields.

Example:

    melange.findMessages(
        [
            "getmelange.com/developer/method",
            "?getmelange.com/developer/optional"
        ],
        {},
        function(msgs) {
            if (msgs.error !== undefined) {
                // Something Happened
                alert(msgs.error.message);
                return;
            }

            // msgs is an array of Message Type
        },
        function(msg) {
            alert("Received New Message");
            // msg is of Message Type
        }
    )

### melange.downloadMessage(address, id, callback)

`downloadMessage` will download a specific message named `id` from the
Melange address `address`. Like `findMessages` callback will either
receive the retrieved Message object *or* an API response that
includes a non-zero error code.

Example:

    melange.downloadMessage("hleath@airdispatch.me",
    "documentation/test-create",
        function(msg) {
            if (msg.error !== undefined) {
                alert(msg.error.message);
                return;
            }

            // msg is of type Message
        })

### melange.downloadPublicMessages(fields, predicate, addr, callback)

`downloadPublicMessages` will download all public messages from an
address that satisfy the same `fields` and `predicate` filtering as in
`findMessages`.

# Viewer Management

These API calls are exclusively for "Viewers" of plugins. We will be
posting a guide on what a viewer is and how to build one shortly. In
essence, a "Viewer" is utilized whenever Melange needs to display a
custom message (in a newsfeed, on a profile). The "Viewer" determines
how to present that content to the user.

### melange.viewer(callback)

`viewer` will initialize Melange's viewer system. `callback` is a
function that accepts one argument - the Message that Melange wants
the viewer to display.

### melange.refreshViewer()

`refreshViewer` will tell Melange that the content is done loading,
and Melange should expand the Viewer's dimensions to its new
content. Since all viewers live in iFrames, this is necessary to size
the window correctly. If you do not call this method after loading
the content, Melange may cut off some of the Viewer's information.

# Melange Management

These methods are for interacting with Melange itself. They are used
when you need more context information, or have to perform an action
that only Melange can do.

### melange.currentUser(callback)

`currentUser` will fetch the "Address" information of the current
user and return it to `callback`. `callback` is a function that
accepts one argument - the "Address" of the current user.

### melange.openLink(url)

`openLink` will open the `url` provided in the user's default web browser.

# Data Proxy

Users can now upload data and images to their servers. Plugins have a
very simple way of accessing this content.

Simply construct a URL of the following form where `addr` is the
address of the person who uploaded the content and `id` is the "name"
of the content. It will be served like normal.

    http://data.melange/addr/id

# Angular Extensions

Since much of Melange is built on Angular, we decided to include some
extras for people who wanted to use Angular in their plugins.

### mlgToField

`mlgToField` is a directive that presents an auto-completing (from
contacts) fields for users to determine who to send a message to. See
its use in the
[Notes Plugin](https://github.com/melange-app/plugin-notes/blob/master/templates/new.html).

It is provided in the Angular module `melangeUi`.

### melange.angularCallback(callback)

`angularCallback` will wrap `callback` in a `$scope.$apply()` so that
your views will update when the callback is finally called.
