


<div style="height: 150px"></div>

# Basic Redis Caching Demo

This app returns the number of repositories a Github account has. When you first search for an account, the server calls Github's API to return the response. This can take 100s of milliseconds. The server then adds the details of this slow response to Redis for future requests. When you search again, the next response comes directly from Redis cache instead of calling Github. The responses are usually usually in a millisecond or so making it blazing fast.

## Try it out

#### Deploy to Heroku

<p>
    <a href="https://heroku.com/deploy" target="_blank">
        <img src="https://www.herokucdn.com/deploy/button.svg" alt="Deploy to Heorku" />
    </a>
</p>

#### Deploy to Google Cloud
<p>
    <a href="https://deploy.cloud.run" target="_blank">
        <img src="https://deploy.cloud.run/button.svg" alt="Run on Google Cloud" width="150px"/>
    </a>
</p>


## How it works?

![How it works](docs/screenshot001.png)


### 1. How the data is stored:
```
SETEX microsoft 3600 1000
```

### 2. How the data is accessed:
```
GET microsoft
```

## How to run it locally?

#### Copy `.env.sample` to create `.env`. And provide the values for environment variables if needed

#### Run application

```sh
docker-compose up -d
```

Follow: http://localhost:5000
