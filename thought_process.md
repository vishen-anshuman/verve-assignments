## Primary Task

1. Understanding the product requirements on top level. It requires <br>
    1. A rest api. <br>
    2. A logger function which will generate create a log file<br>
    3. A function/service which will log the unique count per minute<br>
2. Orchestration of above requirements <br>
    1. Maintaining the unique count requires a cache, though the problem statement says <br>
       processing 10k request per minute which requires a distributed cache like redis, for simplicity <br>
       I assumed a local cache via map.
    2. Cache, Logger, Mutex (for avoiding race conditions) are set at application level.<br>
       Initiating the app will instantiate a Cache and create a log file.
    3. Implementing the handler function for the api, meeting the requirements like, validating input <br>
       adding all ids to cache, and hitting and logging the endpoint response.
    4. Implementing a function which will log the count of unique ids to a log file for each minute.


## Extension 1
Its requirement was pretty simple, for `endpint` coming in api request post calls to be made

## Extension 2
As the service is getting into distributed mode, a local cache will never be good,<br>
so instead of local cache `redis` cache is used. **Extension 2 is implemented on top of Primary task**

## Extension 3
Instead of logging the unique count to log, it has to be streamed. I have used `Kafka` as the streaming service.<br>
**Extension 3 is implemented on top of Extension 2**