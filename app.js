var request = require('request')
var debug = require('debug')('rs')
var async = require('async')
var promclient = require('prometheus-client')

var host = process.env.RANCHER_HOST
var port = process.env.RANCHER_PORT || 8080
var username = process.env.RANCHER_API_ACCESS_KEY
var password = process.env.RANCHER_API_SECRET_KEY

process.on('SIGTERM', function () {
    process.exit(0);
});

if (!host || !username || !password) {
    console.error('Missing parameters for host, username or password')
    process.exit(1)
}

var client = new promclient()
var gauges = {}

function createGauge(name) {
    gauges[name] = client.newGauge({
        namespace: 'rancher_status',
        name: name,
        help: '1 if all containers in an environment are active'
    })
}

function updateGauge(name, value) {
    if (!gauges[name]) {
        createGauge(name)
    }
    gauges[name].set({
        name: name
    }, value)
}

function updateMetrics() {
    debug('updating metrics')
    poll(function(err, results) {
        if (err) {
            throw err
        }
        debug('got results %o', results)
        Object.keys(results).forEach(function(name) {
            var state = results[name]
            name = name.replace('-', '_')
            updateGauge(name, state == 'active' ? 1 : 0)
        });
    });
}

setInterval(updateMetrics, 5000)
updateMetrics()
client.listen(9010)


function poll(callback) {
    var envIdMap = {}

    async.waterfall([
        function(next) {
            var uri = 'http://' + host + ':' + port + '/v1/projects'
            jsonRequest(uri, function(err, json) {
                if (err) {
                    return next(err)
                }
                var environments = json.data[0].links.environments
                next(null, environments)
            })
        },
        function(environmentsUrl, next) {
            jsonRequest(environmentsUrl, function(err, json) {
                if (err) {
                    return next(err)
                }
                var servicesUrl = json.data.map(function(raw) {
                    return raw.links.services
                });
                json.data.forEach(function(env) {
                    envIdMap[env.id] = env.name
                });
                next(null, servicesUrl)
            });
        },
        function(servicesUrls, next) {
            var tasks = servicesUrls.map(function(servicesUrl) {
                return function(next) {
                    jsonRequest(servicesUrl, next)
                }
            });

            async.parallel(tasks, function(err, results) {
                var data = results.map(function(servicesRaw) {
                    return servicesRaw.data
                });

                next(null, data)
            });
        },
        function(servicesData, next) {
            var services = servicesData.map(function(stackServices) {
                return stackServices.map(function(service) {
                    return {
                        name: service.name,
                        state: service.state,
                        environment: envIdMap[service.environmentId]
                    }
                });
            });

            var flattened = []
            services.forEach(function(service) {
                flattened = flattened.concat(service)
            });

            next(null, flattened)
        },
        function(serviceData, next) {
            var envState = {}
            serviceData.forEach(function(service) {
                if (!envState[service.environment]) {
                    envState[service.environment] = service.state
                } else if (service.state != 'active') {
                    envState[service.environment] = service.state
                }
            });
            next(null, envState)
        }
    ], function(err, results) {
        callback(err, results)
    })
}

function jsonRequest(uri, callback) {
    debug('send request: %s', uri)

    request({
        uri: uri,
        headers: {
            'Accept': 'application/json'
        },
        auth: {
            user: username,
            pass: password,
            sendImmediately: true
        }
    }, function(err, response, body) {
        if (err) {
            debug('got error response')
            return callback(err)
        }

        debug('got valid response with status %s', response.statusCode)

        var data;
        try {
            data = JSON.parse(body)
        } catch(e) {
            debug('could not JSON decode response body')
            var error = new Error('could not decode')
            error.response = response
            error.body = body
            return callback(error)
        }

        return callback(null, data)
    })
}
