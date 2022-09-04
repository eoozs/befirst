# befirst

A simple service for getting alerted when there are new website posts.

I built this while searching for a new home, but you can write a crawler 
and use it for any website with similar structure.


## About Supported Websites

Posts should be sorted descending based on creation date on the target website.

If you want to add support for another website add a new crawler under `crawlers/` and register it under `cmd/befirst`.
Similarly, if you want to add support for another notifications platform add it under `notify/`.


## How it works

Here is a brief overview of how it works, it will help you if you want to write a crawler and use this service.

Let's assume that we start the service and this is the first response from a crawler called `houses`:
```
houses = [
    {"listed-at": "2018-01-25", "id": 10},
    {"listed-at": "2018-01-22", "id": 9}
]
```

Then a variable called `last_post_id` will become `10` and nothing will happen.

After some time (duration is configured in settings) the service will fetch again the most recent listings:

```
houses = [
    {"listed-at": "2018-01-27", "id": 12},
    {"listed-at": "2018-01-26", "id": 11},
    {"listed-at": "2018-01-25", "id": 10},
    {"listed-at": "2018-01-22", "id": 9}
]
```

1. Two notifications are sent for posts with id `12` and `11`.
2. Then we see the post with id `10` which is the `last_post_id` so we stop and set `12` as the `last_post_id`.
